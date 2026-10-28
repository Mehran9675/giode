package giode

import (
	"bufio"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mehran9675/giode/update"
)

// The single-instance lock is a file in the temporary directory holding
// the owner's pid, the loopback port it listens on and a random token.
// A second instance hands its command-line arguments to the owner over
// that port, proving it can read the file (which only the user can)
// with the token.

// instanceLock is the owner's side of the lock.
type instanceLock struct {
	ln    net.Listener
	token string
}

// acquireInstanceLock claims a single-instance lock for the given
// application id. It reports whether another live instance already
// holds the lock.
func acquireInstanceLock(id string) (second bool, path string, owner *instanceLock) {
	waitForUpdatedInstance()
	if id == "" {
		id = "giode"
	}
	path = filepath.Join(os.TempDir(), "giode-"+sanitizeID(id)+".lock")

	// The lock owner writes its pid; a stale file (dead process or
	// one that crashed before writing) is reclaimed.
	for attempt := 0; attempt < 5; attempt++ {
		data, err := os.ReadFile(path)
		if err != nil {
			break
		}
		pid, _, _ := parseLock(data)
		if pid > 0 && processAlive(pid) {
			return true, path, nil
		}
		// Empty or dead: give a racing writer a moment to fill in its
		// pid before reclaiming.
		if attempt < 4 {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		os.Remove(path)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		// Lost the race with another instance.
		return true, path, nil
	}
	defer f.Close()
	owner = &instanceLock{}
	port := 0
	if ln, err := net.Listen("tcp", "127.0.0.1:0"); err == nil {
		var tok [16]byte
		rand.Read(tok[:])
		owner.ln, owner.token = ln, hex.EncodeToString(tok[:])
		port = ln.Addr().(*net.TCPAddr).Port
	}
	fmt.Fprintf(f, "%d\n%d\n%s\n", os.Getpid(), port, owner.token)
	return false, path, owner
}

// waitForUpdatedInstance waits, in an app restarted by an update, for
// the old process to exit, so the new one takes over its lock.
func waitForUpdatedInstance() {
	pid, err := strconv.Atoi(os.Getenv(update.WaitEnv))
	if err != nil {
		return
	}
	os.Unsetenv(update.WaitEnv)
	for deadline := time.Now().Add(15 * time.Second); processAlive(pid) && time.Now().Before(deadline); {
		time.Sleep(50 * time.Millisecond)
	}
	update.Cleanup()
}

// parseLock reads a lock file: pid, then port and token.
func parseLock(data []byte) (pid, port int, token string) {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	pid, _ = strconv.Atoi(strings.TrimSpace(lines[0]))
	if len(lines) >= 3 {
		port, _ = strconv.Atoi(strings.TrimSpace(lines[1]))
		token = strings.TrimSpace(lines[2])
	}
	return pid, port, token
}

// releaseInstanceLock removes the lock file of the owning instance.
func releaseInstanceLock(path string, owner *instanceLock) {
	if owner != nil && owner.ln != nil {
		owner.ln.Close()
	}
	if path != "" {
		os.Remove(path)
	}
}

// launchMessage is what a second instance sends the owner.
type launchMessage struct {
	Token string   `json:"token"`
	Args  []string `json:"args"`
}

// forwardLaunch sends args to the instance holding the lock at path.
func forwardLaunch(path string, args []string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, port, token := parseLock(data)
	if port == 0 {
		return fmt.Errorf("giode: the running instance doesn't take launches")
	}
	c, err := net.DialTimeout("tcp", "127.0.0.1:"+strconv.Itoa(port), 2*time.Second)
	if err != nil {
		return err
	}
	defer c.Close()
	c.SetDeadline(time.Now().Add(5 * time.Second))
	if err := json.NewEncoder(c).Encode(launchMessage{Token: token, Args: args}); err != nil {
		return err
	}
	reply, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		return err
	}
	if strings.TrimSpace(reply) != "ok" {
		return fmt.Errorf("giode: the running instance refused the launch")
	}
	return nil
}

// serveLaunches hands launches from second instances to handle.
func (l *instanceLock) serveLaunches(handle func(args []string)) {
	if l == nil || l.ln == nil {
		return
	}
	go func() {
		for {
			c, err := l.ln.Accept()
			if err != nil {
				return
			}
			go func() {
				defer c.Close()
				c.SetDeadline(time.Now().Add(5 * time.Second))
				var m launchMessage
				if err := json.NewDecoder(io.LimitReader(c, 1<<20)).Decode(&m); err != nil {
					return
				}
				if subtle.ConstantTimeCompare([]byte(m.Token), []byte(l.token)) != 1 {
					io.WriteString(c, "no\n")
					return
				}
				io.WriteString(c, "ok\n")
				handle(m.Args)
			}()
		}
	}()
}

// handleLaunch handles a second instance's arguments in the running
// one.
func (a *App) handleLaunch(args []string) {
	if a.runContextAction(args) {
		return
	}
	// A link the system opened with the app (its URL scheme).
	if len(args) == 1 && strings.Contains(args[0], "://") && a.cfg.OnOpenURL != nil {
		if u, err := url.Parse(args[0]); err == nil && u.Scheme != "" {
			a.openedURL(u)
			a.Post(func() {
				a.ShowWindow()
				a.Raise()
			})
			return
		}
	}
	a.Post(func() {
		if l := a.cfg.SingleInstanceLock; l != nil && l.OnLaunchArgs != nil {
			l.OnLaunchArgs(args)
		}
		a.ShowWindow()
		a.Raise()
	})
}

func sanitizeID(id string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			return r
		default:
			return '-'
		}
	}, id)
}

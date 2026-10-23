package giode

import (
	"context"
	"errors"
	"io"
	"slices"
	"strings"
	"testing"
	"time"

	"gioui.org/app"
	"gioui.org/x/explorer"
)

func TestExtensions_AddLeadingDot(t *testing.T) {
	got := extensions([]FileFilter{{Extensions: []string{"wav", ".mp3", ""}}, {Extensions: []string{"ogg"}}})
	if want := []string{".wav", ".mp3", ".ogg"}; !slices.Equal(got, want) {
		t.Errorf("extensions = %v, want %v", got, want)
	}
}

type closer struct{ closed chan struct{} }

func (c *closer) Read([]byte) (int, error) { return 0, io.EOF }
func (c *closer) Close() error             { close(c.closed); return nil }

func TestAwait_MapsErrors(t *testing.T) {
	ctx := context.Background()
	if _, err := await(ctx, func() (int, error) { return 0, explorer.ErrUserDecline }, func(int) {}); !errors.Is(err, ErrCanceled) {
		t.Errorf("decline: err = %v, want ErrCanceled", err)
	}
	if _, err := await(ctx, func() (int, error) { return 0, explorer.ErrNotAvailable }, func(int) {}); !errors.Is(err, ErrNotSupported) {
		t.Errorf("unavailable: err = %v, want ErrNotSupported", err)
	}
	if v, err := await(ctx, func() (int, error) { return 7, nil }, func(int) {}); v != 7 || err != nil {
		t.Errorf("success: %v, %v", v, err)
	}
}

func TestAwait_CancelClosesLateFile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	release := make(chan struct{})
	f := &closer{closed: make(chan struct{})}
	go func() { time.Sleep(10 * time.Millisecond); cancel() }()
	_, err := await(ctx, func() (io.ReadCloser, error) {
		<-release
		return f, nil
	}, closeFile[io.ReadCloser])
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	close(release)
	select {
	case <-f.closed:
	case <-time.After(time.Second):
		t.Error("file chosen after cancellation wasn't closed")
	}
}

func TestShareText_DesktopCopiesToClipboard(t *testing.T) {
	a := &App{window: new(app.Window)}
	if err := a.Share(context.Background(), ShareText{Text: "hello"}); err != nil {
		t.Fatal(err)
	}
	if a.clip.write == nil || *a.clip.write != "hello" {
		t.Error("ShareText didn't queue the clipboard write")
	}
}

func TestShareFile_MissingFile(t *testing.T) {
	a := &App{}
	err := a.Share(context.Background(), ShareFile{Path: "/does/not/exist.wav"})
	if err == nil || !strings.Contains(err.Error(), "exist") {
		t.Errorf("err = %v, want a missing-file error", err)
	}
}

func TestShareFile_MIMEGuess(t *testing.T) {
	if got := (ShareFile{Path: "a.wav"}).mimeType(); !strings.HasPrefix(got, "audio/") {
		t.Errorf("mime for .wav = %q", got)
	}
	if got := (ShareFile{Path: "a.unknownext"}).mimeType(); got != "application/octet-stream" {
		t.Errorf("mime for unknown = %q", got)
	}
}

func TestLifecycle_MinimizeOnDesktop(t *testing.T) {
	var events []string
	a := &App{window: new(app.Window), cfg: Config{
		OnBackground: func() { events = append(events, "bg") },
		OnForeground: func() { events = append(events, "fg") },
	}}
	a.handleConfig(app.Config{Mode: app.Minimized})
	a.handleConfig(app.Config{Mode: app.Windowed})
	a.handleConfig(app.Config{Mode: app.Maximized})
	if want := []string{"bg", "fg"}; !slices.Equal(events, want) {
		t.Errorf("events = %v, want %v", events, want)
	}
}

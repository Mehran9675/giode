# Platform services

Services apps reach outside their window: links, the safe area, the system theme, device
details, haptics, notifications, staying awake, secure storage, pickers and dictation.
File dialogs, sharing, the Back button and background/foreground callbacks are in
[window.md](window.md).

| | Android | iOS | macOS | Windows | Linux | Web |
| --- | --- | --- | --- | --- | --- | --- |
| `OpenURL` | ✓ | ✓ | ✓ | ✓ | ✓ (`xdg-open`) | ✓ |
| `SafeArea`, `EdgeToEdge` | ✓ | ✓ | no insets | no insets | no insets | no insets |
| `FollowSystemTheme`, `SystemDarkMode` | ✓ | ✓ | ✓ | ✓ | ✓ (XDG portal) | ✓ |
| `Device` | ✓ | ✓ | ✓ | ✓ | ✓ | OS only |
| `Haptic` | ✓ | ✓ | – | – | – | – |
| `Notify` | ✓ | – | ✓ | ✓ | ✓ | – |
| `KeepScreenOn` | ✓ | ✓ | ✓ | ✓ | ✓ | – |
| `KeepAlive` | manifest entry (see below) | ~30 s | ✓ | ✓ | ✓ (logind) | – |
| `SecureStore` | Keystore | Keychain | Keychain | DPAPI | Secret Service | – |
| `PickFiles` | ✓ | ✓ | ✓ | ✓ | ✓ (portal) | ✓ |
| `PickPhoto` library / camera | ✓ / ✓ | ✓ / Info.plist key | file / – | file / – | file / – | file / – |
| `Dictate` | ✓ | Info.plist keys | – | – | – | – |
| `Share` | share sheet | share sheet | reveal / clipboard | reveal / clipboard | reveal / clipboard | – |

"–" means the call returns `ErrNotSupported` (or, for `Haptic`, does nothing). The iOS
code is compiled in CI but hasn't run on a device.

## Links

```go
err := app.OpenURL(ctx, "https://gioui.org")
```

Opens the link in the app that handles it. Only `http`, `https`, `mailto`, `tel` and `sms`
links are opened, so a link in untrusted content (a model's reply) can't reach local files
or scripts. `RichText` and `Markdown` links open through it unless they have their own
`OnLink`. An in-app browser (Custom Tabs, `SFSafariViewController`) isn't available:
Custom Tabs needs an AndroidX library gogio doesn't bundle.

## Safe area and keyboard

`giode.SafeArea()` returns the insets covered by system UI (status bar, navigation bar,
notch, and on iOS the keyboard) in dp. By default the view is laid out inside them; with
`Config.EdgeToEdge` it covers the whole window and you pad with `SafeArea()`. On Android
the window shrinks for the keyboard, so content at the bottom stays above it either way.
`app.HideKeyboard()` hides the keyboard and clears focus.

## System theme and device

`Config.FollowSystemTheme` sets the theme mode from the platform's light/dark appearance
at startup and when it changes (checked every few seconds and whenever the window regains
focus); `Config.OnSystemThemeChange` is called on changes. `giode.SystemDarkMode()` reads
it once.

`giode.Device()` returns `TotalRAM`, `Class` (phone, tablet - a window whose shorter side is
at least 600dp - or desktop), `OS`, `OSVersion` and `Model`.

## Haptics and notifications

```go
giode.Haptic(giode.HapticSuccess) // HapticLight, HapticMedium, HapticHeavy, HapticSuccess, HapticWarning
err := app.Notify("Research finished", "3 sources, 2 minutes")
```

`OnLongPress` plays `HapticHeavy`. On Android, haptics follow the user's touch feedback
setting. `Notify` uses `gioui.org/x/notify`; on Android 13+ the system asks for permission
the first time, and an app built with `-targetsdk 33` or higher can't notify at all,
because gogio can't declare `POST_NOTIFICATIONS`. iOS notifications aren't implemented.

## Staying awake

`app.KeepScreenOn(true)` keeps the screen on while the app is shown, for on-device
generation the user is watching.

`release, err := app.KeepAlive("Researching")` keeps work running: desktops don't go to
idle sleep, and iOS grants its short background grace period (about 30 seconds).

On Android it starts a foreground service, which keeps the app running after it leaves
the screen and shows an ongoing notification with the reason; tapping it brings the app
back. Call it while the app is shown (Android doesn't let a backgrounded app start one),
and call `release` when the work is done. The service is in giode's helper, but Android
only runs a service the app's manifest declares:

```xml
<uses-permission android:name="android.permission.FOREGROUND_SERVICE"/>
<uses-permission android:name="android.permission.FOREGROUND_SERVICE_DATA_SYNC"/>
<application ...>
	<service android:name="org.giode.mobile.KeepAliveService"
		android:foregroundServiceType="dataSync" android:exported="false"/>
</application>
```

The [giode CLI](https://github.com/mehran9675/giode-cli) adds them when `giode.toml` has
`keep_alive = true` under `[android]`; gogio alone can't. Without the declaration
`KeepAlive` returns an error wrapping `ErrNotSupported`; keep the app open with
`KeepScreenOn` and save progress in `OnBackground` as a fallback. Android 15
limits `dataSync` services to six hours a day, after which the service stops.

## Secure storage

```go
giode.SecureStore.Service = "com.example.pocketagent" // set once, before use
err := giode.SecureStore.Set("anthropic.key", []byte(key))
key, err := giode.SecureStore.Get("anthropic.key") // ErrNotFound if unset
err = giode.SecureStore.Delete("anthropic.key")
```

Values live in the platform's credential store: the Android Keystore (AES-GCM with a key
that never leaves the device), the Keychain (on iOS accessible after first unlock, this
device only), DPAPI-encrypted files on Windows, and the Secret Service (GNOME Keyring,
KWallet) on Linux. Calls can block - a Linux keyring may prompt to unlock - so keep them
off the UI goroutine.

## Picking files and photos

```go
files, err := app.PickFiles(ctx, giode.PickOptions{Types: []string{"image/*", "application/pdf"}, Multiple: true})
for _, f := range files {
	fmt.Println(f.Name, f.MIME, f.Size) // "Q3 report.pdf", "application/pdf", 182345
}

photo, err := app.PickPhoto(ctx, giode.PhotoOptions{Source: giode.PhotoLibrary, MaxSide: 2048})
```

`PickedFile` is an `io.ReadCloser` with `Name`, `MIME` (from the extension, or sniffed from
the content) and `Size` (-1 when unknown). `PickPhoto` returns an upright JPEG scaled to
`MaxSide`: Android's photo picker (Android 13+, the image picker before) or camera, and
PHPicker on iOS 14+, decode HEIC on the device, so HEIC never reaches the app. On desktops
it picks an image file (JPEG, PNG, GIF, WebP). The iOS camera needs
`NSCameraUsageDescription` in Info.plist (see [Info.plist keys on iOS](#infoplist-keys-on-ios));
without it `PickPhoto` returns an error wrapping `ErrNotSupported` rather than letting iOS
end the app. Like `OpenFile`, these block: call them from a goroutine.

Reading an image from the clipboard isn't supported yet.

## Dictation

```go
results, err := app.Dictate(ctx, "en-US")
for r := range results {
	composer.SetText(r.Text) // the whole transcript so far
	if r.Final { break }
}
```

On Android it uses `SpeechRecognizer`, asking for the microphone permission the first
time; import `gioui.org/app/permission/microphone`. On iOS it uses `SFSpeechRecognizer`
with the microphone, asking for both permissions the first time, and ends after a two
second pause, as Android does. It needs `NSMicrophoneUsageDescription` and
`NSSpeechRecognitionUsageDescription` in Info.plist; without them `Dictate` returns an
error wrapping `ErrNotSupported`.

## Info.plist keys on iOS

iOS ends an app that uses the camera, microphone or speech recognition without a line in
its Info.plist explaining why, shown in the permission prompt. With the
[giode CLI](https://github.com/mehran9675/giode-cli), list them in `giode.toml`:

```toml
[ios.usage]
camera = "To attach photos you take."
microphone = "To dictate messages."
speech = "To dictate messages."
```

gogio writes Info.plist itself and has no option to add keys, so with gogio add them to
the built app before installing it:

```sh
gogio -target ios -arch amd64 -appid com.example.app -o App.app .
plutil -insert NSCameraUsageDescription -string "To attach photos you take." App.app/Info.plist
plutil -insert NSMicrophoneUsageDescription -string "To dictate messages." App.app/Info.plist
plutil -insert NSSpeechRecognitionUsageDescription -string "To dictate messages." App.app/Info.plist
codesign -f -s - App.app   # re-sign after editing; use your identity for a device build
```

For a device `.ipa`, unzip it, edit `Payload/App.app/Info.plist` the same way, and
re-sign the app with your identity and its entitlements before zipping it again.

## Links that open the app

`Config.OnOpenURL` receives the links that open the app, or reach it while it runs, on the
UI goroutine:

```go
giode.Config{
	OnOpenURL: func(u *url.URL) {
		// myapp://chat/42 or https://example.com/app/chat/42
		router.Go(u.Path)
	},
}
```

- **URL schemes** (`myapp://...`): list them in `giode.toml` as `app.schemes`. They work on
  Android, iOS, macOS, Windows and Linux (from the `.desktop` file the CLI writes).
- **Universal links and Android App Links** (`https://example.com/app/...`): list them as
  `app.links`; `giode links` writes the two files your site must serve (see the CLI's
  README). iOS and macOS need the Associated Domains capability in the provisioning profile.

With a `SingleInstanceLock`, a link that starts a second desktop instance is handed to the
running one.

## Updating desktop apps

The `update` package updates a desktop app from a feed you host: it checks for a newer
version, downloads it, verifies it, replaces the app and restarts it.

```go
import "github.com/mehran9675/giode/update"

u := update.New() // feed URL, public key and version from giode.toml, via the giode CLI
if rel, err := u.Check(ctx); err == nil && rel != nil {
	// Offer rel.Version and rel.Notes; then:
	dl, err := u.Download(ctx, rel, func(done, total int64) { progress.Set(done, total) })
	if err == nil {
		err = dl.Install() // replaces the app and restarts it; doesn't return on success
	}
}
```

Set it up with the CLI: `giode update-key` makes an Ed25519 key pair (keep the private key
out of the repository; it signs releases), and `giode release -base-url https://...` writes
`dist/update/` - the signed `feed.json`, `feed.json.sig` and a download per platform - to
upload to that URL. The app accepts only a feed signed with its key and downloads whose
SHA-256 the feed lists.

| Platform | What the feed offers | Installing |
| --- | --- | --- |
| Windows | the program (`exe`) or the setup program (`installer`) | renames the running program aside and puts the new one in place; or runs the setup program silently once the app exits |
| macOS | the zipped universal app (`app.zip`) | swaps the `.app` bundle |
| Linux | the program (`binary`) or AppImage (`appimage`) | replaces the file |

The app restarts with the same arguments; the old version is deleted on the next start.
Apps installed by a package manager or store (`.deb`, `.rpm`, Flatpak, Snap, MSIX, the Mac
App Store) are updated by it: `Check` returns `update.ErrManaged`.

## Crash reports

The `crash` package records the app's crashes and hands them back on the next start, on
every platform:

```go
import "github.com/mehran9675/giode/crash"

func main() {
	crash.Init(crash.Options{Version: "1.4.0"})
	if reports := crash.Reports(); len(reports) > 0 {
		// Offer to send them; then:
		for _, r := range reports {
			if r.Send(ctx, "https://example.com/api/crashes") == nil {
				r.Delete()
			}
		}
	}
	...
}
```

A crash is an unrecovered panic or a fatal error on any goroutine, including background
ones: the Go runtime writes the message and every goroutine's stack to the crash file as
the app dies (`runtime/debug.SetCrashOutput`). A `Report` has the message, the full trace,
the app and Go versions and the platform; `Send` posts it as JSON, and the newest ten are
kept. Crashes in Java (Android) or Objective-C (Apple) threads aren't Go's and aren't
recorded.

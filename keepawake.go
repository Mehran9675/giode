package giode

import "context"

// KeepScreenOn keeps the screen from dimming and locking while on is
// set and the app is shown - for on-device generation the user is
// watching. On desktops it holds off the screen saver and display
// sleep.
func (a *App) KeepScreenOn(on bool) error {
	return a.keepScreenOn(context.Background(), on)
}

// KeepAlive asks the system to keep the app running work it started -
// an agent run, a long download - until release is called:
//
//   - Desktops: the computer doesn't go to sleep while idle.
//   - iOS: the app gets the short grace period iOS allows after it
//     leaves the screen (about 30 seconds), via beginBackgroundTask.
//   - Android: not supported. Keeping work running once the app is in
//     the background takes a foreground service, which must be declared
//     in the app manifest, and gogio has no way to add one; Android
//     freezes the app soon after it leaves the screen. KeepAlive
//     returns ErrNotSupported there, so the app can tell the user to
//     keep it open (with KeepScreenOn) or save its progress in
//     OnBackground.
//
// reason is shown where the platform shows one (Linux's inhibitor
// list, macOS's power assertions).
func (a *App) KeepAlive(reason string) (release func(), err error) {
	return a.keepAlive(reason)
}

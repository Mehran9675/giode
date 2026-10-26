package giode

// HapticKind is a kind of haptic feedback.
type HapticKind int

const (
	// HapticLight is a light tap, for selections and key presses.
	HapticLight HapticKind = iota
	// HapticMedium is a firmer tap, for confirming a control.
	HapticMedium
	// HapticHeavy is the long-press feedback.
	HapticHeavy
	// HapticSuccess confirms a completed action (send, approve).
	HapticSuccess
	// HapticWarning flags a rejected or risky action.
	HapticWarning
)

// androidHaptic maps kinds to android.view.HapticFeedbackConstants.
var androidHaptic = map[HapticKind]int{
	HapticLight:   3,  // KEYBOARD_TAP
	HapticMedium:  1,  // VIRTUAL_KEY
	HapticHeavy:   0,  // LONG_PRESS
	HapticSuccess: 16, // CONFIRM
	HapticWarning: 17, // REJECT
}

// Haptic plays haptic feedback on the device: through the view's
// haptic feedback on Android (which follows the user's touch feedback
// setting) and the feedback generators on iOS. It returns at once and
// does nothing on platforms without haptics.
func (a *App) Haptic(k HapticKind) {
	go a.haptic(k)
}

// Haptic plays haptic feedback through the most recently created App;
// see App.Haptic.
func Haptic(k HapticKind) {
	if a := currentApp.Load(); a != nil {
		a.Haptic(k)
	}
}

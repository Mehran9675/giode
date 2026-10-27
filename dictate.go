package giode

import "context"

// Transcript is one result of Dictate: the text recognized so far,
// replaced by each later result, until the Final one.
type Transcript struct {
	Text  string
	Final bool
	Err   error
}

// Dictate listens through the microphone and transcribes speech, for a
// microphone button in a composer. Results arrive on the channel as the
// user speaks - each holds the whole transcript so far - and the
// channel closes after the Final one, which comes when the user stops
// speaking. Cancelling ctx stops listening early; a final result still
// follows.
//
// It uses the platform's speech recognizer, asking for permission the
// first time: Android's SpeechRecognizer (the app must import
// gioui.org/app/permission/microphone so the permission is in its
// manifest), and SFSpeechRecognizer on iOS (the app's Info.plist needs
// NSMicrophoneUsageDescription and NSSpeechRecognitionUsageDescription,
// or it returns an error wrapping ErrNotSupported; see
// docs/platform.md). Elsewhere it returns ErrNotSupported. lang is a
// BCP 47 tag ("en-US"), or "" for the device's language.
func (a *App) Dictate(ctx context.Context, lang string) (<-chan Transcript, error) {
	return a.dictate(ctx, lang)
}

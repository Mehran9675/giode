# Android library stubs

Compile-only stand-ins for the parts of Google's Android libraries that
giode's Java code calls: Firebase Messaging (push), Play Billing (store)
and Play In-App Review (review). The `go:generate` lines of those
packages compile their Java against these, so a checkout can rebuild the
jars without Gradle or Google's Maven repository. They are never
packaged: apps get the real libraries, which the giode CLI resolves from
each package's `android_dependencies.txt`.

Signatures must match the pinned versions exactly (Firebase Messaging
24.1.0, Play Billing 7.1.1, Play Review 2.0.2): a mismatch builds, then
fails at run time with NoSuchMethodError.

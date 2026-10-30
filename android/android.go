// Package android calls Java from Go on Android: static methods of the
// app's own classes (from a jar in one of its packages, which giode
// build adds to the APK) and of the Android SDK, through the app's
// class loader.
//
//	ok, err := android.CallBool("com/example/usb/Usb", "hasPermission",
//		"(Landroid/content/Context;Ljava/lang/String;)Z", android.Context(), name)
//
// Signatures are JNI type descriptors. Java runs a function with the
// JNI environment for anything the Call helpers don't cover (instance
// methods, fields, arrays), using git.wow.st/gmp/jni.
//
// Java calling Go goes through //export Java_<class>_<method> functions
// in the app; GoString and GoBytes convert their arguments.
//
// On other platforms the Call helpers return errors.ErrUnsupported and
// Available reports false.
package android

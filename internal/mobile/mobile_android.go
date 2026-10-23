package mobile

/*
#include <jni.h>
*/
import "C"

import (
	"git.wow.st/gmp/jni"
)

//go:generate javac --release 8 -classpath $ANDROID_HOME/platforms/android-16/android.jar -d $TEMP/giode_mobile/classes mobile_android.java
//go:generate jar cf mobile_android.jar -C $TEMP/giode_mobile/classes .

const helperClass = "org/giode/mobile/mobile_android"

// call loads the helper class through the app's class loader and runs
// f with it on the current thread, attached to the JVM.
func call(vm, ctx uintptr, f func(env jni.Env, cls jni.Class) error) error {
	return jni.Do(jni.JVMFor(vm), func(env jni.Env) error {
		cls, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(ctx)), helperClass)
		if err != nil {
			return err
		}
		return f(env, cls)
	})
}

// WatchLifecycle starts reporting background and foreground changes to
// the handlers set with SetLifecycleHandlers. ctx is the Android
// application context. Calling it again does nothing.
func WatchLifecycle(vm, ctx uintptr) error {
	return call(vm, ctx, func(env jni.Env, cls jni.Class) error {
		m := jni.GetStaticMethodID(env, cls, "watchLifecycle", "(Landroid/content/Context;)V")
		return jni.CallStaticVoidMethod(env, cls, m, jni.Value(ctx))
	})
}

// Publish copies the file at path into the shared Downloads collection
// and returns its content URI, readable by the app it's shared with.
// It copies the whole file, so call it off the main thread. Android 10
// or later.
func Publish(vm, ctx uintptr, path, mime string) (string, error) {
	var uri string
	err := call(vm, ctx, func(env jni.Env, cls jni.Class) error {
		m := jni.GetStaticMethodID(env, cls, "publish", "(Landroid/content/Context;Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;")
		obj, err := jni.CallStaticObjectMethod(env, cls, m,
			jni.Value(ctx),
			jni.Value(jni.JavaString(env, path)),
			jni.Value(jni.JavaString(env, mime)),
		)
		if err != nil {
			return err
		}
		uri = jni.GoString(env, jni.String(obj))
		return nil
	})
	return uri, err
}

// ShareText opens the share sheet for text from view, the GioView.
// Call it on the main thread.
func ShareText(vm, ctx, view uintptr, text, subject, title string) error {
	return call(vm, ctx, func(env jni.Env, cls jni.Class) error {
		m := jni.GetStaticMethodID(env, cls, "shareText", "(Landroid/view/View;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V")
		return jni.CallStaticVoidMethod(env, cls, m,
			jni.Value(view),
			jni.Value(jni.JavaString(env, text)),
			jni.Value(jni.JavaString(env, subject)),
			jni.Value(jni.JavaString(env, title)),
		)
	})
}

// ShareURI opens the share sheet for a URI returned by Publish. Call
// it on the main thread.
func ShareURI(vm, ctx, view uintptr, uri, mime, title string) error {
	return call(vm, ctx, func(env jni.Env, cls jni.Class) error {
		m := jni.GetStaticMethodID(env, cls, "shareURI", "(Landroid/view/View;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V")
		return jni.CallStaticVoidMethod(env, cls, m,
			jni.Value(view),
			jni.Value(jni.JavaString(env, uri)),
			jni.Value(jni.JavaString(env, mime)),
			jni.Value(jni.JavaString(env, title)),
		)
	})
}

//export Java_org_giode_mobile_mobile_1android_onBackground
func Java_org_giode_mobile_mobile_1android_onBackground(env *C.JNIEnv, cls C.jclass) {
	notify(&onBackground)
}

//export Java_org_giode_mobile_mobile_1android_onForeground
func Java_org_giode_mobile_mobile_1android_onForeground(env *C.JNIEnv, cls C.jclass) {
	notify(&onForeground)
}

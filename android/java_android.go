package android

/*
#include <jni.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"gioui.org/app"
	"git.wow.st/gmp/jni"
)

// Available reports whether Java can be called: true on Android.
func Available() bool { return true }

// Java runs f with a JNI environment attached to the current thread and
// the class (in slash form: "com/example/usb/Usb"), loaded through the
// app's class loader, so classes from the app's own jars resolve.
func Java(class string, f func(env jni.Env, cls jni.Class) error) error {
	return jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		cls, err := jni.LoadClass(env, jni.ClassLoaderFor(env, Context()), class)
		if err != nil {
			return err
		}
		return f(env, cls)
	})
}

// Context returns the application's android.content.Context, a global
// reference valid for the app's lifetime.
func Context() jni.Object {
	return jni.Object(app.AppContext())
}

// Args converts Go values to JNI arguments: strings to Java strings,
// byte slices to byte[], ints, int32s and int64s, bools, float32s and
// float64s, and uintptr, jni.Object or jni.Value object references.
func Args(env jni.Env, args ...any) []jni.Value {
	vs := make([]jni.Value, len(args))
	for i, a := range args {
		switch a := a.(type) {
		case string:
			vs[i] = jni.Value(jni.JavaString(env, a))
		case []byte:
			vs[i] = jni.Value(jni.NewByteArray(env, a))
		case int:
			vs[i] = jni.Value(a)
		case int32:
			vs[i] = jni.Value(a)
		case int64:
			vs[i] = jni.Value(a)
		case bool:
			if a {
				vs[i] = 1
			}
		case float32:
			vs[i] = jni.Value(*(*uint32)(unsafe.Pointer(&a)))
		case float64:
			vs[i] = jni.Value(*(*uint64)(unsafe.Pointer(&a)))
		case uintptr:
			vs[i] = jni.Value(a)
		case jni.Object:
			vs[i] = jni.Value(a)
		case jni.Value:
			vs[i] = a
		default:
			panic(fmt.Sprintf("android: can't pass %T to Java", a))
		}
	}
	return vs
}

func static(env jni.Env, cls jni.Class, method, sig string) jni.MethodID {
	return jni.GetStaticMethodID(env, cls, method, sig)
}

// CallVoid calls a static void method: CallVoid("com/example/Usb",
// "open", "(Landroid/content/Context;I)V", android.Context(), id).
func CallVoid(class, method, sig string, args ...any) error {
	return Java(class, func(env jni.Env, cls jni.Class) error {
		return jni.CallStaticVoidMethod(env, cls, static(env, cls, method, sig), Args(env, args...)...)
	})
}

// CallString calls a static method returning a String; null is "".
func CallString(class, method, sig string, args ...any) (string, error) {
	var s string
	err := Java(class, func(env jni.Env, cls jni.Class) error {
		obj, err := jni.CallStaticObjectMethod(env, cls, static(env, cls, method, sig), Args(env, args...)...)
		if err != nil {
			return err
		}
		s = jni.GoString(env, jni.String(obj))
		return nil
	})
	return s, err
}

// CallBool calls a static method returning a boolean.
func CallBool(class, method, sig string, args ...any) (bool, error) {
	var b bool
	err := Java(class, func(env jni.Env, cls jni.Class) error {
		var err error
		b, err = jni.CallStaticBooleanMethod(env, cls, static(env, cls, method, sig), Args(env, args...)...)
		return err
	})
	return b, err
}

// CallInt calls a static method returning an int.
func CallInt(class, method, sig string, args ...any) (int, error) {
	var n int
	err := Java(class, func(env jni.Env, cls jni.Class) error {
		var err error
		n, err = jni.CallStaticIntMethod(env, cls, static(env, cls, method, sig), Args(env, args...)...)
		return err
	})
	return n, err
}

// CallBytes calls a static method returning a byte[]; null is nil.
func CallBytes(class, method, sig string, args ...any) ([]byte, error) {
	var b []byte
	err := Java(class, func(env jni.Env, cls jni.Class) error {
		obj, err := jni.CallStaticObjectMethod(env, cls, static(env, cls, method, sig), Args(env, args...)...)
		if err != nil || obj == 0 {
			return err
		}
		b = jni.GetByteArrayElements(env, jni.ByteArray(obj))
		return nil
	})
	return b, err
}

// GoString converts a Java string passed to a native method (a
// //export Java_... function) to Go; null is "".
func GoString(env unsafe.Pointer, s uintptr) string {
	if s == 0 {
		return ""
	}
	return jni.GoString(jni.EnvFor(uintptr(env)), jni.String(s))
}

// GoBytes converts a Java byte[] passed to a native method to Go; null
// is nil.
func GoBytes(env unsafe.Pointer, b uintptr) []byte {
	if b == 0 {
		return nil
	}
	return jni.GetByteArrayElements(jni.EnvFor(uintptr(env)), jni.ByteArray(b))
}

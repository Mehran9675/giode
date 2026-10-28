package bridge

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

// Java calls static methods of class (org/giode/...), loaded through
// the app's class loader, on the current thread.
func Java(class string, f func(env jni.Env, cls jni.Class) error) error {
	return jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		cls, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(app.AppContext())), class)
		if err != nil {
			return err
		}
		return f(env, cls)
	})
}

// Args converts Go values to JNI arguments: strings to Java strings,
// ints, bools, and uintptr or jni.Value object references.
func Args(env jni.Env, args ...any) []jni.Value {
	vs := make([]jni.Value, len(args))
	for i, a := range args {
		switch a := a.(type) {
		case string:
			vs[i] = jni.Value(jni.JavaString(env, a))
		case int:
			vs[i] = jni.Value(a)
		case bool:
			if a {
				vs[i] = 1
			}
		case uintptr:
			vs[i] = jni.Value(a)
		case jni.Value:
			vs[i] = a
		default:
			panic(fmt.Sprintf("bridge: can't pass %T to Java", a))
		}
	}
	return vs
}

// CallVoid calls a static void method.
func CallVoid(class, method, sig string, args ...any) error {
	return Java(class, func(env jni.Env, cls jni.Class) error {
		return jni.CallStaticVoidMethod(env, cls, jni.GetStaticMethodID(env, cls, method, sig), Args(env, args...)...)
	})
}

// CallString calls a static method returning a String; null is "".
func CallString(class, method, sig string, args ...any) (string, error) {
	var s string
	err := Java(class, func(env jni.Env, cls jni.Class) error {
		obj, err := jni.CallStaticObjectMethod(env, cls, jni.GetStaticMethodID(env, cls, method, sig), Args(env, args...)...)
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
		b, err = jni.CallStaticBooleanMethod(env, cls, jni.GetStaticMethodID(env, cls, method, sig), Args(env, args...)...)
		return err
	})
	return b, err
}

// GoString converts a Java string passed to a native method; null is
// "".
func GoString(env unsafe.Pointer, s uintptr) string {
	if s == 0 {
		return ""
	}
	return jni.GoString(jni.EnvFor(uintptr(env)), jni.String(s))
}

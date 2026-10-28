//go:build darwin && !ios

@import AppKit;
#include <objc/runtime.h>
#include "_cgo_export.h"

// The status item's owner in fyne.io/systray.
extern id owner;

void giodeTrayStartMain(void) {
	dispatch_async(dispatch_get_main_queue(), ^{
		giodeTrayStart();
	});
}

void giodeTrayHideItem(void) {
	dispatch_async(dispatch_get_main_queue(), ^{
		if (owner == nil) {
			return;
		}
		NSStatusItem *item = [owner valueForKey:@"statusItem"];
		item.visible = NO;
	});
}

static NSWindow *windowOf(uintptr_t view) {
	NSView *v = (__bridge NSView *)(void *)view;
	return v.window;
}

void giodeHideWindow(uintptr_t view) {
	dispatch_async(dispatch_get_main_queue(), ^{
		[windowOf(view) orderOut:nil];
	});
}

void giodeShowWindow(uintptr_t view) {
	dispatch_async(dispatch_get_main_queue(), ^{
		NSWindow *w = windowOf(view);
		if (w.miniaturized) {
			[w deminiaturize:nil];
		}
		[NSApp activateIgnoringOtherApps:YES];
		[w makeKeyAndOrderFront:nil];
	});
}

static BOOL giodeWindowShouldClose(id self, SEL _cmd, NSWindow *sender) {
	if (giodeShouldHideOnClose()) {
		[sender orderOut:nil];
		giodeHiddenOnClose();
		return NO;
	}
	return YES;
}

// giodeInterceptClose gives Gio's window delegate a windowShouldClose:
// that hides the window instead, for CloseToTray.
void giodeInterceptClose(uintptr_t view) {
	dispatch_async(dispatch_get_main_queue(), ^{
		id del = windowOf(view).delegate;
		if (del == nil) {
			return;
		}
		struct objc_method_description d = protocol_getMethodDescription(@protocol(NSWindowDelegate), @selector(windowShouldClose:), NO, YES);
		class_replaceMethod([del class], @selector(windowShouldClose:), (IMP)giodeWindowShouldClose, d.types);
	});
}

static BOOL giodeShouldHandleReopen(id self, SEL _cmd, NSApplication *app, BOOL visible) {
	giodeReopen();
	return YES;
}

// giodeHandleReopen shows a hidden window when the Dock icon is
// clicked.
void giodeHandleReopen(void) {
	dispatch_async(dispatch_get_main_queue(), ^{
		id del = [NSApp delegate];
		if (del == nil) {
			return;
		}
		SEL sel = @selector(applicationShouldHandleReopen:hasVisibleWindows:);
		struct objc_method_description d = protocol_getMethodDescription(@protocol(NSApplicationDelegate), sel, NO, YES);
		class_addMethod([del class], sel, (IMP)giodeShouldHandleReopen, d.types);
	});
}

// giodeMenuBarApp removes the Dock icon and app menu. Gio's delegate
// makes the app a regular one when it finishes launching, so the
// policy is set again after that.
void giodeMenuBarApp(void) {
	dispatch_async(dispatch_get_main_queue(), ^{
		[NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
	});
	[[NSNotificationCenter defaultCenter] addObserverForName:NSApplicationDidFinishLaunchingNotification
	                                                  object:nil
	                                                   queue:[NSOperationQueue mainQueue]
	                                              usingBlock:^(NSNotification *note) {
		[NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
	}];
}

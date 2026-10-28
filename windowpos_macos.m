//go:build darwin && !ios

@import AppKit;
#include <os/lock.h>

// Window placement, read on the main thread whenever the window moves
// or resizes, or the screens change, so the UI goroutine can read it
// without waiting for the main thread - which waits for the UI
// goroutine while a frame is drawn. Rectangles are top-left based, from
// the primary screen's top left, as {x, y, width, height}.

static os_unfair_lock placementLock = OS_UNFAIR_LOCK_INIT;
static double placement[13]; // frame, screen, visible frame, scale
static int placementValid;
static BOOL placementWatched;

static NSWindow *placementWindow(uintptr_t view) {
	NSView *v = (__bridge NSView *)(void *)view;
	return v.window;
}

static void flipRect(NSRect r, CGFloat h, double *out) {
	out[0] = r.origin.x;
	out[1] = h - (r.origin.y + r.size.height);
	out[2] = r.size.width;
	out[3] = r.size.height;
}

static void updatePlacement(NSWindow *w) {
	if (w == nil) {
		return;
	}
	NSScreen *primary = NSScreen.screens.firstObject;
	NSScreen *s = w.screen ?: NSScreen.mainScreen;
	if (primary == nil || s == nil) {
		return;
	}
	CGFloat h = primary.frame.size.height;
	double p[13];
	flipRect(w.frame, h, p);
	flipRect(s.frame, h, p + 4);
	flipRect(s.visibleFrame, h, p + 8);
	p[12] = s.backingScaleFactor;
	os_unfair_lock_lock(&placementLock);
	memcpy(placement, p, sizeof(p));
	placementValid = 1;
	os_unfair_lock_unlock(&placementLock);
}

void giodeWatchPlacement(uintptr_t view) {
	dispatch_async(dispatch_get_main_queue(), ^{
		NSWindow *w = placementWindow(view);
		updatePlacement(w);
		if (placementWatched || w == nil) {
			return;
		}
		placementWatched = YES;
		NSNotificationCenter *nc = NSNotificationCenter.defaultCenter;
		void (^update)(NSNotification *) = ^(NSNotification *n) {
			updatePlacement(w);
		};
		for (NSNotificationName name in @[NSWindowDidMoveNotification, NSWindowDidResizeNotification, NSWindowDidChangeScreenNotification]) {
			[nc addObserverForName:name object:w queue:nil usingBlock:update];
		}
		[nc addObserverForName:NSApplicationDidChangeScreenParametersNotification object:nil queue:nil usingBlock:update];
	});
}

int giodePlacement(double *out) {
	os_unfair_lock_lock(&placementLock);
	memcpy(out, placement, sizeof(placement));
	int ok = placementValid;
	os_unfair_lock_unlock(&placementLock);
	return ok;
}

void giodeSetPosition(uintptr_t view, double x, double y) {
	dispatch_async(dispatch_get_main_queue(), ^{
		NSWindow *w = placementWindow(view);
		NSScreen *primary = NSScreen.screens.firstObject;
		if (w == nil || primary == nil) {
			return;
		}
		[w setFrameTopLeftPoint:NSMakePoint(x, primary.frame.size.height - y)];
		updatePlacement(w);
	});
}

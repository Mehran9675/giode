#import <UIKit/UIKit.h>
#import <PhotosUI/PhotosUI.h>
#include <stdint.h>
#include "_cgo_export.h"

static void onMain(dispatch_block_t b) {
	if ([NSThread isMainThread]) {
		b();
	} else {
		dispatch_async(dispatch_get_main_queue(), b);
	}
}

void giodeOpenURL(const char *url) {
	NSURL *u = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
	onMain(^{
		[[UIApplication sharedApplication] openURL:u options:@{} completionHandler:nil];
	});
}

void giodeShare(CFTypeRef vcRef, const char *text, const char *path, const char *subject) {
	UIViewController *vc = (__bridge UIViewController *)vcRef;
	NSMutableArray *items = [NSMutableArray array];
	if (text != NULL) {
		[items addObject:[NSString stringWithUTF8String:text]];
	}
	if (path != NULL) {
		[items addObject:[NSURL fileURLWithPath:[NSString stringWithUTF8String:path]]];
	}
	NSString *subj = subject != NULL ? [NSString stringWithUTF8String:subject] : nil;
	onMain(^{
		UIActivityViewController *a = [[UIActivityViewController alloc] initWithActivityItems:items applicationActivities:nil];
		if (subj.length > 0) {
			[a setValue:subj forKey:@"subject"];
		}
		// On iPad the sheet is a popover and needs an anchor: the
		// middle of the window.
		UIPopoverPresentationController *pop = a.popoverPresentationController;
		if (pop != nil) {
			pop.sourceView = vc.view;
			CGRect b = vc.view.bounds;
			pop.sourceRect = CGRectMake(CGRectGetMidX(b), CGRectGetMidY(b), 1, 1);
			pop.permittedArrowDirections = 0;
		}
		[vc presentViewController:a animated:YES completion:nil];
	});
}

void giodeHaptic(int kind) {
	onMain(^{
		switch (kind) {
		case 0:
			[[[UIImpactFeedbackGenerator alloc] initWithStyle:UIImpactFeedbackStyleLight] impactOccurred];
			break;
		case 1:
			[[[UIImpactFeedbackGenerator alloc] initWithStyle:UIImpactFeedbackStyleMedium] impactOccurred];
			break;
		case 2:
			[[[UIImpactFeedbackGenerator alloc] initWithStyle:UIImpactFeedbackStyleHeavy] impactOccurred];
			break;
		case 3:
			[[[UINotificationFeedbackGenerator alloc] init] notificationOccurred:UINotificationFeedbackTypeSuccess];
			break;
		default:
			[[[UINotificationFeedbackGenerator alloc] init] notificationOccurred:UINotificationFeedbackTypeWarning];
			break;
		}
	});
}

void giodeIdleTimerDisabled(int disabled) {
	onMain(^{
		[UIApplication sharedApplication].idleTimerDisabled = disabled != 0;
	});
}

uintptr_t giodeBeginBackgroundTask(const char *name) {
	__block UIBackgroundTaskIdentifier task = UIBackgroundTaskInvalid;
	NSString *n = [NSString stringWithUTF8String:name];
	task = [[UIApplication sharedApplication] beginBackgroundTaskWithName:n expirationHandler:^{
		// Time's up: iOS suspends the app once the task ends.
		[[UIApplication sharedApplication] endBackgroundTask:task];
		task = UIBackgroundTaskInvalid;
	}];
	return (uintptr_t)task;
}

void giodeEndBackgroundTask(uintptr_t task) {
	if ((UIBackgroundTaskIdentifier)task != UIBackgroundTaskInvalid) {
		[[UIApplication sharedApplication] endBackgroundTask:(UIBackgroundTaskIdentifier)task];
	}
}

void giodeWatchLifecycle(void) {
	static dispatch_once_t once;
	dispatch_once(&once, ^{
		NSNotificationCenter *nc = [NSNotificationCenter defaultCenter];
		[nc addObserverForName:UIApplicationDidEnterBackgroundNotification object:nil queue:nil usingBlock:^(NSNotification *n) {
			giodeOnBackground();
		}];
		[nc addObserverForName:UIApplicationWillEnterForegroundNotification object:nil queue:nil usingBlock:^(NSNotification *n) {
			giodeOnForeground();
		}];
	});
}

API_AVAILABLE(ios(14))
@interface GiodePhotoDelegate : NSObject <PHPickerViewControllerDelegate>
@property int requestID;
@property int maxSide;
@end

// pending keeps delegates alive until their picker finishes.
static NSMutableSet *pending;

static NSData *giodeJPEG(UIImage *img, int maxSide) {
	CGSize size = img.size;
	CGFloat longest = MAX(size.width, size.height);
	if (maxSide > 0 && longest > maxSide) {
		CGFloat s = maxSide / longest;
		size = CGSizeMake(floor(size.width * s), floor(size.height * s));
	}
	UIGraphicsImageRendererFormat *f = [UIGraphicsImageRendererFormat defaultFormat];
	f.scale = 1;
	UIGraphicsImageRenderer *r = [[UIGraphicsImageRenderer alloc] initWithSize:size format:f];
	// Drawing applies the image's orientation, so the JPEG is upright.
	UIImage *out = [r imageWithActions:^(UIGraphicsImageRendererContext *ctx) {
		[img drawInRect:CGRectMake(0, 0, size.width, size.height)];
	}];
	return UIImageJPEGRepresentation(out, 0.9);
}

@implementation GiodePhotoDelegate
- (void)picker:(PHPickerViewController *)picker didFinishPicking:(NSArray<PHPickerResult *> *)results API_AVAILABLE(ios(14)) {
	[picker dismissViewControllerAnimated:YES completion:nil];
	int rid = self.requestID;
	int maxSide = self.maxSide;
	[pending removeObject:self];
	if (results.count == 0) {
		giodePhotoDone(rid, NULL, 0, NULL, NULL);
		return;
	}
	NSItemProvider *p = results[0].itemProvider;
	NSString *name = p.suggestedName ?: @"photo";
	if (![p canLoadObjectOfClass:[UIImage class]]) {
		giodePhotoDone(rid, NULL, 0, NULL, (char *)"the picked item isn't an image");
		return;
	}
	// UIImage decodes HEIC, so it never reaches the app.
	[p loadObjectOfClass:[UIImage class] completionHandler:^(id<NSItemProviderReading> obj, NSError *err) {
		if (err != nil || obj == nil) {
			NSString *msg = err != nil ? err.localizedDescription : @"can't load the image";
			giodePhotoDone(rid, NULL, 0, NULL, (char *)msg.UTF8String);
			return;
		}
		NSData *jpeg = giodeJPEG((UIImage *)obj, maxSide);
		giodePhotoDone(rid, (void *)jpeg.bytes, (int)jpeg.length, (char *)name.UTF8String, NULL);
	}];
}
@end

int giodePickPhoto(CFTypeRef vcRef, int rid, int maxSide) {
	if (@available(iOS 14, *)) {
		UIViewController *vc = (__bridge UIViewController *)vcRef;
		onMain(^{
			PHPickerConfiguration *cfg = [[PHPickerConfiguration alloc] init];
			cfg.filter = [PHPickerFilter imagesFilter];
			cfg.selectionLimit = 1;
			PHPickerViewController *picker = [[PHPickerViewController alloc] initWithConfiguration:cfg];
			GiodePhotoDelegate *d = [[GiodePhotoDelegate alloc] init];
			d.requestID = rid;
			d.maxSide = maxSide;
			if (pending == nil) {
				pending = [NSMutableSet set];
			}
			[pending addObject:d];
			picker.delegate = d;
			[vc presentViewController:picker animated:YES completion:nil];
		});
		return 1;
	}
	return 0;
}

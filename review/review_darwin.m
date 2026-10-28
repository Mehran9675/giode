// SPDX-License-Identifier: Unlicense OR MIT

#include <TargetConditionals.h>
#import <StoreKit/StoreKit.h>
#if TARGET_OS_IPHONE
#import <UIKit/UIKit.h>
#endif

void giodeRequestReview(void) {
	dispatch_async(dispatch_get_main_queue(), ^{
#if TARGET_OS_IPHONE
		if (@available(iOS 14.0, *)) {
			for (UIScene *s in [UIApplication sharedApplication].connectedScenes) {
				if (s.activationState == UISceneActivationStateForegroundActive && [s isKindOfClass:[UIWindowScene class]]) {
					[SKStoreReviewController requestReviewInScene:(UIWindowScene *)s];
					return;
				}
			}
		}
		[SKStoreReviewController requestReview];
#else
		[SKStoreReviewController requestReview];
#endif
	});
}

package com.google.android.play.core.review;

import com.google.android.gms.tasks.Task;

public interface ReviewManager {
	Task<ReviewInfo> requestReviewFlow();
	Task<Void> launchReviewFlow(android.app.Activity activity, ReviewInfo info);
}

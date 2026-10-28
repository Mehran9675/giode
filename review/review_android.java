// SPDX-License-Identifier: Unlicense OR MIT

package org.giode.review;

import android.app.Activity;
import android.content.Context;
import android.view.View;

import com.google.android.gms.tasks.OnCompleteListener;
import com.google.android.gms.tasks.Task;
import com.google.android.play.core.review.ReviewInfo;
import com.google.android.play.core.review.ReviewManager;
import com.google.android.play.core.review.ReviewManagerFactory;

// review_android asks Google Play's in-app review flow to show. Play
// decides whether it does (it has a quota), and doesn't say.
public final class review_android {
	static native void onDone(int id, String err);

	public static void request(View view, final int id) {
		Context ctx = view.getContext();
		if (!(ctx instanceof Activity)) {
			onDone(id, "no activity");
			return;
		}
		final Activity activity = (Activity) ctx;
		final ReviewManager m = ReviewManagerFactory.create(activity);
		m.requestReviewFlow().addOnCompleteListener(new OnCompleteListener<ReviewInfo>() {
			@Override public void onComplete(Task<ReviewInfo> task) {
				if (!task.isSuccessful()) {
					Exception e = task.getException();
					onDone(id, e != null ? e.toString() : "review flow unavailable");
					return;
				}
				m.launchReviewFlow(activity, task.getResult()).addOnCompleteListener(new OnCompleteListener<Void>() {
					@Override public void onComplete(Task<Void> t) {
						onDone(id, null);
					}
				});
			}
		});
	}
}

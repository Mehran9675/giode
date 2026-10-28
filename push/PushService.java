// SPDX-License-Identifier: Unlicense OR MIT

package org.giode.push;

import com.google.firebase.messaging.FirebaseMessagingService;
import com.google.firebase.messaging.RemoteMessage;

// PushService receives Firebase messages and token changes. giode/push's
// android_manifest.xml declares it.
public class PushService extends FirebaseMessagingService {
	@Override public void onNewToken(String token) {
		push_android.deliverToken(token);
	}

	@Override public void onMessageReceived(RemoteMessage message) {
		push_android.deliverMessage(message);
	}
}

// SPDX-License-Identifier: Unlicense OR MIT

package org.giode.push;

import android.app.Activity;
import android.app.Fragment;
import android.app.FragmentManager;
import android.content.Context;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.os.Build;
import android.os.Bundle;
import android.os.Handler;
import android.os.Looper;
import android.view.View;

import com.google.android.gms.tasks.OnCompleteListener;
import com.google.android.gms.tasks.Task;
import com.google.firebase.messaging.FirebaseMessaging;
import com.google.firebase.messaging.RemoteMessage;

import org.json.JSONObject;

import java.util.ArrayList;
import java.util.Map;

// push_android connects Firebase Cloud Messaging to giode/push. Tokens
// and messages that arrive before the Go side starts (the service can
// run without the app on screen, with no Go library loaded) wait in
// pending.
public final class push_android {
	private static boolean started;
	private static final ArrayList<String[]> pending = new ArrayList<String[]>();

	static native void onToken(String token);
	static native void onMessage(String json, boolean opened);
	static native void onTokenResult(int id, String token, String err);
	static native void onPermission(int id, boolean granted);

	// start is called by Go once its handlers are ready.
	public static void start() {
		ArrayList<String[]> queued;
		synchronized (pending) {
			started = true;
			queued = new ArrayList<String[]>(pending);
			pending.clear();
		}
		for (String[] e : queued) {
			if (e[0].equals("token")) {
				onToken(e[1]);
			} else {
				onMessage(e[1], false);
			}
		}
	}

	static void deliverToken(String token) {
		synchronized (pending) {
			if (!started) {
				pending.add(new String[]{"token", token});
				return;
			}
		}
		onToken(token);
	}

	static void deliverMessage(RemoteMessage m) {
		String json = messageJSON(m);
		synchronized (pending) {
			if (!started) {
				pending.add(new String[]{"message", json});
				return;
			}
		}
		onMessage(json, false);
	}

	private static String messageJSON(RemoteMessage m) {
		JSONObject o = new JSONObject();
		try {
			o.put("id", m.getMessageId());
			RemoteMessage.Notification n = m.getNotification();
			if (n != null) {
				o.put("title", n.getTitle());
				o.put("body", n.getBody());
			}
			JSONObject data = new JSONObject();
			for (Map.Entry<String, String> e : m.getData().entrySet()) {
				data.put(e.getKey(), e.getValue());
			}
			o.put("data", data);
		} catch (Exception e) {
		}
		return o.toString();
	}

	public static void getToken(final int id) {
		FirebaseMessaging.getInstance().getToken().addOnCompleteListener(new OnCompleteListener<String>() {
			@Override public void onComplete(Task<String> task) {
				if (task.isSuccessful()) {
					onTokenResult(id, task.getResult(), null);
				} else {
					Exception e = task.getException();
					onTokenResult(id, null, e != null ? e.toString() : "unknown error");
				}
			}
		});
	}

	public static void deleteToken() {
		FirebaseMessaging.getInstance().deleteToken();
	}

	// openedMessage returns the push message that launched the activity,
	// as JSON, or null. FCM starts the launcher activity with the
	// message's data as intent extras when a notification is tapped.
	public static String openedMessage(View view) {
		Context ctx = view.getContext();
		if (!(ctx instanceof Activity)) {
			return null;
		}
		Intent intent = ((Activity) ctx).getIntent();
		Bundle extras = intent != null ? intent.getExtras() : null;
		if (extras == null || !extras.containsKey("google.message_id")) {
			return null;
		}
		JSONObject o = new JSONObject();
		JSONObject data = new JSONObject();
		try {
			o.put("id", extras.getString("google.message_id"));
			for (String k : extras.keySet()) {
				if (k.startsWith("google.") || k.startsWith("gcm.") || k.equals("from") || k.equals("collapse_key")) {
					continue;
				}
				Object v = extras.get(k);
				if (v != null) {
					data.put(k, v.toString());
				}
			}
			o.put("data", data);
		} catch (Exception e) {
		}
		// Only once: the intent stays with the activity.
		extras.remove("google.message_id");
		intent.removeExtra("google.message_id");
		return o.toString();
	}

	public static boolean permitted(Context ctx) {
		if (Build.VERSION.SDK_INT < 33) {
			return true;
		}
		return ctx.checkSelfPermission("android.permission.POST_NOTIFICATIONS") == PackageManager.PERMISSION_GRANTED;
	}

	// requestPermission asks for POST_NOTIFICATIONS (Android 13 and
	// later) through a headless fragment, which receives the answer the
	// activity doesn't pass on.
	public static void requestPermission(final View view, final int id) {
		final Context ctx = view.getContext();
		if (permitted(ctx) || !(ctx instanceof Activity)) {
			onPermission(id, permitted(ctx));
			return;
		}
		new Handler(Looper.getMainLooper()).post(new Runnable() {
			@Override public void run() {
				FragmentManager fm = ((Activity) ctx).getFragmentManager();
				PermissionFragment f = new PermissionFragment();
				f.id = id;
				fm.beginTransaction().add(f, "giode-push-permission-" + id).commitAllowingStateLoss();
			}
		});
	}

	public static final class PermissionFragment extends Fragment {
		int id;
		private boolean asked;

		@Override public void onStart() {
			super.onStart();
			if (!asked) {
				asked = true;
				requestPermissions(new String[]{"android.permission.POST_NOTIFICATIONS"}, 1);
			}
		}

		@Override public void onRequestPermissionsResult(int code, String[] perms, int[] results) {
			boolean granted = results.length > 0 && results[0] == PackageManager.PERMISSION_GRANTED;
			getFragmentManager().beginTransaction().remove(this).commitAllowingStateLoss();
			onPermission(id, granted);
		}
	}
}

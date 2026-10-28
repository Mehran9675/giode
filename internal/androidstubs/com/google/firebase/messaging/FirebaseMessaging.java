package com.google.firebase.messaging;

import com.google.android.gms.tasks.Task;

public class FirebaseMessaging {
	public static synchronized FirebaseMessaging getInstance() { throw new UnsupportedOperationException(); }
	public Task<String> getToken() { throw new UnsupportedOperationException(); }
	public Task<Void> deleteToken() { throw new UnsupportedOperationException(); }
}

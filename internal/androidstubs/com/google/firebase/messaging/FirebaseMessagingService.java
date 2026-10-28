package com.google.firebase.messaging;

public class FirebaseMessagingService extends android.app.Service {
	public FirebaseMessagingService() {}
	public void onMessageReceived(RemoteMessage message) {}
	public void onNewToken(String token) {}
	public final android.os.IBinder onBind(android.content.Intent intent) { return null; }
}

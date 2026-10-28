package com.google.firebase.messaging;

import java.util.Map;

public final class RemoteMessage {
	public Map<String, String> getData() { throw new UnsupportedOperationException(); }
	public Notification getNotification() { throw new UnsupportedOperationException(); }
	public String getMessageId() { throw new UnsupportedOperationException(); }
	public String getFrom() { throw new UnsupportedOperationException(); }

	public static class Notification {
		public String getTitle() { throw new UnsupportedOperationException(); }
		public String getBody() { throw new UnsupportedOperationException(); }
	}
}

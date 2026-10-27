package org.giode.mobile;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.app.Service;
import android.content.ComponentName;
import android.content.Context;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.content.pm.ServiceInfo;
import android.os.Build;
import android.os.IBinder;

// KeepAliveService is the foreground service behind App.KeepAlive on
// Android: while it runs, Android keeps the app's process running (and
// unfrozen) after it leaves the screen, showing an ongoing
// notification. The app's manifest must declare it:
//
//	<uses-permission android:name="android.permission.FOREGROUND_SERVICE"/>
//	<uses-permission android:name="android.permission.FOREGROUND_SERVICE_DATA_SYNC"/>
//	<service android:name="org.giode.mobile.KeepAliveService"
//		android:foregroundServiceType="dataSync" android:exported="false"/>
public final class KeepAliveService extends Service {
	private static final String CHANNEL = "giode.keepalive";
	private static final int NOTIFICATION_ID = 0x6170;
	private static final String EXTRA_REASON = "reason";

	// declared reports whether the app's manifest declares the service
	// and the foreground service permission, so start can succeed.
	public static boolean declared(Context ctx) {
		try {
			ctx.getPackageManager().getServiceInfo(new ComponentName(ctx, KeepAliveService.class), 0);
		} catch (PackageManager.NameNotFoundException e) {
			return false;
		}
		if (Build.VERSION.SDK_INT >= 28 && ctx.checkSelfPermission("android.permission.FOREGROUND_SERVICE") != PackageManager.PERMISSION_GRANTED) {
			return false;
		}
		return true;
	}

	// start starts the service, or updates its notification text when
	// it's running. Android only lets an app start a foreground service
	// while it's in the foreground, so call it while the app is shown.
	public static void start(Context ctx, String reason) {
		Intent i = new Intent(ctx, KeepAliveService.class);
		i.putExtra(EXTRA_REASON, reason);
		if (Build.VERSION.SDK_INT >= 26) {
			ctx.startForegroundService(i);
		} else {
			ctx.startService(i);
		}
	}

	// stop stops the service and removes its notification.
	public static void stop(Context ctx) {
		ctx.stopService(new Intent(ctx, KeepAliveService.class));
	}

	@Override
	public int onStartCommand(Intent intent, int flags, int startId) {
		String reason = intent != null ? intent.getStringExtra(EXTRA_REASON) : null;
		Notification n = notification(reason);
		if (Build.VERSION.SDK_INT >= 29) {
			startForeground(NOTIFICATION_ID, n, ServiceInfo.FOREGROUND_SERVICE_TYPE_DATA_SYNC);
		} else {
			startForeground(NOTIFICATION_ID, n);
		}
		// If Android kills the process anyway, the work is gone: don't
		// restart the service with nothing to keep alive.
		return START_NOT_STICKY;
	}

	// onTimeout is called on Android 15 and later when the service has
	// used up its daily dataSync time; it must stop.
	public void onTimeout(int startId, int fgsType) {
		stopSelf();
	}

	@Override
	public IBinder onBind(Intent intent) {
		return null;
	}

	@SuppressWarnings("deprecation")
	private Notification notification(String reason) {
		CharSequence label = getApplicationInfo().loadLabel(getPackageManager());
		Notification.Builder b;
		if (Build.VERSION.SDK_INT >= 26) {
			NotificationManager nm = (NotificationManager) getSystemService(Context.NOTIFICATION_SERVICE);
			if (nm.getNotificationChannel(CHANNEL) == null) {
				NotificationChannel ch = new NotificationChannel(CHANNEL, "Background work", NotificationManager.IMPORTANCE_LOW);
				nm.createNotificationChannel(ch);
			}
			b = new Notification.Builder(this, CHANNEL);
		} else {
			b = new Notification.Builder(this);
		}
		b.setSmallIcon(android.R.drawable.stat_notify_sync)
				.setContentTitle(label)
				.setOngoing(true);
		if (reason != null && !reason.isEmpty()) {
			b.setContentText(reason);
		}
		// A tap brings the app back.
		Intent launch = getPackageManager().getLaunchIntentForPackage(getPackageName());
		if (launch != null) {
			int fl = PendingIntent.FLAG_UPDATE_CURRENT;
			if (Build.VERSION.SDK_INT >= 23) {
				fl |= PendingIntent.FLAG_IMMUTABLE;
			}
			b.setContentIntent(PendingIntent.getActivity(this, 0, launch, fl));
		}
		return b.build();
	}
}

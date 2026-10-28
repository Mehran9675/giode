package com.android.billingclient.api;

public interface BillingClientStateListener {
	void onBillingSetupFinished(BillingResult result);
	void onBillingServiceDisconnected();
}

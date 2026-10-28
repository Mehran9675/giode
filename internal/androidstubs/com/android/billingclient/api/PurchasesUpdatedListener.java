package com.android.billingclient.api;

public interface PurchasesUpdatedListener {
	void onPurchasesUpdated(BillingResult result, java.util.List<Purchase> purchases);
}

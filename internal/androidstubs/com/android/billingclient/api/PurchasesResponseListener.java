package com.android.billingclient.api;

public interface PurchasesResponseListener {
	void onQueryPurchasesResponse(BillingResult result, java.util.List<Purchase> purchases);
}

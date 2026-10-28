package com.android.billingclient.api;

public abstract class BillingClient {
	public static Builder newBuilder(android.content.Context context) { throw new UnsupportedOperationException(); }
	public abstract boolean isReady();
	public abstract void startConnection(BillingClientStateListener listener);
	public abstract void endConnection();
	public abstract BillingResult launchBillingFlow(android.app.Activity activity, BillingFlowParams params);
	public abstract void queryProductDetailsAsync(QueryProductDetailsParams params, ProductDetailsResponseListener listener);
	public abstract void queryPurchasesAsync(QueryPurchasesParams params, PurchasesResponseListener listener);
	public abstract void acknowledgePurchase(AcknowledgePurchaseParams params, AcknowledgePurchaseResponseListener listener);
	public abstract void consumeAsync(ConsumeParams params, ConsumeResponseListener listener);

	public static final class Builder {
		public Builder setListener(PurchasesUpdatedListener listener) { throw new UnsupportedOperationException(); }
		public Builder enablePendingPurchases(PendingPurchasesParams params) { throw new UnsupportedOperationException(); }
		public BillingClient build() { throw new UnsupportedOperationException(); }
	}
}

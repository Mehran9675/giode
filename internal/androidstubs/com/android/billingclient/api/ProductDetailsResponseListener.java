package com.android.billingclient.api;

public interface ProductDetailsResponseListener {
	void onProductDetailsResponse(BillingResult result, java.util.List<ProductDetails> details);
}

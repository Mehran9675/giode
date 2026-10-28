package com.android.billingclient.api;

import java.util.List;

public class BillingFlowParams {
	public static Builder newBuilder() { throw new UnsupportedOperationException(); }

	public static class Builder {
		public Builder setProductDetailsParamsList(List<ProductDetailsParams> list) { throw new UnsupportedOperationException(); }
		public BillingFlowParams build() { throw new UnsupportedOperationException(); }
	}

	public static final class ProductDetailsParams {
		public static Builder newBuilder() { throw new UnsupportedOperationException(); }

		public static class Builder {
			public Builder setProductDetails(ProductDetails details) { throw new UnsupportedOperationException(); }
			public Builder setOfferToken(String token) { throw new UnsupportedOperationException(); }
			public ProductDetailsParams build() { throw new UnsupportedOperationException(); }
		}
	}
}

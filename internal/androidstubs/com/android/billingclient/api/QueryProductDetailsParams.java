package com.android.billingclient.api;

import java.util.List;

public final class QueryProductDetailsParams {
	public static Builder newBuilder() { throw new UnsupportedOperationException(); }

	public static final class Builder {
		public Builder setProductList(List<Product> list) { throw new UnsupportedOperationException(); }
		public QueryProductDetailsParams build() { throw new UnsupportedOperationException(); }
	}

	public static final class Product {
		public static Builder newBuilder() { throw new UnsupportedOperationException(); }

		public static final class Builder {
			public Builder setProductId(String id) { throw new UnsupportedOperationException(); }
			public Builder setProductType(String type) { throw new UnsupportedOperationException(); }
			public Product build() { throw new UnsupportedOperationException(); }
		}
	}
}

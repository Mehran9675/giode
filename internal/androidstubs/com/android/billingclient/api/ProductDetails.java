package com.android.billingclient.api;

import java.util.List;

public final class ProductDetails {
	public String getProductId() { throw new UnsupportedOperationException(); }
	public String getProductType() { throw new UnsupportedOperationException(); }
	public String getTitle() { throw new UnsupportedOperationException(); }
	public String getName() { throw new UnsupportedOperationException(); }
	public String getDescription() { throw new UnsupportedOperationException(); }
	public OneTimePurchaseOfferDetails getOneTimePurchaseOfferDetails() { throw new UnsupportedOperationException(); }
	public List<SubscriptionOfferDetails> getSubscriptionOfferDetails() { throw new UnsupportedOperationException(); }

	public static final class OneTimePurchaseOfferDetails {
		public String getFormattedPrice() { throw new UnsupportedOperationException(); }
		public long getPriceAmountMicros() { throw new UnsupportedOperationException(); }
		public String getPriceCurrencyCode() { throw new UnsupportedOperationException(); }
	}

	public static final class SubscriptionOfferDetails {
		public String getOfferToken() { throw new UnsupportedOperationException(); }
		public PricingPhases getPricingPhases() { throw new UnsupportedOperationException(); }
	}

	public static class PricingPhases {
		public List<PricingPhase> getPricingPhaseList() { throw new UnsupportedOperationException(); }
	}

	public static final class PricingPhase {
		public String getFormattedPrice() { throw new UnsupportedOperationException(); }
		public long getPriceAmountMicros() { throw new UnsupportedOperationException(); }
		public String getPriceCurrencyCode() { throw new UnsupportedOperationException(); }
		public String getBillingPeriod() { throw new UnsupportedOperationException(); }
		public int getRecurrenceMode() { throw new UnsupportedOperationException(); }
	}
}

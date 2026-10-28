// SPDX-License-Identifier: Unlicense OR MIT

package org.giode.store;

import android.app.Activity;
import android.content.Context;
import android.os.Handler;
import android.os.Looper;
import android.view.View;

import com.android.billingclient.api.AcknowledgePurchaseParams;
import com.android.billingclient.api.AcknowledgePurchaseResponseListener;
import com.android.billingclient.api.BillingClient;
import com.android.billingclient.api.BillingClientStateListener;
import com.android.billingclient.api.BillingFlowParams;
import com.android.billingclient.api.BillingResult;
import com.android.billingclient.api.ConsumeParams;
import com.android.billingclient.api.ConsumeResponseListener;
import com.android.billingclient.api.PendingPurchasesParams;
import com.android.billingclient.api.ProductDetails;
import com.android.billingclient.api.ProductDetailsResponseListener;
import com.android.billingclient.api.Purchase;
import com.android.billingclient.api.PurchasesResponseListener;
import com.android.billingclient.api.PurchasesUpdatedListener;
import com.android.billingclient.api.QueryProductDetailsParams;
import com.android.billingclient.api.QueryPurchasesParams;

import org.json.JSONArray;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

// store_android connects Google Play Billing to giode/store. Results go
// to Go as JSON with the request's id; purchases that complete outside a
// request (pending ones approved later) go to onUpdate.
public final class store_android {
	// Response codes (BillingClient.BillingResponseCode).
	private static final int OK = 0;
	private static final int USER_CANCELED = 1;
	private static final int ITEM_ALREADY_OWNED = 7;

	private static BillingClient client;
	private static final ArrayList<Runnable> waiting = new ArrayList<Runnable>();
	private static boolean connecting;
	private static final HashMap<String, ProductDetails> products = new HashMap<String, ProductDetails>();
	// purchasing is the id of the purchase flow in progress, or 0.
	private static int purchasing;

	static native void onResult(int id, String json, String err);
	static native void onUpdate(String json);

	private static final PurchasesUpdatedListener updates = new PurchasesUpdatedListener() {
		@Override public void onPurchasesUpdated(BillingResult r, List<Purchase> list) {
			int id;
			synchronized (store_android.class) {
				id = purchasing;
				purchasing = 0;
			}
			int code = r.getResponseCode();
			if (code != OK || list == null || list.isEmpty()) {
				if (id != 0) {
					onResult(id, null, code == USER_CANCELED ? "canceled" : error(r));
				}
				return;
			}
			for (int i = 0; i < list.size(); i++) {
				String js = purchaseJSON(list.get(i)).toString();
				if (i == 0 && id != 0) {
					onResult(id, js, null);
				} else {
					onUpdate(js);
				}
			}
		}
	};

	private static String error(BillingResult r) {
		String msg = r.getDebugMessage();
		return "Play Billing error " + r.getResponseCode() + (msg != null && !msg.isEmpty() ? ": " + msg : "");
	}

	// connected runs r once the billing client is connected, or reports
	// the failure to request id.
	private static void connected(Context ctx, final int id, final Runnable r) {
		synchronized (store_android.class) {
			if (client == null) {
				client = BillingClient.newBuilder(ctx.getApplicationContext())
					.setListener(updates)
					.enablePendingPurchases(PendingPurchasesParams.newBuilder().enableOneTimeProducts().build())
					.build();
			}
			if (client.isReady()) {
				r.run();
				return;
			}
			waiting.add(new Runnable() {
				@Override public void run() {
					if (client.isReady()) {
						r.run();
					} else {
						onResult(id, null, "can't connect to Google Play");
					}
				}
			});
			if (connecting) {
				return;
			}
			connecting = true;
		}
		client.startConnection(new BillingClientStateListener() {
			@Override public void onBillingSetupFinished(BillingResult result) {
				flush();
			}

			@Override public void onBillingServiceDisconnected() {
				flush();
			}
		});
	}

	private static void flush() {
		ArrayList<Runnable> rs;
		synchronized (store_android.class) {
			connecting = false;
			rs = new ArrayList<Runnable>(waiting);
			waiting.clear();
		}
		for (Runnable r : rs) {
			r.run();
		}
	}

	// products queries ids, both one-time products and subscriptions.
	public static void products(Context ctx, final int id, String idsJSON) {
		final List<String> ids = new ArrayList<String>();
		try {
			JSONArray a = new JSONArray(idsJSON);
			for (int i = 0; i < a.length(); i++) {
				ids.add(a.getString(i));
			}
		} catch (Exception e) {
			onResult(id, null, e.toString());
			return;
		}
		connected(ctx, id, new Runnable() {
			@Override public void run() {
				final JSONArray out = new JSONArray();
				query(ids, "inapp", out, new Runnable() {
					@Override public void run() {
						query(ids, "subs", out, new Runnable() {
							@Override public void run() {
								onResult(id, out.toString(), null);
							}
						});
					}
				});
			}
		});
	}

	private static void query(List<String> ids, String type, final JSONArray out, final Runnable next) {
		List<QueryProductDetailsParams.Product> list = new ArrayList<QueryProductDetailsParams.Product>();
		for (String pid : ids) {
			list.add(QueryProductDetailsParams.Product.newBuilder().setProductId(pid).setProductType(type).build());
		}
		client.queryProductDetailsAsync(QueryProductDetailsParams.newBuilder().setProductList(list).build(), new ProductDetailsResponseListener() {
			@Override public void onProductDetailsResponse(BillingResult r, List<ProductDetails> details) {
				if (details != null) {
					for (ProductDetails d : details) {
						synchronized (store_android.class) {
							products.put(d.getProductId(), d);
						}
						out.put(productJSON(d));
					}
				}
				next.run();
			}
		});
	}

	private static JSONObject productJSON(ProductDetails d) {
		JSONObject o = new JSONObject();
		try {
			o.put("id", d.getProductId());
			o.put("title", d.getName());
			o.put("description", d.getDescription());
			ProductDetails.OneTimePurchaseOfferDetails one = d.getOneTimePurchaseOfferDetails();
			if (one != null) {
				o.put("price", one.getFormattedPrice());
				o.put("micros", one.getPriceAmountMicros());
				o.put("currency", one.getPriceCurrencyCode());
			}
			List<ProductDetails.SubscriptionOfferDetails> offers = d.getSubscriptionOfferDetails();
			if (offers != null && !offers.isEmpty()) {
				o.put("subscription", true);
				// The base plan's recurring phase: the last one.
				List<ProductDetails.PricingPhase> phases = offers.get(0).getPricingPhases().getPricingPhaseList();
				if (!phases.isEmpty()) {
					ProductDetails.PricingPhase p = phases.get(phases.size() - 1);
					o.put("price", p.getFormattedPrice());
					o.put("micros", p.getPriceAmountMicros());
					o.put("currency", p.getPriceCurrencyCode());
					o.put("period", p.getBillingPeriod());
				}
			}
		} catch (Exception e) {
		}
		return o;
	}

	private static JSONObject purchaseJSON(Purchase p) {
		JSONObject o = new JSONObject();
		try {
			o.put("id", p.getOrderId());
			List<String> ps = p.getProducts();
			o.put("product", ps.isEmpty() ? "" : ps.get(0));
			o.put("token", p.getPurchaseToken());
			o.put("time", p.getPurchaseTime());
			// Purchase.PurchaseState: 1 purchased, 2 pending.
			o.put("pending", p.getPurchaseState() == 2);
			o.put("acknowledged", p.isAcknowledged());
			o.put("receipt", p.getOriginalJson());
			o.put("signature", p.getSignature());
		} catch (Exception e) {
		}
		return o;
	}

	// purchase starts the purchase flow for a product queried before.
	public static void purchase(View view, final int id, final String productID) {
		final Context ctx = view.getContext();
		if (!(ctx instanceof Activity)) {
			onResult(id, null, "no activity");
			return;
		}
		connected(ctx, id, new Runnable() {
			@Override public void run() {
				final ProductDetails d;
				synchronized (store_android.class) {
					d = products.get(productID);
				}
				if (d == null) {
					onResult(id, null, "product " + productID + " wasn't loaded with Products");
					return;
				}
				BillingFlowParams.ProductDetailsParams.Builder pb = BillingFlowParams.ProductDetailsParams.newBuilder().setProductDetails(d);
				List<ProductDetails.SubscriptionOfferDetails> offers = d.getSubscriptionOfferDetails();
				if (offers != null && !offers.isEmpty()) {
					pb.setOfferToken(offers.get(0).getOfferToken());
				}
				List<BillingFlowParams.ProductDetailsParams> list = new ArrayList<BillingFlowParams.ProductDetailsParams>();
				list.add(pb.build());
				final BillingFlowParams params = BillingFlowParams.newBuilder().setProductDetailsParamsList(list).build();
				new Handler(Looper.getMainLooper()).post(new Runnable() {
					@Override public void run() {
						synchronized (store_android.class) {
							if (purchasing != 0) {
								onResult(id, null, "another purchase is in progress");
								return;
							}
							purchasing = id;
						}
						BillingResult r = client.launchBillingFlow((Activity) ctx, params);
						if (r.getResponseCode() != OK) {
							synchronized (store_android.class) {
								purchasing = 0;
							}
							onResult(id, null, r.getResponseCode() == ITEM_ALREADY_OWNED ? "owned" : error(r));
						}
					}
				});
			}
		});
	}

	// finish acknowledges a purchase, or consumes it so it can be bought
	// again.
	public static void finish(Context ctx, final int id, final String token, final boolean consume) {
		connected(ctx, id, new Runnable() {
			@Override public void run() {
				if (consume) {
					client.consumeAsync(ConsumeParams.newBuilder().setPurchaseToken(token).build(), new ConsumeResponseListener() {
						@Override public void onConsumeResponse(BillingResult r, String t) {
							onResult(id, "{}", r.getResponseCode() == OK ? null : error(r));
						}
					});
					return;
				}
				client.acknowledgePurchase(AcknowledgePurchaseParams.newBuilder().setPurchaseToken(token).build(), new AcknowledgePurchaseResponseListener() {
					@Override public void onAcknowledgePurchaseResponse(BillingResult r) {
						onResult(id, "{}", r.getResponseCode() == OK ? null : error(r));
					}
				});
			}
		});
	}

	// owned lists the purchases the user holds: unconsumed products and
	// active subscriptions.
	public static void owned(Context ctx, final int id) {
		connected(ctx, id, new Runnable() {
			@Override public void run() {
				final JSONArray out = new JSONArray();
				client.queryPurchasesAsync(QueryPurchasesParams.newBuilder().setProductType("inapp").build(), new PurchasesResponseListener() {
					@Override public void onQueryPurchasesResponse(BillingResult r, List<Purchase> list) {
						if (r.getResponseCode() != OK) {
							onResult(id, null, error(r));
							return;
						}
						for (Purchase p : list) {
							out.put(purchaseJSON(p));
						}
						client.queryPurchasesAsync(QueryPurchasesParams.newBuilder().setProductType("subs").build(), new PurchasesResponseListener() {
							@Override public void onQueryPurchasesResponse(BillingResult r, List<Purchase> list) {
								if (r.getResponseCode() == OK) {
									for (Purchase p : list) {
										out.put(purchaseJSON(p));
									}
								}
								onResult(id, out.toString(), null);
							}
						});
					}
				});
			}
		});
	}
}

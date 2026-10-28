package com.google.android.gms.tasks;

public abstract class Task<TResult> {
	public abstract boolean isSuccessful();
	public abstract TResult getResult();
	public abstract Exception getException();
	public Task<TResult> addOnCompleteListener(OnCompleteListener<TResult> listener) { throw new UnsupportedOperationException(); }
}

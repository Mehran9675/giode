//go:build !((linux && !android) || freebsd || openbsd || netbsd || windows || darwin || android)

package i18n

func systemLocale() string { return "" }

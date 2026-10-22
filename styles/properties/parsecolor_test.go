package properties

import (
	"image/color"
	"testing"
)

func TestParseColor(t *testing.T) {
	cases := map[string]color.NRGBA{
		"#fff":                     {255, 255, 255, 255},
		"#11223344":                {0x11, 0x22, 0x33, 0x44},
		"3b82f6":                   {0x3b, 0x82, 0xf6, 255},
		"rgb(59, 130, 246)":        {59, 130, 246, 255},
		"rgb(59 130 246 / 50%)":    {59, 130, 246, 128},
		"rgba(255,0,0,0.25)":       {255, 0, 0, 64},
		"rgb(100%, 0%, 0%)":        {255, 0, 0, 255},
		"hsl(0, 100%, 50%)":        {255, 0, 0, 255},
		"hsl(120deg 100% 25%)":     {0, 128, 0, 255},
		"hsla(240, 100%, 50%, .5)": {0, 0, 255, 128},
		"RebeccaPurple":            {0x66, 0x33, 0x99, 255},
		"transparent":              {},
	}
	for in, want := range cases {
		got, err := ParseColor(in)
		if err != nil {
			t.Errorf("ParseColor(%q): %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("ParseColor(%q) = %v, want %v", in, got, want)
		}
	}
	for _, bad := range []string{"", "#zzz", "#12345", "rgb(1,2)", "hsl(x, 1%, 1%)", "notacolor"} {
		if _, err := ParseColor(bad); err == nil {
			t.Errorf("ParseColor(%q): want error", bad)
		}
	}
}

func TestCalcColor_NeverPanics(t *testing.T) {
	if got := CalcColor("#zz"); got != fallbackColor {
		t.Errorf("invalid color = %v, want white fallback", got)
	}
	if got := CalcColor(""); got.A != 0 {
		t.Errorf("empty color = %v, want transparent", got)
	}
}

func TestCalcColor_ThemeAlpha(t *testing.T) {
	resetTheme(t)
	full := CalcColor("theme.primary")
	half := CalcColor("theme.primary/0.5")
	if half.R != full.R || half.A != 128 {
		t.Errorf("theme.primary/0.5 = %v, want primary at alpha 128", half)
	}
	if got := CalcColor("theme.primary/50%"); got != half {
		t.Errorf("percent alpha = %v, want %v", got, half)
	}
}

func TestSetPalette_EmptyFieldFallsBack(t *testing.T) {
	resetTheme(t)
	SetPalette(ModeLight, Palette{Primary: "#ff00ff"})
	if got, want := CalcColor("theme.danger"), CalcColor(DefaultLightPalette.Danger); got != want {
		t.Errorf("unset palette field = %v, want default %v", got, want)
	}
}

func TestThemeChangeHookAndVersion(t *testing.T) {
	resetTheme(t)
	calls := 0
	SetThemeChangeHook(func() { calls++ })
	defer SetThemeChangeHook(nil)
	v := ThemeVersion()
	SetMode(ModeDark)
	SetMode(ModeDark) // no change: no notification
	if calls != 1 || ThemeVersion() != v+1 {
		t.Errorf("calls = %d, version delta = %d, want 1, 1", calls, ThemeVersion()-v)
	}
}

func TestCalcColorReverse_RoundTrip(t *testing.T) {
	c := color.NRGBA{1, 2, 254, 128}
	if got := CalcColor(CalcColorReverse(c)); got != c {
		t.Errorf("round trip = %v, want %v", got, c)
	}
}

func TestLerpColor(t *testing.T) {
	if got := CalcColor(LerpColor("#000000", "#ffffff", 0.5)); got.R < 126 || got.R > 129 {
		t.Errorf("midpoint = %v", got)
	}
	if got := CalcColor(LerpColor("", "#ff0000", 0.5)); got.R != 255 || got.A < 126 || got.A > 129 {
		t.Errorf("fade in from unset = %v, want red at half alpha", got)
	}
}

func TestEasing_CSSCurves(t *testing.T) {
	for _, e := range []Easing{EaseInOut, Ease, EaseIn, EaseOut, EaseLinear, "cubic-bezier(0.4, 0, 0.2, 1)"} {
		if a, b := e.Apply(0), e.Apply(1); a > 1e-3 || b < 1-1e-3 {
			t.Errorf("%s: endpoints %v, %v, want 0 and 1", e, a, b)
		}
	}
	if v := EaseIn.Apply(0.5); v > 0.4 {
		t.Errorf("ease-in at 0.5 = %v, want slow start (< 0.4)", v)
	}
	if v := EaseOut.Apply(0.5); v < 0.6 {
		t.Errorf("ease-out at 0.5 = %v, want fast start (> 0.6)", v)
	}
	if v := EaseInOut.Apply(0.5); v < 0.49 || v > 0.51 {
		t.Errorf("ease-in-out at 0.5 = %v, want 0.5 (symmetric)", v)
	}
	// CSS ease-in-out at 25% is ~0.129.
	if v := EaseInOut.Apply(0.25); v < 0.12 || v > 0.14 {
		t.Errorf("ease-in-out at 0.25 = %v, want ~0.129", v)
	}
	peak := float32(0)
	for i := 0; i <= 100; i++ {
		peak = max(peak, EaseOutBack.Apply(float32(i)/100))
	}
	if peak <= 1 {
		t.Error("ease-out-back doesn't overshoot")
	}
	if Easing("cubic-bezier(2, 0, 0, 1)").Apply(0.5) != EaseInOut.Apply(0.5) {
		t.Error("invalid cubic-bezier doesn't fall back to ease-in-out")
	}
}

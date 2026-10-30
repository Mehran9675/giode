package properties

import "testing"

func TestMix(t *testing.T) {
	if got := Mix("#ffffff", "#101115", 0.05); got != "#1C1D21FF" {
		t.Errorf("Mix 5%% white = %s, want #1C1D21FF", got)
	}
	if got := Mix("#ff0000", "#0000ff", 0); got != "#0000FFFF" {
		t.Errorf("t=0 gave %s", got)
	}
	if got := Mix("#ff0000", "#0000ff", 2); got != "#FF0000FF" {
		t.Errorf("t>1 gave %s", got)
	}
}

func TestTabWidth(t *testing.T) {
	for _, tc := range []struct {
		w        TabWidth
		in, want string
	}{
		{0, "a\tb", "a   b"},
		{0, "\tx\n\ty", "    x\n    y"},
		{2, "ab\tc", "ab  c"},
		{8, "no tabs", "no tabs"},
	} {
		if got := tc.w.Expand(tc.in); got != tc.want {
			t.Errorf("TabWidth(%d).Expand(%q) = %q, want %q", tc.w, tc.in, got, tc.want)
		}
	}
}

func TestMixFollowsTheme(t *testing.T) {
	defer SetMode(CurrentMode())
	c := Mix("theme.text", "theme.background", 0.5)
	SetMode(ModeLight)
	light := CalcColor(c)
	SetMode(ModeDark)
	dark := CalcColor(c)
	if light == dark {
		t.Errorf("%s resolved to %v in both modes", c, light)
	}
	want := mixNRGBA(CalcColor("theme.text"), CalcColor("theme.background"), 0.5)
	if dark != want {
		t.Errorf("%s = %v in dark mode, want %v", c, dark, want)
	}
	if got := CalcColor(Mix(c, "#000000", 1)); got != dark {
		t.Errorf("nested mix = %v, want %v", got, dark)
	}
}

func TestColorMix(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"color-mix(in srgb, #ffffff 5%, #101115)", "#1C1D21FF"},
		{"color-mix(in srgb, #ff0000, #0000ff 100%)", "#0000FFFF"},
		{"color-mix(in srgb, rgb(255 0 0), #0000ff)", "#800080FF"},
	} {
		c, err := ParseColor(tc.in)
		if err != nil {
			t.Fatal(err)
		}
		if got := CalcColorReverse(c); got != tc.want {
			t.Errorf("%s = %s, want %s", tc.in, got, tc.want)
		}
	}
}

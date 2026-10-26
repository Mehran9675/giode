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

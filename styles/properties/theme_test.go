package properties

import "testing"

// resetTheme restores the default mode and palettes after a test that
// changes global theme state, so tests don't leak into each other.
func resetTheme(t *testing.T) {
	t.Cleanup(func() {
		SetMode(ModeLight)
		SetPalette(ModeLight, DefaultLightPalette)
		SetPalette(ModeDark, DefaultDarkPalette)
	})
}

func TestCalcColor_ThemeToken(t *testing.T) {
	resetTheme(t)

	if got, want := CalcColor("theme.primary"), CalcColor(DefaultLightPalette.Primary); got != want {
		t.Errorf("theme.primary in light mode = %v, want %v", got, want)
	}

	SetMode(ModeDark)
	if got, want := CalcColor("theme.primary"), CalcColor(DefaultDarkPalette.Primary); got != want {
		t.Errorf("theme.primary in dark mode = %v, want %v", got, want)
	}
	if got, want := CalcColor("theme.background"), CalcColor(DefaultDarkPalette.Background); got != want {
		t.Errorf("theme.background in dark mode = %v, want %v", got, want)
	}
}

func TestSetPalette_OverridesResolution(t *testing.T) {
	resetTheme(t)

	custom := DefaultLightPalette
	custom.Primary = "#ff00ff"
	SetPalette(ModeLight, custom)

	if got, want := CalcColor("theme.primary"), Hex("#ff00ff"); got != want {
		t.Errorf("theme.primary after SetPalette = %v, want %v", got, want)
	}
	// Fields left untouched still resolve to the built-in default.
	if got, want := CalcColor("theme.background"), CalcColor(DefaultLightPalette.Background); got != want {
		t.Errorf("theme.background after partial SetPalette = %v, want %v", got, want)
	}
}

func TestCalcColor_UnknownThemeToken(t *testing.T) {
	resetTheme(t)

	if got, want := CalcColor("theme.nonexistent"), CalcColor("white"); got != want {
		t.Errorf("unknown theme token = %v, want white fallback %v", got, want)
	}
}

func TestCalcColor_PlainColorsUnaffectedByMode(t *testing.T) {
	resetTheme(t)
	SetMode(ModeDark)

	if got, want := CalcColor("#3b82f6"), Hex("#3b82f6"); got != want {
		t.Errorf("plain hex color changed under dark mode: got %v, want %v", got, want)
	}
	if got, want := CalcColor("red"), Hex("#ff0000"); got != want {
		t.Errorf("plain named color changed under dark mode: got %v, want %v", got, want)
	}
}

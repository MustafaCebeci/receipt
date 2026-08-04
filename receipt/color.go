package receipt

// Color represents a PDF color.
type Color struct {
	R, G, B uint8
}

// Predefined colors.
var (
	// Black is pure black (0, 0, 0).
	Black = Color{0, 0, 0}

	// White is pure white (255, 255, 255).
	White = Color{255, 255, 255}

	// Transparent represents no fill color.
	Transparent = Color{0, 0, 0}

	// LightGray is RGB(238, 238, 238).
	LightGray = Color{238, 238, 238}

	// DarkGray is RGB(128, 128, 128).
	DarkGray = Color{128, 128, 128}
)

// RGB creates a new color from RGB values (0-255).
func RGB(r, g, b uint8) Color {
	return Color{r, g, b}
}

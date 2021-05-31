package tui

import "github.com/muesli/termenv"

var (
	term = termenv.ColorProfile()
)

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Bold().Styled
}

var (
	verylightBlue = makeFgStyle("45")
	lightBlue     = makeFgStyle("39")
	blue          = makeFgStyle("27")
	darkBlue      = makeFgStyle("20")
	grey          = makeFgStyle("#888888")
)

// ColorizeWhale the whale base on index
func ColorizeWhale(i int, s string) string {
	switch i {
	case 0:
		return verylightBlue(s)
	case 1:
		return lightBlue(s)
	case 2:
		return blue(s)
	case 3:
		return darkBlue(s)
	default:
		return grey(s)
	}
}

// ColorizeWater the water base on index
func ColorizeWater(i int, s string) string {
	switch i {
	case 0:
		return darkBlue(s)
	case 1:
		return blue(s)
	case 2:
		return lightBlue(s)
	case 3:
		return verylightBlue(s)
	case 4:
		return grey(s)
	default:
		return s
	}
}

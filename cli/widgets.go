package cli

import (
	"github.com/gdamore/tcell/v2"
)

type Region struct {
	W, H    int
	SetCell func(x, y int, style tcell.Style, ch rune)
}

func TuiRegion(tui tcell.Screen, x, y, w, h int) Region {
	return Region{
		W: w, H: h,
		SetCell: func(dx, dy int, style tcell.Style, ch rune) {
			if dx >= 0 && dx < w && dy >= 0 && dy < h {
				tui.SetCell(x+dx, y+dy, style, ch)
			}
		},
	}
}

var (
	whiteOnBlue = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)
)

func drawText(region Region, style tcell.Style, text string) {
	for x, ch := range text {
		region.SetCell(x, 0, style, ch)
	}
}

type funcReader func([]byte) (int, error)

func (f funcReader) Read(p []byte) (int, error) { return f(p) }

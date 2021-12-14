package cli

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type BufView struct {
	// TODO: Wrap bool
	Y   int // Y of the view in the Buf, for down/up scrolling
	X   int // X of the view in the Buf, for left/right scrolling
	Buf string
}

func (v *BufView) DrawTo(region Region) {
	r := bufio.NewReader(strings.NewReader(v.Buf))

	// PgDn/PgUp etc. support
	for y := v.Y; y > 0; y-- {
		line, err := r.ReadBytes('\n')
		switch err {
		case nil:
			// skip line
			continue
		case io.EOF:
			r = bufio.NewReader(bytes.NewReader(line))
			y = 0
			break
		default:
			panic(err)
		}
	}

	lclip := false
	drawch := func(x, y int, ch rune) {
		if x <= v.X && v.X != 0 {
			x, ch = 0, '«'
			lclip = true
		} else {
			x -= v.X
		}
		if x >= region.W {
			x, ch = region.W-1, '»'
		}
		region.SetCell(x, y, tcell.StyleDefault, ch)
	}
	endline := func(x, y int) {
		x -= v.X
		if x < 0 {
			x = 0
		}
		if x == 0 && lclip {
			x++
		}
		lclip = false
		for ; x < region.W; x++ {
			region.SetCell(x, y, tcell.StyleDefault, ' ')
		}
	}

	x, y := 0, 0
	// TODO: handle runes properly, including their visual width (mattn/go-runewidth)
	for {
		ch, _, err := r.ReadRune()
		if y >= region.H || err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		switch ch {
		case '\n':
			endline(x, y)
			x, y = 0, y+1
			continue
		case '\t':
			const tabwidth = 8
			drawch(x, y, ' ')
			for x%tabwidth < (tabwidth - 1) {
				x++
				if x >= region.W {
					break
				}
				drawch(x, y, ' ')
			}
		default:
			drawch(x, y, ch)
		}
		x++
	}
	for ; y < region.H; y++ {
		endline(x, y)
		x = 0
	}
}

func (v *BufView) HandleKey(ev *tcell.EventKey, scrollY int) bool {
	const scrollX = 8 // When user scrolls horizontally, move by this many characters
	switch getKey(ev) {
	//
	// Vertical scrolling
	//
	case key(tcell.KeyUp):
		v.Y--
		v.normalizeY()
	case key(tcell.KeyDown):
		v.Y++
		v.normalizeY()
	case key(tcell.KeyPgDn):
		// TODO: in top-right corner of Buf area, draw current line number & total # of lines
		v.Y += scrollY
		v.normalizeY()
	case key(tcell.KeyPgUp):
		v.Y -= scrollY
		v.normalizeY()
	//
	// Horizontal scrolling
	//
	case altKey(tcell.KeyLeft),
		ctrlKey(tcell.KeyLeft):
		v.X -= scrollX
		if v.X < 0 {
			v.X = 0
		}
	case altKey(tcell.KeyRight),
		ctrlKey(tcell.KeyRight):
		v.X += scrollX
	case altKey(tcell.KeyHome),
		ctrlKey(tcell.KeyHome):
		v.X = 0
	default:
		// Unknown key/combination, not handled
		return false
	}
	return true
}

func (v *BufView) normalizeY() {
	nlines := count(strings.NewReader(v.Buf), '\n') + 1
	if v.Y >= nlines {
		v.Y = nlines - 1
	}
	if v.Y < 0 {
		v.Y = 0
	}
}

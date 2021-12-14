package cli

import (
	"bytes"
	"io"

	"github.com/gdamore/tcell/v2"
)

func NewEditor(prompt string) *Editor {
	return &Editor{prompt: []rune(prompt)}
}

type Editor struct {
	// TODO: make editor multiline. Reuse gocui or something for this?
	prompt []rune
	value  []rune
	cursor int
	// lastw is length of value on last Draw; we need it to know how much to erase after backspace
	lastw int
}

func (e *Editor) String() string { return string(e.value) }

func (e *Editor) DrawTo(region Region, setcursor func(x, y int)) {
	// Draw prompt & the edited value - use white letters on blue background
	style := whiteOnBlue
	for i, ch := range e.prompt {
		region.SetCell(i, 0, style, ch)
	}
	for i, ch := range e.value {
		region.SetCell(len(e.prompt)+i, 0, style, ch)
	}

	// Clear remains of last value if needed
	for i := len(e.value); i < e.lastw; i++ {
		region.SetCell(len(e.prompt)+i, 0, tcell.StyleDefault, ' ')
	}
	e.lastw = len(e.value)

	// Show cursor if requested
	if setcursor != nil {
		setcursor(len(e.prompt)+e.cursor, 0)
	}
}

func (e *Editor) HandleKey(ev *tcell.EventKey) bool {
	// If a character is entered, with no modifiers except maybe shift, then just insert it
	if ev.Key() == tcell.KeyRune && ev.Modifiers()&(^tcell.ModShift) == 0 {
		e.insert(ev.Rune())
		return true
	}
	// Handle editing & movement keys
	switch getKey(ev) {
	case key(tcell.KeyBackspace), key(tcell.KeyBackspace2):
		// See https://github.com/nsf/termbox-go/issues/145
		e.delete(-1)
	case key(tcell.KeyDelete):
		e.delete(0)
	case key(tcell.KeyLeft):
		if e.cursor > 0 {
			e.cursor--
		}
	case key(tcell.KeyRight):
		if e.cursor < len(e.value) {
			e.cursor++
		}
	default:
		// Unknown key/combination, not handled
		return false
	}
	return true
}

func (e *Editor) insert(ch rune) {
	// Insert character into value (https://github.com/golang/go/wiki/SliceTricks#insert)
	e.value = append(e.value, 0)
	copy(e.value[e.cursor+1:], e.value[e.cursor:])
	e.value[e.cursor] = ch
	e.cursor++
}

func (e *Editor) delete(dx int) {
	pos := e.cursor + dx
	if pos < 0 || pos >= len(e.value) {
		return
	}
	e.value = append(e.value[:pos], e.value[pos+1:]...)
	e.cursor = pos
}

func count(r io.Reader, b byte) (n int) {
	buf := [256]byte{}
	for {
		i, err := r.Read(buf[:])
		n += bytes.Count(buf[:i], []byte{b})
		if err != nil {
			return
		}
	}
}

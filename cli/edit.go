package cli

import (
	"github.com/awesome-gocui/gocui"
)

type customEditor struct {
	keypressCh chan string
}

func NewCustomEditor(keypressCh chan string) *customEditor {
	return &customEditor{keypressCh}
}

// simpleEditor is used as the default gocui editor.
func (ce *customEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if ch != 0 && mod == 0 {
		v.EditWrite(ch)
		ce.keypressCh <- v.Buffer()
		return
	}

	switch key {
	case gocui.KeyCtrlC:
		return
	case gocui.KeySpace:
		v.EditWrite(' ')
	case gocui.KeyBackspace, gocui.KeyBackspace2:
		v.EditDelete(true)
	case gocui.KeyDelete:
		v.EditDelete(false)
	// case gocui.KeyInsert:
	// 	v.Overwrite = !v.Overwrite
	// case gocui.KeyEnter:
	// 	v.EditNewLine()
	case gocui.KeyArrowDown:
		v.MoveCursor(0, 1)
		return
	case gocui.KeyArrowUp:
		v.MoveCursor(0, -1)
		return
	case gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0)
		return
	case gocui.KeyArrowRight:
		v.MoveCursor(1, 0)
		return
	case gocui.KeyTab:
		v.EditWrite('\t')
	case gocui.KeyEsc:
		// If not here the esc key will act like the KeySpace
	default:
		v.EditWrite(ch)
	}

	ce.keypressCh <- v.Buffer()
}

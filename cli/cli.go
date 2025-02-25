package cli

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/owenthereal/jqplay/config"
	"github.com/owenthereal/jqplay/jq"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

const defaultMessage = ` [jqplay] ^C exit  PgUp/PgDn/Up/Dn/^</^> scroll`

type ScreenSize int

const (
	SizeDefault ScreenSize = iota
	SizeTiny
	SizeSmall
	SizeNormal
)

type Cli struct {
	conf *config.Config

	app *tview.Application

	grid        *tview.Grid
	editorField *tview.InputField
	footer      *tview.TextView
	inputView   *tview.TextView
	outputView  *tview.TextView

	editorCh    chan string
	lastFilter  string
	lastCommand string

	screenW    int
	screenSize ScreenSize

	opts map[string]*jq.JQOpt

	log         []string
	debugLogger *log.Logger
	debugBuf    *bytes.Buffer
}

func New(conf *config.Config) *Cli {
	dl := log.New()
	dl.Level = log.DebugLevel
	buf := &bytes.Buffer{}
	dl.Out = buf

	opts := make(map[string]*jq.JQOpt)
	opts["slurp"] = &jq.JQOpt{"slurp", false}
	opts["null-input"] = &jq.JQOpt{"null-input", false}
	opts["compact-output"] = &jq.JQOpt{"compact-output", false}
	opts["raw-input"] = &jq.JQOpt{"raw-input", false}
	opts["raw-output"] = &jq.JQOpt{"raw-output", false}

	return &Cli{conf: conf,
		editorCh:    make(chan string, 10),
		opts:        opts,
		debugLogger: dl,
		debugBuf:    buf,
	}
}

func (c *Cli) exitDebug() {
	fmt.Println("Caught exit")
	fmt.Println("replaying log:")
	for i, l := range c.log {
		fmt.Printf("%d: %s\n", i, l)
	}
	fmt.Println(c.debugBuf)
}

func (c *Cli) toggleOpt(name string) {
	opt := c.opts[name]
	opt.Enabled = !opt.Enabled
	c.updateFooter()
	c.runJq(c.lastFilter)
}

func (c *Cli) updateFooter() {
	msg := defaultMessage

	o := func(opt *jq.JQOpt, key string) string {
		if key == "" {
			key = opt.Name[0:1]
		}
		color := "[white]"
		if opt.Enabled {
			color = "[green::b]"
		}
		switch c.screenSize {
		case SizeTiny:
			return fmt.Sprintf("  [%s%s[white::]]", color, "--"+opt.Name)
		case SizeSmall:
			return fmt.Sprintf("  [%s=%s%t[white::]]", opt.Name, color, opt.Enabled)
		default:
			return fmt.Sprintf("  [alt+%s %s=%s%t[white::]]", key, opt.Name, color, opt.Enabled)
		}
	}

	msg += o(c.opts["slurp"], "")
	msg += o(c.opts["null-input"], "")
	msg += o(c.opts["compact-output"], "")
	msg += o(c.opts["raw-input"], "i")
	msg += o(c.opts["raw-output"], "o")

	c.footer.SetText(msg)
}

func (c *Cli) createViews() {
	c.editorField = tview.NewInputField().SetLabel("JQ Filter: ")
	c.editorField.SetChangedFunc(func(text string) {
		// c.debugLogger.Infof("added text: %s", text)
		c.editorCh <- text
	})
	c.editorField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			c.editorCh <- c.editorField.GetText()
		case tcell.KeyTAB:
			c.app.SetFocus(c.inputView)
			c.inputView.SetBorderColor(tcell.ColorRed)
			c.outputView.SetBorderColor(tcell.ColorWhite)
			c.editorField.SetFieldBackgroundColor(tcell.ColorBlack)
		case tcell.KeyBacktab:
			c.app.SetFocus(c.outputView)
			c.outputView.SetBorderColor(tcell.ColorRed)
			c.inputView.SetBorderColor(tcell.ColorWhite)
			c.editorField.SetFieldBackgroundColor(tcell.ColorBlack)
		}
	})

	c.footer = tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignLeft)
	c.footer.SetBackgroundColor(tcell.ColorBlue)
	c.updateFooter()

	c.grid = tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0, 0, 0).
		SetBorders(false).
		AddItem(c.editorField, 0, 0, 1, 3, 0, 0, true).
		AddItem(c.footer, 2, 0, 1, 3, 0, 0, false)

	c.inputView = tview.NewTextView().SetText(c.conf.JSON).SetWrap(true).SetScrollable(true)
	c.inputView.SetBorder(true).SetTitle("JSON Input").SetTitleAlign(tview.AlignLeft)
	c.inputView.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTAB:
			c.app.SetFocus(c.outputView)
			c.outputView.SetBorderColor(tcell.ColorRed)
			c.inputView.SetBorderColor(tcell.ColorWhite)
			c.editorField.SetFieldBackgroundColor(tcell.ColorBlack)
		case tcell.KeyBacktab:
			c.app.SetFocus(c.editorField)
			c.inputView.SetBorderColor(tcell.ColorWhite)
			c.outputView.SetBorderColor(tcell.ColorWhite)
			c.editorField.SetFieldBackgroundColor(tcell.ColorBlue)
		}
	})

	c.outputView = tview.NewTextView().SetWrap(true).SetScrollable(true)
	c.outputView.SetBorder(true).SetTitle("JQ Output").SetTitleAlign(tview.AlignLeft)
	c.outputView.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTAB:
			c.app.SetFocus(c.editorField)
			c.inputView.SetBorderColor(tcell.ColorWhite)
			c.outputView.SetBorderColor(tcell.ColorWhite)
			c.editorField.SetFieldBackgroundColor(tcell.ColorBlue)

		case tcell.KeyBacktab:
			c.app.SetFocus(c.inputView)
			c.inputView.SetBorderColor(tcell.ColorRed)
			c.outputView.SetBorderColor(tcell.ColorWhite)
			c.editorField.SetFieldBackgroundColor(tcell.ColorBlack)
		}
	})

	flex := tview.NewFlex().
		AddItem(c.inputView, 0, 1, false).
		AddItem(c.outputView, 0, 1, false)

	c.grid.AddItem(flex, 1, 0, 1, 3, 0, 0, false)
}

func (c *Cli) Start() error {
	// defer c.exitDebug()
	defer func() {
		if c.lastCommand != "" {
			fmt.Println(c.lastCommand)
		}
	}()

	c.app = tview.NewApplication()

	c.createViews()

	// Add a hook to catch screen resizing
	drawDebouncer := &Debouncer{}
	c.app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		drawDebouncer.Do(time.Millisecond*100, func() {
			c.resizeHook(screen)
		})
		return false
	})

	go debounce(time.Millisecond*150, c.editorCh, c.runJq)

	// custom keybinds at the app level, global shortcuts
	c.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers()&tcell.ModAlt != 0 {
			switch event.Rune() {
			case 's':
				c.toggleOpt("slurp")
				return nil
			case 'n':
				c.toggleOpt("null-input")
				return nil
			case 'c':
				c.toggleOpt("compact-output")
				return nil
			case 'i':
				c.toggleOpt("raw-input")
				return nil
			case 'o':
				c.toggleOpt("raw-output")
				return nil
			}
		}
		return event
	})

	if err := c.app.SetRoot(c.grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil
}

func (c *Cli) resizeHook(screen tcell.Screen) {
	w, _ := screen.Size()
	if w != c.screenW {
		c.screenW = w
		var size ScreenSize
		if w < 138 {
			size = SizeTiny
		} else if w < 166 {
			size = SizeSmall
		} else {
			size = SizeNormal
		}
		if size != c.screenSize {
			// changed
			c.screenSize = size
			go c.app.QueueUpdateDraw(func() {
				c.updateFooter()
			})
		}
	}
}

func (c *Cli) updateOutput(str string) {
	go c.app.QueueUpdateDraw(func() {
		c.outputView.Clear().SetText(str).ScrollToBeginning()
	})
}

func (c *Cli) runJq(filter string) {
	c.lastFilter = filter
	if filter != "" {
		var opts []jq.JQOpt
		for _, opt := range c.opts {
			opts = append(opts, *opt)
		}
		j := &jq.JQ{J: c.conf.JSON, Q: filter, O: opts}
		c.lastCommand = j.CommandString()
		buff := &bytes.Buffer{}
		err := j.Eval(context.Background(), buff)
		str := buff.String()
		if err != nil {
			c.updateOutput(fmt.Sprintf("jq error: %s\n--\n%s", err, str))
		} else {
			c.updateOutput(str)
		}
	} else {
		c.updateOutput("[enter jq filter]")
	}
}

// debounce events received on the given channel, acting only only the last item
// received in the interval.
func debounce(interval time.Duration, input chan string, cb func(arg string)) {
	var item string
	timer := time.NewTimer(interval)
	for {
		select {
		case item = <-input:
			timer.Reset(interval)
		case <-timer.C:
			cb(item)
		}
	}
}

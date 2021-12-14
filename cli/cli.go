package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/gdamore/tcell/v2"
	"github.com/owenthereal/jqplay/config"
	"github.com/owenthereal/jqplay/jq"
	log "github.com/sirupsen/logrus"
)

const defaultMessage = ` ^C exit  PgUp/PgDn/Up/Dn/^</^> scroll  [jqplay]`

type Cli struct {
	conf *config.Config

	g          *gocui.Gui
	editorView *gocui.View
	inputView  *gocui.View
	outputView *gocui.View

	editorCh chan string

	log         []string
	debugLogger *log.Logger
	debugBuf    *bytes.Buffer
}

func New(conf *config.Config) *Cli {
	dl := log.New()
	dl.Level = log.DebugLevel
	buf := &bytes.Buffer{}
	dl.Out = buf
	return &Cli{conf: conf, editorCh: make(chan string, 10), debugLogger: dl, debugBuf: buf}
}

func (c *Cli) exitDebug() {
	// c.tui.Fini()
	c.g.Close()
	fmt.Println("Caught exit")
	fmt.Println("replaying log:")
	for i, l := range c.log {
		fmt.Printf("%d: %s\n", i, l)
	}
	fmt.Println(c.debugBuf)
}

func (c *Cli) Start() error {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManager(c)

	if err := c.keybindings(g); err != nil {
		log.Panicln(err)
	}

	go debounce(time.Millisecond*150, c.editorCh, c.runJq)

	// debug
	defer c.exitDebug()
	c.g = g

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}

	return nil
}

// quit program
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// vertically scroll the given view
func scrollView(v *gocui.View, dx, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox+dx, oy+dy); err != nil {
			return err
		}
	}
	return nil
}

// set up keybindings in various views
func (c *Cli) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	// scroll pane
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 0, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 0, 1)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

func (c *Cli) Layout(g *gocui.Gui) error {
	// c.debugLogger.Infof("layout called")

	maxX, maxY := g.Size()

	// editor panel
	if v, err := g.SetView("editor", 0, 0, maxX-1, 2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Enter JQ Filter"
		v.Editable = true
		v.Wrap = false
		v.Editor = NewCustomEditor(c.editorCh)
		c.editorView = v
	}

	// input panel
	if v, err := g.SetView("input", 0, 3, maxX/2-1, maxY-3, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "JSON Input"
		v.Highlight = true
		v.Editable = false
		v.Wrap = true
		v.WriteString(c.conf.JSON)
		c.inputView = v
	}

	// output panel
	if v, err := g.SetView("main", maxX/2, 3, maxX-1, maxY-3, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "JQ Output"
		v.Editable = false
		v.Wrap = true
		v.WriteString("[enter jq filter]")
		c.outputView = v
	}

	// help panel
	if v, err := g.SetView("help", 0, maxY-3, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Frame = false
		v.Editable = false
		v.Wrap = true
		v.BgColor = gocui.Attribute(tcell.ColorBlue)
		v.FgColor = gocui.Attribute(tcell.ColorWhite)
		v.WriteString(defaultMessage)
	}

	if _, err := g.SetCurrentView("editor"); err != nil {
		return err
	}
	return nil
}

func (c *Cli) updateOutput(str string) {
	c.g.Update(func(g *gocui.Gui) error {
		c.outputView.Clear()
		c.outputView.WriteString(str)
		return nil
	})
}

func (c *Cli) runJq(command string) {
	// c.log = append(c.log, fmt.Sprintf("running: %s", command))
	// c.debugLogger.Infof("running: %s", command)
	if command != "" {
		// TODO: wrap in commandSubprocess? do we need to kill it?
		j := &jq.JQ{J: c.conf.JSON, Q: command}
		buff := &bytes.Buffer{}
		err := j.Eval(context.Background(), buff)
		str := buff.String()
		// c.debugLogger.Infof("ran jq with command %s, got error? %s", command, err)
		// c.debugLogger.Infof("jq output:\n%s", str)
		if err != nil {
			c.updateOutput(fmt.Sprintf("jq error: %s\n--\n%s", err, str))
		} else {
			c.updateOutput(str)
		}
	} else {
		// c.debugLogger.Infof("skipped running jq, writing default string")
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

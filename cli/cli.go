package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/owenthereal/jqplay/config"
	"github.com/owenthereal/jqplay/jq"
)

const defaultMessage = `^X exit (^C nosave)  PgUp/PgDn/Up/Dn/^</^> scroll  ^S pause (^Q end)  [jqplay]`

type Cli struct {
	conf *config.Config

	tui           tcell.Screen
	commandEditor *Editor
	commandOutput BufView
	message       string

	log []string
}

func New(conf *config.Config) *Cli {
	return &Cli{conf: conf}

}

func (c *Cli) exitDebug() {
	c.tui.Fini()
	fmt.Println("Caught exit")
	fmt.Println("replaying log:")
	for i, l := range c.log {
		fmt.Printf("%d: %s\n", i, l)
	}
}

func (c *Cli) Start() error {
	c.tui = initTUI()
	defer c.tui.Fini()

	// Initialize 3 main UI parts
	// The top line of the TUI is an editable command, which will be used
	// as a pipeline for data we read from stdin
	c.commandEditor = NewEditor("jq filter> ")
	// The rest of the screen is a view of the results of the command
	c.commandOutput = BufView{Buf: c.conf.JSON}
	// Sometimes, a message may be displayed at the bottom of the screen, with help or other info
	c.message = defaultMessage

	// debug
	defer c.exitDebug()

	// Main loop
	c.runLoop()

	return nil
}

func (c *Cli) runLoop() error {
	commandsCh := make(chan string)
	go debounce(time.Millisecond*150, commandsCh, c.runJq)

	lastCommand := ""
	for {
		// If user edited the command, immediately run it in background, and
		// kill the previously running command.
		command := c.commandEditor.String()
		if command != lastCommand {
			c.log = append(c.log, command)
			commandsCh <- command

			// commandSubprocess.Kill()
			// if command != "" {
			// 	commandSubprocess = StartSubprocess(command, stdinCapture, func() { triggerRefresh(tui) })
			// 	commandOutput.Buf = commandSubprocess.Buf
			// } else {
			// 	// If command is empty, show original input data again (~ equivalent of typing `cat`)
			// 	commandSubprocess = nil
			// 	commandOutput.Buf = stdinCapture
			// }
		}

		lastCommand = command

		// Draw UI
		_, h := c.tui.Size()
		c.drawUI()

		// TODO: reset message to default at some point

		// Handle UI events
		switch ev := c.tui.PollEvent().(type) {
		// Key pressed
		case *tcell.EventKey:
			// Is it a command editor key?
			if c.commandEditor.HandleKey(ev) {
				// c.message = ""
				continue
			}
			// Is it a command output view key?
			if c.commandOutput.HandleKey(ev, h-1) {
				// c.message = ""
				continue
			}
			// Some other global key combinations
			switch getKey(ev) {
			case key(tcell.KeyCtrlS),
				ctrlKey(tcell.KeyCtrlS):
				// stdinCapture.Pause(true)
				c.triggerRefresh()
			case key(tcell.KeyCtrlQ),
				ctrlKey(tcell.KeyCtrlQ):
				// stdinCapture.Pause(false)
				lastCommand = ":" // Make sure we restart current command
			case key(tcell.KeyCtrlC),
				ctrlKey(tcell.KeyCtrlC),
				key(tcell.KeyCtrlD),
				ctrlKey(tcell.KeyCtrlD):
				// Quit
				return nil
			case key(tcell.KeyCtrlX),
				ctrlKey(tcell.KeyCtrlX):
				// Write script 'upN.sh' and quit
				// writeScript(commandEditor.String(), tui)
				return nil
			}
		}
	}
}

func (c *Cli) runJq(command string) {
	c.log = append(c.log, fmt.Sprintf("running: %s", command))
	if command != "" {
		// TODO: wrap in commandSubprocess? do we need to kill it?
		j := &jq.JQ{J: c.conf.JSON, Q: command}
		buff := &bytes.Buffer{}
		if err := j.Eval(context.Background(), buff); err != nil {
			c.commandOutput.Buf = fmt.Sprintf("jq error: %s\n--\n%s", err, buff.String())
		} else {
			c.commandOutput.Buf = buff.String()
		}
	} else {
		c.commandOutput.Buf = c.conf.JSON
	}
	// manually redraw since we run jq async
	c.triggerRefresh()
}

func (c *Cli) drawUI() {
	// Draw UI
	w, h := c.tui.Size()
	// stdinCapture.DrawStatus(TuiRegion(tui, 0, 0, 1, 1))
	c.commandEditor.DrawTo(TuiRegion(c.tui, 1, 0, w-1, 1), func(x, y int) { c.tui.ShowCursor(x+1, 0) })
	c.commandOutput.DrawTo(TuiRegion(c.tui, 0, 1, w, h-1))
	drawText(TuiRegion(c.tui, 0, h-1, w, 1), whiteOnBlue, c.message)
	c.tui.Show()
}

func (c *Cli) triggerRefresh() {
	c.tui.PostEvent(tcell.NewEventInterrupt(nil))
}

func initTUI() tcell.Screen {
	// if isatty.IsTerminal(os.Stdin.Fd()) {
	// 	die("up requires some data piped on standard input, for example try: `echo hello world | up`")
	// }

	// Init TUI code
	// TODO: maybe try gocui or termbox?
	tui, err := tcell.NewScreen()
	if err != nil {
		die(err.Error())
	}
	err = tui.Init()
	if err != nil {
		die(err.Error())
	}
	return tui
}

func die(message string) {
	os.Stderr.WriteString("error: " + message + "\n")
	os.Exit(1)
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

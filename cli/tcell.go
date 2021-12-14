package cli

import "github.com/gdamore/tcell/v2"

type key int32

func getKey(ev *tcell.EventKey) key { return key(ev.Modifiers())<<16 + key(ev.Key()) }
func altKey(base tcell.Key) key     { return key(tcell.ModAlt)<<16 + key(base) }
func ctrlKey(base tcell.Key) key    { return key(tcell.ModCtrl)<<16 + key(base) }

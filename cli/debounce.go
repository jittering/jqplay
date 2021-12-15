package cli

import (
	"sync"
	"time"
)

type Debouncer struct {
	lock  sync.Mutex
	timer *time.Timer
}

func (d *Debouncer) Do(duration time.Duration, f func()) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.timer != nil {
		stop := d.timer.Stop()
		if !stop {
			d.timer = nil
			return
		}
	}

	d.timer = time.AfterFunc(duration, func() {
		d.lock.Lock()
		defer d.lock.Unlock()
		f()
		d.timer = nil
	})
}

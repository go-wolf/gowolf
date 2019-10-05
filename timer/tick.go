package timer

import "time"

func Tick(d time.Duration,handler func())  {
	t := time.NewTicker(d)
	for {
		select {
		case <-t.C:
			handler()
		}
	}
}
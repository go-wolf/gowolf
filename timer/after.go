package timer

import (
	"time"
)

func After(d time.Duration,handler func())  {
	t := time.NewTicker(d)
	<-t.C
	handler()
}
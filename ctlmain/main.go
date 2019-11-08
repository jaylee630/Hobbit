package ctlmain

import (
	"os"
	"os/signal"
)

// Main is real entry-point
func Main() {
	if !parseFlags() {
		return
	}

	ctl := Controller{}
	ctl.Init()
	ctl.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctl.Stop()
}

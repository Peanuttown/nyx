package tap

import (
	"os"
	"os/signal"
	"testing"
)

func TestNewTap(t *testing.T) {
	_, err := NewTap("testDev")
	if err != nil {
		t.Fatal(err)
	}
	notify := make(chan os.Signal)
	signal.Notify(notify, os.Interrupt)
	<-notify
}

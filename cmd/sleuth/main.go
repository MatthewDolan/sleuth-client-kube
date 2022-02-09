package main

import (
	sleuthkube "github.com/MatthewDolan/sleuth-client-kube"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stopCh := make(chan struct{})

	// Shutdown gracefully on SIGTERM
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		close(stopCh)
	}()

	// Run blocks until stopCh is closed
	if err := sleuthkube.Run(stopCh, os.Args[1:]); err != nil {
		panic(err)
	}
}

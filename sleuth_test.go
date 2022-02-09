package sleuthkube

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestRun(t *testing.T) {
	stopCh := make(chan struct{})

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		require.NoError(t, Run(stopCh, []string{}))
	}()

	close(stopCh)
	waitGroup.Wait()
}

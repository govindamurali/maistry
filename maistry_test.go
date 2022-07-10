package maistry_test

import (
	"github.com/stretchr/testify/assert"
	"maistry"
	"testing"
	"time"
)

type TestClient struct {
	jobManager        *maistry.JobManager
	internalFluPerSec int
	maxClientQps      int
	waitMiliSec       int
	name              string
}

var count1 int
var count2 int

var dummyLogger maistry.ILogger

func TestDispatcher_Start(t *testing.T) {
	t.Skip()

	clients := []TestClient{

		{
			jobManager:        maistry.NewJobManager(1000, "1"),
			internalFluPerSec: 100,
			maxClientQps:      2,
			waitMiliSec:       200,
			name:              "1",
		},
		{
			jobManager:        maistry.NewJobManager(1, "2"),
			internalFluPerSec: 10,
			maxClientQps:      10,
			waitMiliSec:       300,
			name:              "2",
		},
	}

	dispatcher := maistry.NewDispatcher(1, dummyLogger)

	for _, c := range clients {
		dispatcher.AddJobManager(c.jobManager)
		c.jobManager.Run()
	}

	dispatcher.Start()

	for _, c := range clients {
		SendData(c)
	}

	time.Sleep(time.Second * time.Duration(5))
	assert.True(t, count1 > 16 && count1 < 20)
	assert.True(t, count2 > 3 && count2 < 6)
}

func SendData(c TestClient) {

	go func() {
		ticker := time.Tick(time.Duration(1000/c.internalFluPerSec) * time.Millisecond)

		for {
			<-ticker
			c.jobManager.PushJob(maistry.NewJob(func() {
				time.Sleep(time.Duration(c.waitMiliSec) * time.Millisecond)
				if c.name == "1" {
					count1++
				}
				if c.name == "2" {
					count2++
				}
			}, dummyLogger))
		}
	}()

}

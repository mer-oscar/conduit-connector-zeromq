package zeromq

import (
	"context"
	"fmt"
	"testing"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/zeromq/goczmq"
	"go.uber.org/goleak"
)

type zeroMQAcceptanceTestDriver struct {
	sdk.ConfigurableAcceptanceTestDriver
}

func TestAcceptance(t *testing.T) {
	go func() {
		pubChannel := goczmq.NewPubChanneler("tcp://127.0.0.1:5556")
		defer pubChannel.Destroy()
		ctx := context.Background()
		for {
			select {
			case msg := <-pubChannel.RecvChan:
				fmt.Println(msg)
			case <-ctx.Done():
				return
			}
		}
	}()

	sdk.AcceptanceTest(t, zeroMQAcceptanceTestDriver{
		ConfigurableAcceptanceTestDriver: sdk.ConfigurableAcceptanceTestDriver{
			Config: sdk.ConfigurableAcceptanceTestDriverConfig{
				Connector:         Connector,
				GoleakOptions:     []goleak.Option{goleak.IgnoreCurrent()},
				SourceConfig:      map[string]string{"portBindings": "tcp://*:5556", "topic": "a"},
				DestinationConfig: map[string]string{"routerEndpoints": "tcp://127.0.0.1:5555", "topic": "a"},
				GenerateDataType:  sdk.GenerateRawData,
				Skip: []string{
					"TestSource_Configure_RequiredParams",
					"TestDestination_Configure_RequiredParams",
					"TestSource_Open_ResumeAtPositionCDC",
					"TestSource_Open_ResumeAtPositionSnapshot",
				},
				WriteTimeout: 500 * time.Millisecond,
				ReadTimeout:  3000 * time.Millisecond,
			},
		},
	})
}

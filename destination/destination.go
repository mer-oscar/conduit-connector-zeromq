package destination

import (
	"context"
	"fmt"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/zeromq/goczmq"
)

type Destination struct {
	sdk.UnimplementedDestination

	config     DestinationConfig
	pubChannel *goczmq.Channeler
}

type DestinationConfig struct {
	// Config includes parameters that are the same in the source and destination.
	Config
}

func New() sdk.Destination {
	// Create Destination and wrap it in the default middleware.
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

func (d *Destination) Parameters() map[string]sdk.Parameter {
	// Parameters is a map of named Parameters that describe how to configure
	// the Destination. Parameters can be generated from DestinationConfig with
	// paramgen.
	return d.config.Parameters()
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	sdk.Logger(ctx).Info().Msg("Configuring Destination...")
	err := sdk.Util.ParseConfig(cfg, &d.config)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	return nil
}

func (d *Destination) Open(ctx context.Context) error {
	d.pubChannel = goczmq.NewPubChanneler(d.config.RouterEndpoints)

	return nil
}

func (d *Destination) Write(ctx context.Context, records []sdk.Record) (int, error) {
	var written int
	for _, rec := range records {
		d.pubChannel.SendChan <- [][]byte{[]byte(d.config.Topic), rec.Bytes()}
	}

	return written, nil
}

func (d *Destination) Teardown(ctx context.Context) error {
	if d.pubChannel != nil {
		d.pubChannel.Destroy()
	}
	return nil
}

package source

import (
	"context"
	"fmt"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/oklog/ulid/v2"
	"github.com/zeromq/goczmq"
)

type Source struct {
	sdk.UnimplementedSource

	config        SourceConfig
	routerChannel *goczmq.Channeler
	readBuffer    chan sdk.Record
}

type SourceConfig struct {
	Config
}

func New() sdk.Source {
	return sdk.SourceWithMiddleware(&Source{
		readBuffer: make(chan sdk.Record, 1),
	}, sdk.DefaultSourceMiddleware()...)
}

func (s *Source) Parameters() map[string]sdk.Parameter {
	return s.config.Parameters()
}

func (s *Source) Configure(ctx context.Context, cfg map[string]string) error {
	sdk.Logger(ctx).Info().Msg("Configuring Source...")
	err := sdk.Util.ParseConfig(cfg, &s.config)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	return nil
}

func (s *Source) Open(ctx context.Context, pos sdk.Position) error {
	s.routerChannel = goczmq.NewSubChanneler(s.config.PortBindings, s.config.Topic)
	return nil
}

func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	go s.listen(ctx)
	select {
	case rec := <-s.readBuffer:
		return rec, nil
	case <-ctx.Done():
		return sdk.Record{}, ctx.Err()
	}
}

func (s *Source) Ack(ctx context.Context, position sdk.Position) error {
	return nil
}

func (s *Source) Teardown(ctx context.Context) error {
	if s.routerChannel != nil {
		s.routerChannel.Destroy()
	}

	return nil
}

func (s *Source) listen(ctx context.Context) {
	select {
	case msg := <-s.routerChannel.RecvChan:
		if msg != nil {
			recFrame := msg[0]
			fmt.Println(len(msg))

			for _, frame := range msg[1:] {
				recBytes := frame

				recUlid := ulid.Make()

				rec := sdk.Util.Source.NewRecordCreate(
					sdk.Position(string(recFrame)+"_"+recUlid.String()),
					sdk.Metadata{
						"frame": string(recFrame),
					},
					nil,
					sdk.RawData(recBytes))

				s.readBuffer <- rec
			}
		}
	case <-ctx.Done():
		return
	}
}

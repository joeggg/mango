package embedded

import (
	"context"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

type EmbeddedPacket struct {
	Kind    int
	Size    int
	Command string
	RawData []byte
	Data    proto.Message
}

/*
	Parse the packet into a proto struct
*/
func (p *EmbeddedPacket) Parse(gatherers []Gatherer) error {
	name, result, err := GetEmbdeddedType(p.Kind)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(p.RawData, result)
	if err != nil {
		return err
	}
	p.Command = name
	p.Data = result
	// Check for registered handlers on gatherers
	if gatherers == nil {
		return nil
	}
	errs, _ := errgroup.WithContext(context.Background())
	for _, g := range gatherers {
		handler, ok := g.GetHandlers()[p.Kind]
		if !ok {
			return nil
		}
		// Run handlers concurrently
		errs.Go(func() error {
			return handler(p.Data)
		})
	}
	return errs.Wait()
}

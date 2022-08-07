package embedded

import (
	"context"

	"github.com/joeggg/mango/mappings"
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
func (p *EmbeddedPacket) Parse(gatherers map[string]Gatherer, lk *mappings.LookupObjects) error {
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
			continue
		}
		// Run handlers concurrently
		errs.Go(func() error {
			return handler(p.Data, lk)
		})
	}
	return errs.Wait()
}

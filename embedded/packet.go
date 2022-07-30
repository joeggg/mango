package embedded

import (
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
	for _, g := range gatherers {
		handlers := g.GetHandlers()
		handler, ok := handlers[p.Kind]
		if !ok {
			return nil
		}
		err = handler(p.Data)
		if err != nil {
			return err
		}
	}
	return nil
}

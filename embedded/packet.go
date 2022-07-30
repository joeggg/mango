package embedded

import (
	"google.golang.org/protobuf/proto"
)

type EmbeddedPacket struct {
	Kind    int
	Command string
	RawData []byte
	Data    proto.Message
}

/*
	Parse the packet into a proto struct
*/
func (p *EmbeddedPacket) Parse() error {
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
	// Check for registered handlers
	handlers, ok := EmbdeddedHandlers[p.Kind]
	if !ok {
		return nil
	}
	// Run each handler
	for _, handler := range handlers {
		err = handler(p.Data)
		if err != nil {
			return err
		}
	}
	return nil
}

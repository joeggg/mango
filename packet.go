package mango

import (
	"mango/embedded"
	"mango/pb"

	"google.golang.org/protobuf/proto"
)

type Packet struct {
	Kind         int
	Tick         int
	Size         int
	IsCompressed bool
	RawMessage   []byte
	Command      pb.EDemoCommands
	Message      proto.Message
	Embed        *embedded.EmbeddedPacket
}

func (p *Packet) Parse() error {
	p.Command = pb.EDemoCommands(p.Kind)
	result, err := GetPacketType(p.Command)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(p.RawMessage, result)
	if err != nil {
		return err
	}
	p.Message = result
	return nil
}

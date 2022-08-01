package packet

import (
	"github.com/joeggg/mango/embedded"
	"github.com/joeggg/mango/pb"
	"google.golang.org/protobuf/proto"
)

type Packet struct {
	Kind         int
	Tick         int
	Size         int
	IsCompressed bool
	RawMessage   []byte
	RawEmbed     []byte
	Command      pb.EDemoCommands
	Message      proto.Message
	Embed        *embedded.EmbeddedPacket
}

/*
	Parse the packet and process the result, dependent on the type
*/
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
	// Handle message and set embedded if there is one
	handler := PacketHandlers[p.Command]
	err = handler(p)
	if err != nil {
		return err
	}
	return nil
}

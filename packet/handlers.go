package packet

import (
	"mango/embedded"
	"mango/pb"

	"google.golang.org/protobuf/proto"
)

func HandleEmbeddedPacket(p *Packet, data proto.Message) error {
	payload := data.(*pb.CDemoPacket)
	decoder := embedded.NewEmbeddedDecoder(payload.Data)
	packet, err := decoder.Decode()
	if err != nil {
		return err
	}
	err = packet.Parse()
	if err != nil {
		return err
	}
	p.Embed = packet
	return nil
}

func HandlePlaceHolder(p *Packet, data proto.Message) error {
	return nil
}

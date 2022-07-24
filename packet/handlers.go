package packet

import (
	"mango/embedded"
	"mango/pb"

	"google.golang.org/protobuf/proto"
)

/*
	Process a packet with embedded data by parsing it
*/
func HandleEmbeddedPacket(message proto.Message) (*embedded.EmbeddedPacket, error) {
	info := message.(*pb.CDemoPacket)
	decoder := embedded.NewEmbeddedDecoder(info.Data)
	packet, err := decoder.Decode()
	if err != nil {
		return nil, err
	}
	err = packet.Parse()
	if err != nil {
		return nil, err
	}
	return packet, nil
}

/*
	Placeholder for unimplemented message types
*/
func HandlePlaceHolder(message proto.Message) (*embedded.EmbeddedPacket, error) {
	return nil, nil
}

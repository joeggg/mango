package embedded

import "google.golang.org/protobuf/proto"

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
	return nil
}

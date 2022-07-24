package embedded

import "google.golang.org/protobuf/proto"

type EmbeddedPacket struct {
	Kind    int
	RawData []byte
	Data    proto.Message
}

func (p *EmbeddedPacket) Parse() error {
	result, err := GetEmbdeddedType(p.Kind)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(p.RawData, result)
	if err != nil {
		return err
	}
	p.Data = result
	return nil
}

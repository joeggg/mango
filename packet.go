package mango

import (
	"fmt"
	"mango/pb"

	"google.golang.org/protobuf/proto"
)

type Packet struct {
	Kind         int
	Tick         int
	Size         int
	IsCompressed bool
	Message      []byte
}

func (p *Packet) Parse() (proto.Message, error) {
	if pb.EDemoCommands(p.Kind) == pb.EDemoCommands_DEM_FileInfo {
		data := &pb.CDemoFileInfo{}
		err := proto.Unmarshal(p.Message, data)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		return nil, fmt.Errorf("unknown protobuf type: %v", p.Kind)
	}
}

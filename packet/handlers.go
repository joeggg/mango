package packet

import (
	"encoding/json"
	"mango/embedded"
	"mango/pb"
	"os"

	"google.golang.org/protobuf/proto"
)

var StringTables = map[string][]*pb.CDemoStringTablesItemsT{}

/*
	Process a packet with embedded data by parsing it
*/
func HandleEmbedded(message proto.Message) (*embedded.EmbeddedPacket, error) {
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
	Process a string tables packet by saving the tables and putting them into memory
*/
func HandleStringTables(message proto.Message) (*embedded.EmbeddedPacket, error) {
	info := message.(*pb.CDemoStringTables)
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return nil, err
	}
	err = os.WriteFile("string_tables.json", data, 0755)
	if err != nil {
		return nil, err
	}
	for _, table := range info.Tables {
		StringTables[*table.TableName] = table.Items
	}
	return nil, nil
}

func HandleClassinfo(message proto.Message) (*embedded.EmbeddedPacket, error) {
	info := message.(*pb.CDemoClassInfo)
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return nil, err
	}
	err = os.WriteFile("class_info.json", data, 0755)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func HandleSendTables(message proto.Message) (*embedded.EmbeddedPacket, error) {
	info := message.(*pb.CDemoSendTables)
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return nil, err
	}
	err = os.WriteFile("send_tables.json", data, 0755)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

/*
	Placeholder for unimplemented message types
*/
func HandlePlaceHolder(message proto.Message) (*embedded.EmbeddedPacket, error) {
	return nil, nil
}

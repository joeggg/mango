package packet

import (
	"encoding/json"
	"mango/pb"
	"os"
)

var StringTables = map[string][]*pb.CDemoStringTablesItemsT{}

/*
	Process a packet with embedded data
*/
func HandleEmbedded(p *Packet) error {
	info := p.Message.(*pb.CDemoPacket)
	p.RawEmbed = info.Data
	return nil
}

/*
	Process a packet with embedded data
*/
func HandleFullEmbedded(p *Packet) error {
	info := p.Message.(*pb.CDemoPacket)
	p.RawEmbed = info.Data
	return nil
}

/*
	Process a string tables packet by saving the tables and putting them into memory
*/
func HandleStringTables(p *Packet) error {
	info := p.Message.(*pb.CDemoStringTables)
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile("string_tables.json", data, 0755)
	if err != nil {
		return err
	}
	for _, table := range info.Tables {
		StringTables[*table.TableName] = table.Items
	}
	return nil
}

func HandleClassinfo(p *Packet) error {
	info := p.Message.(*pb.CDemoClassInfo)
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile("class_info.json", data, 0755)
	if err != nil {
		return err
	}
	return nil
}

func HandleSendTables(p *Packet) error {
	info := p.Message.(*pb.CDemoSendTables)
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile("send_tables.json", data, 0755)
	if err != nil {
		return err
	}
	return nil
}

/*
	Placeholder for unimplemented message types
*/
func HandlePlaceHolder(p *Packet) error {
	return nil
}

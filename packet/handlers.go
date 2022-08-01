package packet

import (
	"encoding/json"
	"os"

	"github.com/joeggg/mango/pb"
	"github.com/joeggg/mango/sendtables"
)

var StringTables = map[string][]*pb.CDemoStringTablesItemsT{}
var Players = map[int]string{}

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

func HandleFileInfo(p *Packet) error {
	info := p.Message.(*pb.CDemoFileInfo)
	players := info.GameInfo.Dota.PlayerInfo
	for i, player := range players {
		Players[i] = player.GetHeroName()
	}
	Players[-1] = "no one"
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
	decoder := sendtables.TableDecoder{}
	flat, err := decoder.Decode(info.Data)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(flat, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile("send_tables.json", data, 0755)
	if err != nil {
		return err
	}
	return err
}

/*
	Placeholder for unimplemented message types
*/
func HandlePlaceHolder(p *Packet) error {
	return nil
}

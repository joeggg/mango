package packet

import (
	"encoding/json"
	"os"

	"github.com/joeggg/mango/mappings"
	"github.com/joeggg/mango/pb"
)

/*
	Process a packet with embedded data
*/
func handleEmbedded(p *Packet, lk *mappings.LookupObjects) error {
	info := p.Message.(*pb.CDemoPacket)
	p.RawEmbed = info.Data
	return nil
}

/*
	Save raw summary info and create lookup of ID to player
*/
func handleFileInfo(p *Packet, lk *mappings.LookupObjects) error {
	info := p.Message.(*pb.CDemoFileInfo)
	lk.Summary = info
	players := info.GameInfo.Dota.PlayerInfo
	for i, player := range players {
		lk.Players[i] = player
	}
	lk.Players[-1] = nil
	return nil
}

/*
	Save string tables as map of table name to list of items
*/
func handleStringTables(p *Packet, lk *mappings.LookupObjects) error {
	info := p.Message.(*pb.CDemoStringTables)
	for _, table := range info.Tables {
		lk.StringTables[*table.TableName] = table.Items
	}
	return nil
}

/*
	Create class ID to network name map
*/
func handleClassinfo(p *Packet, lk *mappings.LookupObjects) error {
	info := p.Message.(*pb.CDemoClassInfo)
	for _, class := range info.Classes {
		lk.ClassInfo[int(class.GetClassId())] = class.GetNetworkName()
	}
	return nil
}

/*
	Decode send tables into flattened serialiser ..tbc..
*/
func handleSendTables(p *Packet, lk *mappings.LookupObjects) error {
	info := p.Message.(*pb.CDemoSendTables)
	decoder := mappings.TableDecoder{}
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
func handlePlaceHolder(p *Packet, lk *mappings.LookupObjects) error {
	return nil
}

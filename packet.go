package mango

import (
	"errors"
	"fmt"
	"mango/pb"

	"google.golang.org/protobuf/proto"
)

type Packet struct {
	Kind         int
	Tick         int
	Size         int
	IsCompressed bool
	RawMessage   []byte
	Command      pb.EDemoCommands
	Message      proto.Message
}

var PacketTypeMap = map[pb.EDemoCommands]string{
	pb.EDemoCommands_DEM_Error:               "",
	pb.EDemoCommands_DEM_Stop:                "mango.CDemoStop",
	pb.EDemoCommands_DEM_FileHeader:          "mango.CDemoFileHeader",
	pb.EDemoCommands_DEM_FileInfo:            "mango.CDemoFileInfo",
	pb.EDemoCommands_DEM_SyncTick:            "mango.CDemoSyncTick",
	pb.EDemoCommands_DEM_SendTables:          "mango.CDemoSendTables",
	pb.EDemoCommands_DEM_ClassInfo:           "mango.CDemoClassInfo",
	pb.EDemoCommands_DEM_StringTables:        "mango.CDemoStringTables",
	pb.EDemoCommands_DEM_Packet:              "mango.CDemoPacket",
	pb.EDemoCommands_DEM_SignonPacket:        "mango.CDemoPacket",
	pb.EDemoCommands_DEM_ConsoleCmd:          "mango.CDemoConsoleCmd",
	pb.EDemoCommands_DEM_CustomData:          "mango.CDemoCustomData",
	pb.EDemoCommands_DEM_CustomDataCallbacks: "mango.CDemoCustomDataCallbacks",
	pb.EDemoCommands_DEM_UserCmd:             "mango.CDemoUserCmd",
	pb.EDemoCommands_DEM_FullPacket:          "mango.CDemoFullPacket",
	pb.EDemoCommands_DEM_SaveGame:            "mango.CDemoSaveGame",
	pb.EDemoCommands_DEM_SpawnGroups:         "mango.CDemoSpawnGroups",
	pb.EDemoCommands_DEM_Max:                 "",
}

func (p *Packet) Parse() error {
	p.Command = pb.EDemoCommands(p.Kind)
	switch p.Command {

	case pb.EDemoCommands_DEM_Error:
		err := errors.New("found a replay error packet")
		return err

	case pb.EDemoCommands_DEM_Stop:
		result := &pb.CDemoStop{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_FileHeader:
		result := &pb.CDemoFileHeader{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_FileInfo:
		result := &pb.CDemoFileInfo{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_SyncTick:
		result := &pb.CDemoSyncTick{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_SendTables:
		result := &pb.CDemoSendTables{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_ClassInfo:
		result := &pb.CDemoClassInfo{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_StringTables:
		result := &pb.CDemoStringTables{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_Packet, pb.EDemoCommands_DEM_SignonPacket:
		result := &pb.CDemoPacket{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_ConsoleCmd:
		result := &pb.CDemoConsoleCmd{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_CustomData:
		result := &pb.CDemoCustomData{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_CustomDataCallbacks:
		result := &pb.CDemoCustomDataCallbacks{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_UserCmd:
		result := &pb.CDemoUserCmd{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_FullPacket:
		result := &pb.CDemoFullPacket{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_SaveGame:
		result := &pb.CDemoSaveGame{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_SpawnGroups:
		result := &pb.CDemoSpawnGroups{}
		err := proto.Unmarshal(p.RawMessage, result)
		p.Message = result
		return err

	case pb.EDemoCommands_DEM_Max:
		err := errors.New("found a DEM Max packet")
		return err

	default:
		err := fmt.Errorf("unknown protobuf type: %v", p.Kind)
		return err
	}
}

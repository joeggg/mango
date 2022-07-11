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

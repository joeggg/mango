package mango

import (
	"errors"
	"fmt"
	"mango/pb"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type PacketHandler func(*Packet, proto.Message) error

var PacketHandlerMap = map[pb.EDemoCommands]PacketHandler{
	pb.EDemoCommands_DEM_Stop:                HandlePlaceHolder,
	pb.EDemoCommands_DEM_FileHeader:          HandlePlaceHolder,
	pb.EDemoCommands_DEM_FileInfo:            HandlePlaceHolder,
	pb.EDemoCommands_DEM_SyncTick:            HandlePlaceHolder,
	pb.EDemoCommands_DEM_SendTables:          HandlePlaceHolder,
	pb.EDemoCommands_DEM_ClassInfo:           HandlePlaceHolder,
	pb.EDemoCommands_DEM_StringTables:        HandlePlaceHolder,
	pb.EDemoCommands_DEM_Packet:              HandleEmbeddedPacket,
	pb.EDemoCommands_DEM_SignonPacket:        HandlePlaceHolder,
	pb.EDemoCommands_DEM_ConsoleCmd:          HandlePlaceHolder,
	pb.EDemoCommands_DEM_CustomData:          HandlePlaceHolder,
	pb.EDemoCommands_DEM_CustomDataCallbacks: HandlePlaceHolder,
	pb.EDemoCommands_DEM_UserCmd:             HandlePlaceHolder,
	pb.EDemoCommands_DEM_FullPacket:          HandlePlaceHolder,
	pb.EDemoCommands_DEM_SaveGame:            HandlePlaceHolder,
	pb.EDemoCommands_DEM_SpawnGroups:         HandlePlaceHolder,
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

func GetPacketType(command pb.EDemoCommands) (proto.Message, error) {
	t, ok := PacketTypeMap[command]
	if !ok {
		return nil, errors.New("unknown embedded message type")
	}
	if t == "" {
		return nil, errors.New(fmt.Sprintf("received a %s packet", command))
	}
	name := protoreflect.FullName(t)
	cls, err := protoregistry.GlobalTypes.FindMessageByName(name)
	if err != nil {
		return nil, err
	}
	data := cls.New().Interface()
	return data, nil
}
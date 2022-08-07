package mappings

import "github.com/joeggg/mango/pb"

type LookupObjects struct {
	Summary *pb.CDemoFileInfo
	Players map[int]*pb.CGameInfo_CDotaGameInfo_CPlayerInfo

	ClassInfo    map[int]string
	SendTables   *pb.CSVCMsg_FlattenedSerializer
	StringTables map[string][]*pb.CDemoStringTablesItemsT
}

func NewLookupObjects() *LookupObjects {
	return &LookupObjects{
		Players:      make(map[int]*pb.CGameInfo_CDotaGameInfo_CPlayerInfo),
		ClassInfo:    make(map[int]string),
		StringTables: make(map[string][]*pb.CDemoStringTablesItemsT),
	}
}

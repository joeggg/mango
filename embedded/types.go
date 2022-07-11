package embedded

import (
	"mango/pb"
	"reflect"
)

var EmbeddedTypeMap = map[int]reflect.Type{
	int(pb.NET_Messages_net_NOP):        reflect.TypeOf(pb.CNETMsg_NOP{}),
	int(pb.SVC_Messages_svc_ServerInfo): reflect.TypeOf(pb.CSVCMsg_ServerInfo{}),
}

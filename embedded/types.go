package embedded

import (
	"fmt"
	"mango/pb"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// Embedded type to proto struct name
var EmbeddedTypeMap = map[int]string{
	int(pb.NET_Messages_net_NOP):                  "mango.CNETMsg_NOP",
	int(pb.NET_Messages_net_Tick):                 "mango.CNETMsg_Tick",
	int(pb.NET_Messages_net_SetConVar):            "mango.CNETMsg_SetConVar",
	int(pb.SVC_Messages_svc_ServerInfo):           "mango.CSVCMsg_ServerInfo",
	int(pb.SVC_Messages_svc_ClearAllStringTables): "mango.CSVCMsg_ClearAllStringTables",
}

func GetEmbdeddedType(kind int) (string, proto.Message, error) {
	t, ok := EmbeddedTypeMap[kind]
	if !ok {
		return t, nil, fmt.Errorf("unknown embedded message type: %d", kind)
	}
	name := protoreflect.FullName(t)
	cls, err := protoregistry.GlobalTypes.FindMessageByName(name)
	if err != nil {
		return t, nil, err
	}
	data := cls.New().Interface()
	return t, data, nil
}

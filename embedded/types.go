package embedded

import (
	"errors"
	"mango/pb"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var EmbeddedTypeMap = map[int]string{
	int(pb.NET_Messages_net_NOP):        "mango.CNETMsg_NOP",
	int(pb.SVC_Messages_svc_ServerInfo): "mango.CSVCMsg_ServerInfo",
}

func GetEmbdeddedType(kind int) (proto.Message, error) {
	t, ok := EmbeddedTypeMap[kind]
	if !ok {
		return nil, errors.New("unknown embedded message type")
	}
	name := protoreflect.FullName(t)
	cls, err := protoregistry.GlobalTypes.FindMessageByName(name)
	if err != nil {
		return nil, err
	}
	data := cls.New().Interface()
	return data, nil
}

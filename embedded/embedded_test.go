package embedded_test

import (
	"encoding/base64"
	"fmt"
	"mango"
	"mango/embedded"
	"mango/pb"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestEmbeddedParse(t *testing.T) {
	packetMap := mango.LoadExamplePacketData("embedded_packets")
	for name, info := range packetMap {
		fmt.Printf("Testing packet type %s: \n\n", name)
		p, err := processEmbeddedFromJson(info)
		if err != nil {
			t.Error(err)
		}
		mango.PrintStruct(p.Data)
		fmt.Println()
	}
}

func TestRegisterHandler(t *testing.T) {
	packetMap := mango.LoadExamplePacketData("embedded_packets")
	info := packetMap["ServerInfo"]
	count := 0
	total := 3
	// Register multiple functions
	h := func(p proto.Message) error {
		count++
		return nil
	}
	for i := 0; i < total; i++ {
		embedded.RegisterHandler(int(pb.SVC_Messages_svc_ServerInfo), h)
	}
	// Parse and check all functions ran
	_, err := processEmbeddedFromJson(info)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%d/%d handlers ran\n", count, total)
	if count != total {
		t.Error("Not all handlers ran!")
	}
}

func processEmbeddedFromJson(
	info map[string]interface{},
) (*embedded.EmbeddedPacket, error) {
	data, _ := base64.StdEncoding.DecodeString(info["data"].(string))
	d := embedded.NewEmbeddedDecoder(data)
	p, err := d.Decode()
	if err != nil {
		return nil, err
	}
	err = p.Parse()
	if err != nil {
		return nil, err
	}
	return p, nil
}

package packet_test

import (
	"encoding/base64"
	"fmt"
	"mango"
	"mango/packet"
	"testing"
)

func TestPacketParse(t *testing.T) {
	packetMap := mango.LoadExamplePacketData("packets")
	for name, info := range packetMap {
		fmt.Printf("Testing packet type %s: \n", name)
		b, err := base64.StdEncoding.DecodeString(info["data"].(string))
		if err != nil {
			t.Error(err)
		}
		p := packet.Packet{RawMessage: b, Kind: int(info["code"].(float64))}
		err = p.Parse()
		if err != nil {
			t.Error(err)
		}
		mango.PrintStruct(p.Message)
		if p.Embed != nil {
			mango.PrintStruct(p.Embed)
		}
	}
}

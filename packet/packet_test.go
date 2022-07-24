package packet_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mango"
	"mango/packet"
	"os"
	"testing"
)

func TestPacketParse(t *testing.T) {
	packetMap := loadPacketData("packets")
	for kind, data := range packetMap {
		fmt.Printf("Testing packet type %d: ", kind)
		b, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			t.Error(err)
		}
		p := packet.Packet{RawMessage: b, Kind: kind}
		err = p.Parse()
		if err != nil {
			t.Error(err)
		}
		fmt.Println(p.Command)
		mango.PrintStruct(p.Message)
		if p.Embed != nil {
			mango.PrintStruct(p.Embed)
		}
	}
}

func loadPacketData(name string) map[int]string {
	file, err := os.ReadFile(fmt.Sprintf("../example_data/%s.json", name))
	if err != nil {
		panic(err)
	}
	data := make(map[int]string, 0)
	if err = json.Unmarshal(file, &data); err != nil {
		panic(err)
	}
	return data
}

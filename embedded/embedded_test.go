package embedded_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mango"
	"mango/embedded"
	"os"
	"testing"
)

func TestEmbeddedParse(t *testing.T) {
	data := loadPacketData("sign_on_packet")
	d := embedded.NewEmbeddedDecoder(data)
	p, err := d.Decode()
	if err != nil {
		t.Error(err)
	}
	err = p.Parse()
	if err != nil {
		t.Error(err)
	}
	mango.PrintStruct(p.Data)
}

func loadPacketData(name string) []byte {
	file, err := os.ReadFile(fmt.Sprintf("../example_data/%s.json", name))
	if err != nil {
		panic(err)
	}
	data := make(map[string]string, 0)
	if err = json.Unmarshal(file, &data); err != nil {
		panic(err)
	}
	b, err := base64.StdEncoding.DecodeString(data["data"])
	if err != nil {
		panic(err)
	}
	return b
}

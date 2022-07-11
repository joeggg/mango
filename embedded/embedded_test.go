package embedded_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mango/embedded"
	"os"
	"reflect"
	"testing"
)

func TestEmbeddedParser(t *testing.T) {
	data := loadPacketData("sign_on_packet")
	p := embedded.NewEmbeddedParser(data)
	err := p.Parse()
	if err != nil {
		t.Error(err)
	}
}

func TestEmbeddedTypeMap(t *testing.T) {
	kind := embedded.EmbeddedTypeMap[0]
	msg := reflect.New(kind).Elem().Interface()
	fmt.Println(msg)
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

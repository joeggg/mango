package mango

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintStruct(data any) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func LoadExamplePacketData(name string) map[string]map[string]interface{} {
	file, err := os.ReadFile(fmt.Sprintf("../testdata/%s.json", name))
	if err != nil {
		panic(err)
	}
	data := make(map[string]map[string]interface{}, 0)
	if err = json.Unmarshal(file, &data); err != nil {
		panic(err)
	}
	return data
}

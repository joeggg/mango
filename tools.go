package mango

import (
	"encoding/json"
	"fmt"
	"os"
)

const maxPrintSize = 1000

func PrintStruct(data any) {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	strOut := string(out)
	if len(strOut) > maxPrintSize {
		strOut = strOut[:maxPrintSize] + "\n\n...continued\n"
	}
	fmt.Println(strOut)
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

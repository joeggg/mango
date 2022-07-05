package mango

import (
	"encoding/json"
	"fmt"
)

func PrintStruct(data any) {
	out, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(out))
}

package mango

import (
	"encoding/json"
	"fmt"
)

func PrintStruct(data any) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

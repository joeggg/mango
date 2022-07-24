package mango_test

import (
	"fmt"
	"mango"
	"mango/pb"
	"testing"
)

func TestErrorPacket(t *testing.T) {
	_, err := mango.GetPacketType(pb.EDemoCommands_DEM_Error)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Error("No error given for error packet")
}

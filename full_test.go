package mango_test

import (
	"fmt"
	"mango"
	"testing"
)

func TestFull(*testing.T) {
	if p, err := mango.NewReplayParser("data/test.dem"); err != nil {
		fmt.Println(err)
		return
	} else if err = p.Initialise(); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Summary packet:")
		packet, err := p.GetPacket()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(packet)
		err = packet.Parse()
		if err != nil {
			fmt.Println(err)
		}
	}
}

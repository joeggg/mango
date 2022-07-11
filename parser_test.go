package mango_test

import (
	"fmt"
	"mango"
	"testing"
)

const testReplayFilename = "example_data/test.dem"

func TestSummary(t *testing.T) {
	if p, err := mango.NewReplayParser(testReplayFilename); err != nil {
		t.Error(err)
	} else if err = p.Initialise(); err != nil {
		t.Error(err)
	} else if summary, err := p.GetSummary(); err != nil {
		t.Error(err)
	} else {
		fmt.Println("Summary packet:")
		err = mango.PrintStruct(summary)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestParse(t *testing.T) {
	if p, err := mango.NewReplayParser(testReplayFilename); err != nil {
		t.Error(err)
	} else if err = p.Initialise(); err != nil {
		t.Error(err)
	} else if packets, err := p.ParseReplay(); err != nil {
		t.Error(err)
	} else {
		for _, packet := range packets {
			fmt.Println(packet.Command)
			if packet.Size < 1000 {
				err = mango.PrintStruct(packet.Message)
				if err != nil {
					t.Error(err)
				}
			} else {
				fmt.Println("Too big to show :(")
			}
			fmt.Println()
		}
	}
}

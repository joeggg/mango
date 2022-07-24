package mango_test

import (
	"fmt"
	"mango"
	"testing"
)

const testReplayFilename = "testdata/test.dem"

func TestSummary(t *testing.T) {
	if p, err := mango.NewReplayParser(testReplayFilename); err != nil {
		t.Error(err)
	} else if err = p.Initialise(); err != nil {
		t.Error(err)
	} else if summary, err := p.GetSummary(); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("Summary packet:\n\n")
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
		fmt.Println("All replay parsed through without errors!")
		fmt.Printf("Sample packets: \n\n")
		for _, packet := range packets {
			fmt.Printf("%s:\n", packet.Command)
			if packet.Size < 1000 {
				if packet.Embed != nil {
					fmt.Printf("Embedded packet! Type: %s\n", packet.Embed.Command)
					err = mango.PrintStruct(packet.Embed.Data)
				} else {
					err = mango.PrintStruct(packet.Message)
				}
				if err != nil {
					t.Error(err)
				}
			} else {
				fmt.Printf("Too big to show :(\n\n")
			}
			fmt.Println()
		}
	}
}

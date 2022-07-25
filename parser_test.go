package mango_test

import (
	"fmt"
	"mango"
	"testing"

	"google.golang.org/protobuf/proto"
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
		mango.PrintStruct(summary)
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
		count := 0
		for _, packet := range packets {
			var toShow proto.Message
			show := false
			if packet.Embed != nil {
				toShow = packet.Embed.Data
				if packet.Embed.Kind != 4 && packet.Embed.Kind != 145 {
					show = true
					count++
					fmt.Printf("%s:\n", packet.Command)
					fmt.Printf("Embedded packet! Type: %s\n", packet.Embed.Command)
				}
			} else {
				fmt.Printf("%s:\n", packet.Command)
				toShow = packet.Message
				show = true
				count++
			}

			if packet.Size < 10000 {
				if show {
					mango.PrintStruct(toShow)
					fmt.Println()
				}
			} else {
				fmt.Printf("Too big to show :(\n\n")
			}
		}
		fmt.Printf("Showing %d packets", count)
	}
}

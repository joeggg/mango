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
		for i := 0; i < 10; i++ {
			var toShow proto.Message
			show := false
			if packets[i].Embed != nil {
				toShow = packets[i].Embed.Data
				if packets[i].Embed.Kind != 4 && packets[i].Embed.Kind != 145 {
					show = true
					count++
					fmt.Printf("%s:\n", packets[i].Command)
					fmt.Printf("Embedded packet! Type: %s\n", packets[i].Embed.Command)
				}
			} else {
				fmt.Printf("%s:\n", packets[i].Command)
				toShow = packets[i].Message
				show = true
				count++
			}

			if packets[i].Size < 10000 {
				if show {
					mango.PrintStruct(toShow)
					fmt.Println()
				}
			} else {
				fmt.Printf("Too big to show :(\n\n")
			}
		}
		fmt.Printf("Showing %d packets\n", count)
	}
}

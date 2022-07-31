package mango_test

import (
	"fmt"
	"mango"
	"mango/gatherers"
	"mango/packet"
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
			if packets[i].Embed != nil {
				fmt.Printf("%s:\n", packets[i].Command)
				fmt.Printf("Embedded packet! Type: %s\n", packets[i].Embed.Command)
				toShow = packets[i].Embed.Data
				count++
			} else {
				fmt.Printf("%s:\n", packets[i].Command)
				toShow = packets[i].Message
				count++
			}

			if toShow != nil {
				if packets[i].Size < 10000 {
					mango.PrintStruct(toShow)
					fmt.Println()
				} else {
					fmt.Printf("Too big to show :(\n\n")
				}
			}
		}
		fmt.Printf("Showing %d packets\n", count)
	}
}

func TestParseWithGatherers(t *testing.T) {
	rp, err := mango.NewReplayParser(testReplayFilename)
	if err != nil {
		t.Error(err)
	}
	cg := gatherers.NewChatGatherer()
	rp.RegisterGatherer(cg)

	if err = rp.Initialise(); err != nil {
		t.Error(err)
	} else if _, err := rp.ParseReplay(); err != nil {
		t.Error(err)
	} else {
		for player, messages := range cg.GetResults().(map[int][]*gatherers.Message) {
			fmt.Printf("%s:\n", packet.Players[player])
			mango.PrintStruct(messages)
		}
	}
}

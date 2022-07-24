package mango_test

import (
	"encoding/base64"
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
		fmt.Println("All replay parsed through without errors!")
		fmt.Println("Sample packets: ")
		for i := 0; i < 5; i++ {
			fmt.Println(packets[i].Command)
			fmt.Println(base64.StdEncoding.EncodeToString(packets[i].RawMessage))
			if packets[i].Size < 1000 {
				err = mango.PrintStruct(packets[i].Message)
				if err != nil {
					t.Error(err)
				}
				if packets[i].Embed == nil {
					continue
				}
				err = mango.PrintStruct(packets[i].Embed.Data)
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

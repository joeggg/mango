package mango_test

import (
	"fmt"
	"mango"
	"testing"
)

func TestSummary(t *testing.T) {
	if p, err := mango.NewReplayParser("example_data/test.dem"); err != nil {
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
	if p, err := mango.NewReplayParser("example_data/test.dem"); err != nil {
		t.Error(err)
	} else if err = p.Initialise(); err != nil {
		t.Error(err)
	} else {

		err = p.ParseReplay()
		if err != nil {
			t.Error(err)
		}
		err = mango.PrintStruct(p)
		if err != nil {
			t.Error(err)
		}
	}
}

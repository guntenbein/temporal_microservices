package domain

import (
	"context"
	"testing"
)

func TestGetActivityName(t *testing.T) {
	if GetActivityName(Activity1) != "Activity1" {
		t.Fatalf("wrong defined name for separate activity function")
	}
	if GetActivityName(ActivityHandler1{}.Activity2) != "Activity2" {
		t.Fatalf("wrong defined name for method of the by value receiver")
	}
	if GetActivityName((&ActivityHandler2{}).Activity3) != "Activity3" {
		t.Fatalf("wrong defined name for method of the pointer receoiver")
	}
}

func Activity1(_ context.Context) {
}

type ActivityHandler1 struct{}

func (ah ActivityHandler1) Activity2(_ context.Context) {
}

type ActivityHandler2 struct{}

func (ah *ActivityHandler2) Activity3(_ context.Context) {
}

package grpc_test

import (
	model "github.com/FlowingSPDG/gotv-plus-go/grpc_client/src"
	"testing"
)

func TestGetMatches(t *testing.T) {
	g := &model.GOTVPLUS{}
	g, err := g.Init("localhost:50055")
	if err != nil {
		t.Fatal(err)
	}
	matches, err := g.GetMatches()
	if err != nil {
		t.Fatal(err)
	}
	match := matches.GetMatch()
	t.Logf("match : %v\n", match)
}

func TestMarkID(t *testing.T) {
	g := &model.GOTVPLUS{}
	g, err := g.Init("localhost:50055")
	if err != nil {
		t.Fatal(err)
	}
	err = g.MarkID("s90132533918272518t1581172258", "MATCH_ID_1") // Mark first match as "MATCH_1". so you can play them by http://localhost:8080/match/id/TOKEN/sync
	if err != nil {
		t.Fatal(err)
	}
}

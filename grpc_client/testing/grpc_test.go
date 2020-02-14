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
	err = g.MarkID("s90132662860548102t1581661145", "MATCH_ID_1") // Mark first match as "MATCH_1". so you can play them by http://localhost:8080/match/id/TOKEN/sync
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMatchByID(t *testing.T) {
	g := &model.GOTVPLUS{}
	g, err := g.Init("localhost:50055")
	if err != nil {
		t.Fatal(err)
	}
	match, err := g.GetMatchByID("MATCH_ID_1")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("match : %v\n", match)
}

func TestGetMatchByToken(t *testing.T) {
	g := &model.GOTVPLUS{}
	g, err := g.Init("localhost:50055")
	if err != nil {
		t.Fatal(err)
	}
	match, err := g.GetMatchByToken("s90132662860548102t1581661145")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("match : %v\n", match)
}

func TestDeleteMatchByID(t *testing.T) {
	g := &model.GOTVPLUS{}
	g, err := g.Init("localhost:50055")
	if err != nil {
		t.Fatal(err)
	}
	err = g.DeleteMatchByID("MATCH_ID_1")
	if err != nil {
		t.Fatal(err)
	}
}
func TestDeleteMatchByToken(t *testing.T) {
	g := &model.GOTVPLUS{}
	g, err := g.Init("localhost:50055")
	if err != nil {
		t.Fatal(err)
	}
	err = g.DeleteMatchByToken("s90132533918272518t1581172258")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSaveMatchToFileByToken(t *testing.T) {
	g := &model.GOTVPLUS{}
	g, err := g.Init("localhost:50055")
	if err != nil {
		t.Fatal(err)
	}
	err = g.SaveMatchToFileByToken("s90132662860548102t1581661145", "matches")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSaveMatchToFileByID(t *testing.T) {
	g := &model.GOTVPLUS{}
	g, err := g.Init("localhost:50055")
	if err != nil {
		t.Fatal(err)
	}
	err = g.SaveMatchToFileByToken("MATCH_ID_1", "matches")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadMatchFromFile(t *testing.T) {
	g := &model.GOTVPLUS{}
	g, err := g.Init("localhost:50055")
	if err != nil {
		t.Fatal(err)
	}
	err = g.LoadMatchFromFile("s90132662860548102t1581661145")
	if err != nil {
		t.Fatal(err)
	}
}

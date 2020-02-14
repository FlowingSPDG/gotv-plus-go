package src

import (
	"context"
	pb "github.com/FlowingSPDG/gotv-plus-go/server/src/grpc/protogen"
	"google.golang.org/grpc"
	"log"
)

type GOTVPLUS struct {
	Addr   string
	Client *pb.GOTV_PlusClient
	conn   *grpc.ClientConn
}

func (g *GOTVPLUS) Init(addr string) (*GOTVPLUS, error) {
	g = &GOTVPLUS{}
	g.Addr = addr
	conn, err := grpc.Dial(g.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("client connection error:", err)
		return nil, err
	}
	g.conn = conn
	client := pb.NewGOTV_PlusClient(conn)
	g.Client = &client
	return g, nil
}

func (g *GOTVPLUS) Disconnect() error { // ID,IP,err
	log.Println("Disconnected Talos.")
	err := g.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (g *GOTVPLUS) GetMatches() (*pb.GetMatchesReply, error) {
	c := *g.Client
	option := &pb.GetMatchesRequest{}
	matches, err := c.GetMatches(context.TODO(), option)
	if err != nil {
		return nil, err
	}
	log.Printf("Matches : %v\n", *matches)
	return matches, nil
}

func (g *GOTVPLUS) GetMatchByID(id string) (*pb.Match, error) {
	c := *g.Client
	option := &pb.GetMatchRequest{
		Ids: &pb.GetMatchRequest_Id{
			Id: id,
		},
	}
	matches, err := c.GetMatch(context.TODO(), option)
	if err != nil {
		return nil, err
	}
	log.Printf("Match : %v\n", *matches)
	return matches, nil
}

func (g *GOTVPLUS) GetMatchByToken(token string) (*pb.Match, error) {
	c := *g.Client
	option := &pb.GetMatchRequest{
		Ids: &pb.GetMatchRequest_Token{
			Token: token,
		},
	}
	matches, err := c.GetMatch(context.TODO(), option)
	if err != nil {
		return nil, err
	}
	log.Printf("Match : %v\n", *matches)
	return matches, nil
}

func (g *GOTVPLUS) DeleteMatchByID(id string) error {
	c := *g.Client
	option := &pb.DeleteMatchRequest{
		Ids: &pb.DeleteMatchRequest_Id{
			Id: id,
		},
	}
	_, err := c.DeleteMatch(context.TODO(), option)
	if err != nil {
		return err
	}
	return nil
}

func (g *GOTVPLUS) DeleteMatchByToken(token string) error {
	c := *g.Client
	option := &pb.DeleteMatchRequest{
		Ids: &pb.DeleteMatchRequest_Token{
			Token: token,
		},
	}
	_, err := c.DeleteMatch(context.TODO(), option)
	if err != nil {
		return err
	}
	return nil
}

func (g *GOTVPLUS) MarkID(token string, id string) error {
	c := *g.Client
	option := &pb.MarkIDRequest{
		Token: token,
		Id:    id,
	}
	_, err := c.MarkID(context.TODO(), option)
	if err != nil {
		return err
	}
	return nil
}

func (g *GOTVPLUS) SaveMatchToFileByID(id string, path string) error {
	c := *g.Client
	option := &pb.SaveMatchToFileRequest{
		Ids: &pb.SaveMatchToFileRequest_Id{
			Id: id,
		},
		Path: path,
	}
	_, err := c.SaveMatchToFile(context.TODO(), option)
	if err != nil {
		return err
	}
	return nil
}

func (g *GOTVPLUS) SaveMatchToFileByToken(token string, path string) error {
	c := *g.Client
	option := &pb.SaveMatchToFileRequest{
		Ids: &pb.SaveMatchToFileRequest_Token{
			Token: token,
		},
		Path: path,
	}
	_, err := c.SaveMatchToFile(context.TODO(), option)
	if err != nil {
		return err
	}
	return nil
}

func (g *GOTVPLUS) LoadMatchFromFile(token string) error {
	c := *g.Client
	option := &pb.LoadMatchFromFileRequest{
		Token: token,
	}
	_, err := c.LoadMatchFromFile(context.TODO(), option)
	if err != nil {
		return err
	}
	return nil
}

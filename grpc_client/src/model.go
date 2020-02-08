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

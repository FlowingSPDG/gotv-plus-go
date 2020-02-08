package grpc

import (
	"context"
	pb "github.com/FlowingSPDG/gotv-plus-go/server/src/grpc/protogen"
	"github.com/FlowingSPDG/gotv-plus-go/server/src/handlers"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s server) GetMatches(ctx context.Context, message *pb.GetMatchesRequest) (*pb.GetMatchesReply, error) {
	log.Println("GetMatches")

	matches, err := handlers.Matches.GetAll()
	if err != nil {
		return &pb.GetMatchesReply{
			Match:        nil,
			Error:        true,
			Errormessage: err.Error(),
		}, nil
	}
	pbmatches := make([]*pb.Match, 0, len(matches))
	for _, v := range matches {
		pbmatches = append(pbmatches, &pb.Match{
			Token: v.Token,
			Id:    v.ID,
		})
	}
	return &pb.GetMatchesReply{
		Match:        pbmatches,
		Error:        false,
		Errormessage: "",
	}, nil
}

func (s server) MarkID(ctx context.Context, message *pb.MarkIDRequest) (*pb.MarkIDReply, error) {
	log.Println("MarkID")
	m, err := handlers.Matches.GetMatchByToken(message.GetToken())
	if err != nil {
		return &pb.MarkIDReply{
			Error:        true,
			Errormessage: err.Error(),
		}, nil
	}
	err = m.TagID(message.GetId())
	if err != nil {
		return &pb.MarkIDReply{
			Error:        true,
			Errormessage: err.Error(),
		}, nil
	}
	return &pb.MarkIDReply{
		Error:        false,
		Errormessage: "",
	}, nil
}

// StartGRPC Starts gRPC API Server on specified ADDR
func StartGRPC(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	// create grpc server
	s := grpc.NewServer()
	pb.RegisterGOTV_PlusServer(s, server{})

	log.Printf("Listening on : %s", addr)

	// and start...
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}

package grpc

import (
	"context"
	"fmt"
	pb "github.com/FlowingSPDG/gotv-plus-go/server/src/grpc/protogen"
	"github.com/FlowingSPDG/gotv-plus-go/server/src/handlers"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s server) GetMatches(ctx context.Context, message *pb.GetMatchesRequest) (*pb.GetMatchesReply, error) {
	log.Println("[gRPC] GetMatches")

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

func (s server) GetMatch(ctx context.Context, message *pb.GetMatchRequest) (*pb.Match, error) {
	log.Println("[gRPC] GetMatch")

	ids := message.GetIds()
	switch i := ids.(type) {
	case *pb.GetMatchRequest_Id:
		match, err := handlers.Matches.GetMatchByID(i.Id)
		if err != nil {
			return nil, err
		}
		return &pb.Match{
			Token: match.Token,
			Id:    match.ID,
		}, nil
	case *pb.GetMatchRequest_Token:
		match, err := handlers.Matches.GetMatchByToken(i.Token)
		if err != nil {
			return nil, err
		}
		return &pb.Match{
			Token: match.Token,
			Id:    match.ID,
		}, nil
	}
	return nil, fmt.Errorf("Something went wrong")
}

func (s server) DeleteMatch(ctx context.Context, message *pb.DeleteMatchRequest) (*pb.DeleteMatchReply, error) {
	log.Println("[gRPC] DeleteMatch")

	ids := message.GetIds()
	var match *handlers.Match
	var err error
	switch i := ids.(type) {
	case *pb.DeleteMatchRequest_Id:
		match, err = handlers.Matches.GetMatchByID(i.Id)
	case *pb.DeleteMatchRequest_Token:
		match, err = handlers.Matches.GetMatchByToken(i.Token)
	}
	if err != nil {
		return nil, err
	}
	err = handlers.Matches.Delete(match)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteMatchReply{
		Error:        false,
		Errormessage: "",
	}, nil
}

func (s server) MarkID(ctx context.Context, message *pb.MarkIDRequest) (*pb.MarkIDReply, error) {
	log.Println("[gRPC] MarkID")
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

	log.Printf("[gRPC] Listening on : %s", addr)

	// and start...
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}

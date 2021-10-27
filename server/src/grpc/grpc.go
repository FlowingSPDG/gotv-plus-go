package grpc

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/FlowingSPDG/gotv-plus-go/server/src/grpc/protogen"
	"github.com/FlowingSPDG/gotv-plus-go/server/src/handlers"
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

	var match *handlers.Match
	var err error
	ids := message.GetIds()
	switch i := ids.(type) {
	case *pb.GetMatchRequest_Id:
		match, err = handlers.Matches.GetMatchByID(i.Id)
		if err != nil {
			return nil, err
		}
	case *pb.GetMatchRequest_Token:
		match, err = handlers.Matches.GetMatchByToken(i.Token)
		if err != nil {
			return nil, err
		}
	}
	log.Printf("[gRPC] Match : Token[%s] Latest Fragment[%d]\n", match.Token, match.Latest)
	return &pb.Match{
		Token: match.Token,
		Id:    match.ID,
	}, nil
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

func (s server) SaveMatchToFile(ctx context.Context, message *pb.SaveMatchToFileRequest) (*pb.SaveMatchToFileReply, error) {
	log.Println("[gRPC] SaveMatchToFile")
	ids := message.GetIds()
	var match *handlers.Match
	var err error
	switch i := ids.(type) {
	case *pb.SaveMatchToFileRequest_Id:
		match, err = handlers.Matches.GetMatchByID(i.Id)
	case *pb.SaveMatchToFileRequest_Token:
		match, err = handlers.Matches.GetMatchByToken(i.Token)
	}
	if err != nil {
		log.Printf("ERR on saving match : %v\n", err)
		return nil, err
	}
	log.Printf("Match : %v\n", match)
	err = match.SaveMatchToFile(message.GetPath())
	if err != nil {
		log.Printf("ERR on saving match : %v\n", err)
		return nil, err
	}
	return &pb.SaveMatchToFileReply{
		Error:        false,
		Errormessage: "",
	}, nil
}

func (s server) LoadMatchFromFile(ctx context.Context, message *pb.LoadMatchFromFileRequest) (*pb.LoadMatchFromFileReply, error) {
	log.Println("[gRPC] LoadMatchFromFile")
	token := message.GetToken()
	token, err := handlers.Matches.LoadMatchFromFile(token)
	if err != nil {
		log.Printf("ERR on saving match : %v\n", err)
		return &pb.LoadMatchFromFileReply{
			Error:        true,
			Errormessage: err.Error(),
		}, nil
	}
	return &pb.LoadMatchFromFileReply{
		Token:        token,
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

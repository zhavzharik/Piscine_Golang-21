package main

import (
	"github.com/google/uuid"
	"gonum.org/v1/gonum/stat/distuv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"net"
	pb "nlo/nlo"
	"time"
)

const (
	port = ":50051"
)

type nloServer struct {
	pb.UnimplementedGetRespServer
}

func (s *nloServer) Do(in *pb.Request, stream pb.GetResp_DoServer) error {
	rand.Seed(time.Now().Unix())
	mean := rand.Float64()*21 - 10
	std := rand.Float64()*1.2 + 0.3
	var res pb.Response
	res.SessionId = uuid.New().String()
	log.Printf("Session_ID: %s, Mean: %f, STD: %f", res.GetSessionId(), mean, std)
	for {
		dist := distuv.Normal{
			Mu:    mean,
			Sigma: std,
		}
		res.Frequency = dist.Rand()
		res.Timestamp = timestamppb.New(time.Now().UTC())
		if err := stream.Send(&res); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	pb.RegisterGetRespServer(s, &nloServer{})
	s.Serve(lis)
}

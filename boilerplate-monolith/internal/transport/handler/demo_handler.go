package handler

import (
	"codegen/api/pb"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type demoHandler struct {
	pb.UnimplementedDemoServiceServer
}

func NewDemoHandler() *demoHandler {
	return &demoHandler{}
}

var demos []*pb.Demo = make([]*pb.Demo, 3)

func init() {
	demos[0] = &pb.Demo{Id: 1, Name: "Bir"}
	demos[1] = &pb.Demo{Id: 2, Name: "İki"}
	demos[2] = &pb.Demo{Id: 3, Name: "Üç"}
}

func (h *demoHandler) ListDemos(ctx context.Context, _ *emptypb.Empty) (*pb.ListDemosResponse, error) {
	return &pb.ListDemosResponse{Demos: demos}, nil
}

func (h *demoHandler) GetDemo(ctx context.Context, req *pb.DemoIdentifier) (*pb.GetDemoResponse, error) {
	demo := demos[req.Id]
	return &pb.GetDemoResponse{Demo: demo}, nil
}

func (h *demoHandler) CreateDemo(ctx context.Context, req *pb.CreateDemoRequest) (*pb.DemoIdentifier, error) {
	newId := demos[len(demos)-1].Id + 1
	newDemo := &pb.Demo{Id: newId, Name: req.Name, Description: req.Description}

	demos = append(demos, newDemo)

	return &pb.DemoIdentifier{Id: newId}, nil
}

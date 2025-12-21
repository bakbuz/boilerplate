package handler

import (
	demov1 "codegen/api/gen/demo/v1"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type demoHandler struct {
	demov1.UnimplementedDemoServiceServer
}

func NewDemoHandler() *demoHandler {
	return &demoHandler{}
}

var demos []*demov1.Demo = make([]*demov1.Demo, 3)

func init() {
	demos[0] = &demov1.Demo{Id: 1, Name: "Bir"}
	demos[1] = &demov1.Demo{Id: 2, Name: "İki"}
	demos[2] = &demov1.Demo{Id: 3, Name: "Üç"}
}

func getDemoById(id int32) *demov1.Demo {
	for _, demo := range demos {
		if demo.Id == id {
			return demo
		}
	}
	return nil
}

func (h *demoHandler) List(context.Context, *demov1.ListDemosRequest) (*demov1.ListDemosResponse, error) {
	return &demov1.ListDemosResponse{Items: demos}, nil
}

func (h *demoHandler) Get(ctx context.Context, req *demov1.GetDemoRequest) (*demov1.Demo, error) {
	demo := getDemoById(req.Id)
	return demo, nil
}

func (h *demoHandler) Create(ctx context.Context, req *demov1.CreateDemoRequest) (*demov1.Demo, error) {
	newId := demos[len(demos)-1].Id + 1
	newDemo := &demov1.Demo{Id: newId, Name: req.Name, Description: req.Description}

	demos = append(demos, newDemo)

	return newDemo, nil
}

func (h *demoHandler) Update(ctx context.Context, req *demov1.UpdateDemoRequest) (*demov1.Demo, error) {
	demo := getDemoById(req.Id)

	demo.Name = req.Name
	demo.Description = req.Description

	return demo, nil
}
func (h *demoHandler) Delete(ctx context.Context, req *demov1.DeleteDemoRequest) (*emptypb.Empty, error) {

	// 1. Önce silinecek elemanın index'ini bul
	index := -1
	for i, d := range demos {
		if d.Id == req.Id {
			index = i
			break
		}
	}

	// 2. Eğer bulunduysa silme işlemini yap
	if index != -1 {
		// [0...index-1] + [index+1...son]
		demos = append(demos[:index], demos[index+1:]...)
	}

	return &emptypb.Empty{}, nil
}

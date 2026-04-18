package services

import (
	"context"
	"mangahub/proto/sample"
)

type GRPCSampleService struct {
	sample.UnimplementedSampleServiceServer
}

func (s *GRPCSampleService) SampleMethod(ctx context.Context, req *sample.SampleRequest) (*sample.SampleResponse, error) {
	// Xử lý logic của phương thức ở đây
	response := &sample.SampleResponse{
		Message: "Hello, " + req.Name + "! This is a sample response.",
	}
	return response, nil
}	
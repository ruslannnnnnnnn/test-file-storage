package grpc_controller

import (
	"context"

	pb "github.com/ruslannnnnnnnn/test-file-storage/api/gen/go/service/v1"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcFileController struct {
	fileService service.IFileService
	pb.UnimplementedFileServiceServer
}

func (g *GrpcFileController) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (g *GrpcFileController) Download(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (g *GrpcFileController) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFiles not implemented")
}

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

func NewGrpcFileController(fileService service.IFileService) *GrpcFileController {
	return &GrpcFileController{fileService: fileService}
}

func (g *GrpcFileController) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (g *GrpcFileController) Download(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (g *GrpcFileController) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	listFilesResponse, err := g.fileService.ListFiles()
	if err != nil {
		return nil, err
	}

	var fileInfo []*pb.FileInfo

	for _, file := range listFilesResponse.Files {
		fileInfo = append(fileInfo, &pb.FileInfo{
			Id:        file.Id,
			Filename:  file.Name,
			CreatedAt: file.CreatedAt.String(),
			UpdatedAt: file.UpdatedAt.String(),
		})
	}

	return &pb.ListFilesResponse{Files: fileInfo}, nil
}

package grpc_controller

import (
	"context"
	"time"

	pb "github.com/ruslannnnnnnnn/test-file-storage/api/gen/go/service/v1"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/api/grpc_controller/helper"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/common"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/service"
)

type GrpcFileController struct {
	fileService service.IFileService
	pb.UnimplementedFileServiceServer
}

func NewGrpcFileController(fileService service.IFileService) *GrpcFileController {
	return &GrpcFileController{fileService: fileService}
}

func (g *GrpcFileController) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	fileId, err := g.fileService.Upload(common.UploadRequest{FileContent: req.GetData(), FileName: req.GetFilename()})
	if err != nil {
		return nil, helper.ToGrpcError(err)
	}

	return &pb.UploadResponse{Id: fileId}, nil
}

func (g *GrpcFileController) Download(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	result, err := g.fileService.Download(common.DownloadRequest{FileId: req.GetId()})
	if err != nil {
		return nil, helper.ToGrpcError(err)
	}

	return &pb.DownloadResponse{Data: result.FileContent, Filename: result.FileName}, nil
}

func (g *GrpcFileController) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	listFilesResponse, err := g.fileService.ListFiles()
	if err != nil {
		return nil, helper.ToGrpcError(err)
	}

	var fileInfo []*pb.FileInfo

	for _, file := range listFilesResponse.Files {
		fileInfo = append(fileInfo, &pb.FileInfo{
			Id:        file.Id,
			Filename:  file.Name,
			CreatedAt: file.CreatedAt.Format(time.RFC3339),
			UpdatedAt: file.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListFilesResponse{Files: fileInfo}, nil
}

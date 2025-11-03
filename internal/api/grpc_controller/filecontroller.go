package grpc_controller

import (
	"context"
	"io"
	"sync"
	"time"

	pb "github.com/ruslannnnnnnnn/test-file-storage/api/gen/go/service/v1"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/api/grpc_controller/helper"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/service"
	"google.golang.org/grpc"
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

// Upload запрос на выгрузку файла к нам, надо получить и записать файл
func (g *GrpcFileController) Upload(bidiStream grpc.BidiStreamingServer[pb.UploadRequest, pb.UploadResponse]) error {

	r, w := io.Pipe()

	fileNameChan := make(chan string)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer w.Close()
		firstChunk := true

		for {
			select {
			case <-bidiStream.Context().Done():
				return
			default:
			}

			req, err := bidiStream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				w.CloseWithError(err)
				return
			}

			if firstChunk {
				select {
				case fileNameChan <- req.Filename:
				case <-bidiStream.Context().Done():
					return
				}
				firstChunk = false
			}

			_, err = w.Write(req.Data)
			if err != nil {
				w.CloseWithError(err)
				return
			}

		}
	}()

	fileName := <-fileNameChan
	fileIdCh := make(chan string)

	wg.Add(1)
	go func() {
		defer wg.Done()

		g.fileService.Upload(bidiStream.Context(), fileName, fileIdCh, r)
	}()

	// надо клиенту узнать id файла как только будет сгенерирован, чтобы не ждать полной выгрузки файла
	fileId := <-fileIdCh
	bidiStream.Send(&pb.UploadResponse{
		Id: fileId,
	})

	wg.Wait()

	return nil
}

// Download запрос на скачивание файла к нам, надо отдать файл
func (g *GrpcFileController) Download(req *pb.DownloadRequest, stream grpc.ServerStreamingServer[pb.DownloadResponse]) error {

	fileName, downloadResponseReader, err := g.fileService.Download(stream.Context(), req.Id)

	if err != nil {
		return helper.ToGrpcError(err)
	}

	fileChunk := make([]byte, 1024)

	for {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Stream ended")
		default:
			n, err := downloadResponseReader.Read(fileChunk)
			if err == io.EOF {
				return nil
			}

			err = stream.Send(&pb.DownloadResponse{
				Filename: fileName,
				Data:     fileChunk[:n],
			})

			if err != nil {
				return status.Error(codes.Canceled, "Stream ended")
			}
		}
	}

}

func (g *GrpcFileController) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	listFilesResponse, err := g.fileService.ListFiles(ctx)
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

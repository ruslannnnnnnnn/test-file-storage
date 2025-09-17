package limiter_test

import (
	"context"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	pb "github.com/ruslannnnnnnnn/test-file-storage/api/gen/go/service/v1"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/api/grpc_controller"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/common"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/limiter"
	"google.golang.org/grpc"
	grpc_codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// fakeFileService реализует IFileService
type fakeFileService struct{}

func (f *fakeFileService) Upload(req common.UploadRequest) (string, error) {
	time.Sleep(10 * time.Millisecond)
	return "file-id-123", nil
}

func (f *fakeFileService) Download(req common.DownloadRequest) (common.DownloadResponse, error) {
	time.Sleep(500 * time.Millisecond) // имитация долгого запроса
	return common.DownloadResponse{FileContent: []byte("test")}, nil
}

func (f *fakeFileService) ListFiles() (common.ListFilesResponse, error) {
	time.Sleep(10 * time.Millisecond)
	return common.ListFilesResponse{}, nil
}

func startTestGRPCServer() (addr string, stopFunc func()) {
	lis, _ := net.Listen("tcp", ":0") // случайный порт

	limits := map[string]int{
		"/service.v1.FileService/Download":  10,
		"/service.v1.FileService/ListFiles": 5,
	}
	reqLimiter := limiter.NewLimiter(limits, time.Minute)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(reqLimiter.UnaryInterceptor),
		grpc.StreamInterceptor(reqLimiter.StreamInterceptor),
	)

	ctrl := grpc_controller.NewGrpcFileController(&fakeFileService{})
	pb.RegisterFileServiceServer(grpcServer, ctrl)

	go grpcServer.Serve(lis)

	return lis.Addr().String(), func() {
		reqLimiter.Stop()
		grpcServer.Stop()
	}
}

func TestDownloadLimit(t *testing.T) {
	addr, stop := startTestGRPCServer()
	defer stop()

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewFileServiceClient(conn)

	concurrent := 15
	var wg sync.WaitGroup
	wg.Add(concurrent)

	results := make([]error, concurrent)

	for i := 0; i < concurrent; i++ {
		go func(idx int) {
			defer wg.Done()
			_, err := client.Download(context.Background(), &pb.DownloadRequest{Id: "1"})
			results[idx] = err
		}(i)
	}

	wg.Wait()

	successCount := 0
	for _, r := range results {
		if r == nil {
			successCount++
		} else {
			st, ok := status.FromError(r)
			if !ok {
				t.Fatalf("unexpected error type: %v", r)
			}
			if st.Code() != grpc_codes.ResourceExhausted {
				t.Fatalf("unexpected grpc code: %v", st.Code())
			}
		}
	}

	if successCount != 10 { // лимит
		t.Fatalf("expected 10 successful downloads, got %d", successCount)
	}

	log.Println("TestDownloadLimit passed")
}

func TestListFilesLimit(t *testing.T) {
	addr, stop := startTestGRPCServer()
	defer stop()

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewFileServiceClient(conn)

	concurrent := 10
	var wg sync.WaitGroup
	wg.Add(concurrent)

	results := make([]error, concurrent)

	for i := 0; i < concurrent; i++ {
		go func(idx int) {
			defer wg.Done()
			_, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
			results[idx] = err
		}(i)
	}

	wg.Wait()

	successCount := 0
	for _, r := range results {
		if r == nil {
			successCount++
		} else {
			st, ok := status.FromError(r)
			if !ok {
				t.Fatalf("unexpected error type: %v", r)
			}
			if st.Code() != grpc_codes.ResourceExhausted {
				t.Fatalf("unexpected grpc code: %v", st.Code())
			}
		}
	}

	if successCount != 5 { // лимит
		t.Fatalf("expected 5 successful ListFiles, got %d", successCount)
	}

	log.Println("TestListFilesLimit passed")
}

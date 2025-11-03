package helper

import (
	"errors"

	"github.com/ruslannnnnnnnn/test-file-storage/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGrpcError(err error) error {

	if errors.Is(err, service.InvalidUUidError) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, service.FileNotFoundError) {
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Internal, "internal error")
}

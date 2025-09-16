package helper

import (
	"reflect"

	"github.com/ruslannnnnnnnn/test-file-storage/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGrpcError маппинг сервисных ошибок в grpcшные
func ToGrpcError(err error) error {

	if ErrorInstanceOf(err, service.InvalidUUidError{}) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if ErrorInstanceOf(err, service.FileNotFoundError{}) {
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Internal, "internal error")
}

// ErrorInstanceOf имеют ли оба параметра один и тот же тип, хотел сделать аналог instanceof из php
func ErrorInstanceOf(instance interface{}, of interface{}) bool {
	if instance == nil {
		return false
	}

	// Типы через reflect
	instType := reflect.TypeOf(instance)
	ofType := reflect.TypeOf(of)

	// Сравниваем не только указатели, но и базовые типы
	if instType == ofType {
		return true
	}
	if instType.Kind() == reflect.Ptr && instType.Elem() == ofType {
		return true
	}
	if ofType.Kind() == reflect.Ptr && ofType.Elem() == instType {
		return true
	}

	return false
}

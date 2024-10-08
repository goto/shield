package v1beta1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goto/shield/pkg/errors"
)

// HTTP Codes defined here:
// https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/errors.go#L36
var (
	ErrInternalServer        = errors.New("internal server error")
	ErrBadRequest            = errors.New("invalid syntax in body")
	ErrConflictRequest       = errors.New("already exist")
	ErrRequestBodyValidation = errors.New("invalid format for field(s)")
	ErrEmptyEmailID          = errors.New("email id is empty")
	ErrUnsupportes           = errors.New("the operation is not implemented or is not supported/enabled in this service.")

	grpcInternalServerError = status.Errorf(codes.Internal, ErrInternalServer.Error())
	grpcConflictError       = status.Errorf(codes.AlreadyExists, ErrConflictRequest.Error())
	grpcBadBodyError        = status.Error(codes.InvalidArgument, ErrBadRequest.Error())
	grpcUnauthenticated     = status.Error(codes.Unauthenticated, errors.ErrUnauthenticated.Error())
	grpcPermissionDenied    = status.Error(codes.PermissionDenied, errors.ErrForbidden.Error())
	grpcUnsupportedError    = status.Error(codes.Unimplemented, ErrUnsupportes.Error())
)

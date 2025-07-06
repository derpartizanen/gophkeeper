package grpc

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Craft gRPC Status with additional details regarding bad request's fields.
func composeBadRequestError(details *errdetails.BadRequest) *status.Status {
	st := status.New(codes.InvalidArgument, "invalid request")

	st, err := st.WithDetails(details)
	if err != nil {
		return status.New(codes.Internal, err.Error())
	}

	return st
}

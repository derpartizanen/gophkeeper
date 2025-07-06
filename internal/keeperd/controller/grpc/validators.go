package grpc

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/derpartizanen/gophkeeper/proto"
)

const (
	MissingField = "not set"

	DefaultMaxUsernameLength   = 128
	DefaultMaxSecretNameLength = 256

	DefaultMetadataLimit = 2 * 1024 * 1024

	DefaultDataLimit = 4 * 1024 * 1024
)

// validateUsername validates provided username.
func validateUsername(name string) (string, bool) {
	if name == "" {
		return MissingField, false
	}

	if len(name) > DefaultMaxUsernameLength {
		return fmt.Sprintf("should be <= %d characters", DefaultMaxUsernameLength), false
	}

	return "", true
}

// validateSecurityKey validates provided security key.
func validateSecurityKey(key string) (string, bool) {
	if key == "" {
		return MissingField, false
	}

	return "", true
}

// validateCredentials validates provided credentials.
func validateCredentials(username, key string) (*errdetails.BadRequest, bool) {
	br := &errdetails.BadRequest{}

	if reason, ok := validateUsername(username); !ok {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "username",
			Description: reason,
		}

		br.FieldViolations = append(br.FieldViolations, v)
	}

	if reason, ok := validateSecurityKey(key); !ok {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "security_key",
			Description: reason,
		}

		br.FieldViolations = append(br.FieldViolations, v)
	}

	if len(br.FieldViolations) == 0 {
		return nil, true
	}

	return br, false
}

// validateSecretName validates provided secret name.
func validateSecretName(name string) (string, bool) {
	if name == "" {
		return MissingField, false
	}

	if len(name) > DefaultMaxSecretNameLength {
		return fmt.Sprintf("should be <= %d characters", DefaultMaxSecretNameLength), false
	}

	return "", true
}

// validateMetadata validates provided metadata.
func validateMetadata(metadata []byte) (string, bool) {
	if len(metadata) > DefaultMetadataLimit {
		return fmt.Sprintf("should be <= %d characters", DefaultMetadataLimit), false
	}

	return "", true
}

// validateSecretData validates provided secret data.
func validateSecretData(data []byte) (string, bool) {
	if len(data) == 0 {
		return MissingField, false
	}

	if len(data) > DefaultDataLimit {
		return fmt.Sprintf("should be <= %d characters", DefaultDataLimit), false
	}

	return "", true
}

// validateCreateSecretReq validates goph.validateCreateSecretReq.
func validateCreateSecretReq(
	req *proto.CreateSecretRequest,
) (*errdetails.BadRequest, bool) {
	br := &errdetails.BadRequest{}

	if reason, ok := validateSecretName(req.GetName()); !ok {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "name",
			Description: reason,
		}

		br.FieldViolations = append(br.FieldViolations, v)
	}

	if reason, ok := validateMetadata(req.GetMetadata()); !ok {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "metadata",
			Description: reason,
		}

		br.FieldViolations = append(br.FieldViolations, v)
	}

	if reason, ok := validateSecretData(req.GetData()); !ok {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "data",
			Description: reason,
		}

		br.FieldViolations = append(br.FieldViolations, v)
	}

	if len(br.FieldViolations) == 0 {
		return nil, true
	}

	return br, false
}

// validateUpdateSecretReq validates goph.validateUpdateSecretReq.
func validateUpdateSecretReq(
	req *proto.UpdateSecretRequest,
) (uuid.UUID, *errdetails.BadRequest) {
	var id uuid.UUID

	br := &errdetails.BadRequest{}

	mask := req.GetUpdateMask()
	if mask == nil || len(mask.GetPaths()) == 0 {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "update_mask",
			Description: MissingField,
		}

		br.FieldViolations = append(br.FieldViolations, v)

		return id, br
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "id",
			Description: err.Error(),
		}

		br.FieldViolations = append(br.FieldViolations, v)
	}

	for _, field := range mask.GetPaths() {
		var (
			ok     bool
			reason string
		)

		switch field {
		case "name":
			reason, ok = validateSecretName(req.GetName())

		case "metadata":
			reason, ok = validateMetadata(req.GetMetadata())

		case "data":
			reason, ok = validateSecretData(req.GetData())
		}

		if !ok {
			v := &errdetails.BadRequest_FieldViolation{
				Field:       field,
				Description: reason,
			}

			br.FieldViolations = append(br.FieldViolations, v)
		}
	}

	if len(br.FieldViolations) == 0 {
		return id, nil
	}

	return id, br
}

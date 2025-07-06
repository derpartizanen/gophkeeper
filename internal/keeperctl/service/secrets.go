package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/encryption"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/repo"
	p "github.com/derpartizanen/gophkeeper/proto"
)

var _ Secrets = (*SecretsService)(nil)

var ErrKindMismatch = errors.New("secret kind doesn't match")

// SecretsService contains business logic related to secrets management.
type SecretsService struct {
	key         encryption.Key
	secretsRepo repo.Secrets
}

// NewSecretsService create and initializes new SecretsService object.
func NewSecretsService(key encryption.Key, secrets repo.Secrets) *SecretsService {
	return &SecretsService{key, secrets}
}

// push is low level function sending generic secret creation message to keeper.
func (s *SecretsService) push(
	ctx context.Context,
	token string,
	name string,
	kind p.DataKind,
	description string,
	data proto.Message,
) (uuid.UUID, error) {
	var id uuid.UUID

	rawData, err := proto.Marshal(data)
	if err != nil {
		return id, fmt.Errorf("SecretsService - push - proto.Marshal: %w", err)
	}

	encData, err := s.key.Encrypt(rawData)
	if err != nil {
		return id, fmt.Errorf("SecretsService - push - uc.key.Encrypt(data): %w", err)
	}

	encDescription, err := s.key.Encrypt([]byte(description))
	if err != nil {
		return id, fmt.Errorf("SecretsService - push - uc.key.Encrypt(description): %w", err)
	}

	id, err = s.secretsRepo.Push(ctx, token, name, kind, encDescription, encData)
	if err != nil {
		return id, fmt.Errorf("SecretsService - push - uc.secretsRepo.Push: %w", err)
	}

	return id, nil
}

// PushBinary creates new secret with arbitrary binary data.
func (s *SecretsService) PushBinary(
	ctx context.Context,
	token, name, description string,
	binary []byte,
) (uuid.UUID, error) {
	data := &p.Binary{
		Binary: binary,
	}

	return s.push(ctx, token, name, p.DataKind_BINARY, description, data)
}

// PushCard creates new secret containing bank card data.
func (s *SecretsService) PushCard(
	ctx context.Context,
	token, name, description string,
	number, expiration, holder string,
	cvv int32,
) (uuid.UUID, error) {
	data := &p.Card{
		Number:     number,
		Expiration: expiration,
		Holder:     holder,
		Cvv:        cvv,
	}

	return s.push(ctx, token, name, p.DataKind_CARD, description, data)
}

// PushCreds creates new secret containing credentials.
func (s *SecretsService) PushCreds(
	ctx context.Context,
	token, name, description, login, password string,
) (uuid.UUID, error) {
	data := &p.Credentials{
		Login:    login,
		Password: password,
	}

	return s.push(ctx, token, name, p.DataKind_CREDENTIALS, description, data)
}

// PushText creates new secret with arbitrary text.
func (s *SecretsService) PushText(
	ctx context.Context,
	token, name, description, text string,
) (uuid.UUID, error) {
	data := &p.Text{
		Text: text,
	}

	return s.push(ctx, token, name, p.DataKind_TEXT, description, data)
}

// List returns list of user's secrets.
// All sensitive parts are decrypted.
func (s *SecretsService) List(ctx context.Context, token string) ([]*p.Secret, error) {
	data, err := s.secretsRepo.List(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("SecretsService - List - uc.secretsRepo.List: %w", err)
	}

	for i, val := range data {
		data[i].Metadata, err = s.key.Decrypt(val.GetMetadata())
		if err != nil {
			return nil, fmt.Errorf("SecretsService - List - uc.key.Decrypt: %w", err)
		}
	}

	return data, nil
}

// update is low level function sending generic secret update message to keeper.
func (s *SecretsService) update(
	ctx context.Context,
	token string,
	id uuid.UUID,
	name string,
	description string,
	noDescription bool,
	data proto.Message,
) error {
	var encData []byte

	if data != nil && !reflect.ValueOf(data).IsNil() {
		rawData, err := proto.Marshal(data)
		if err != nil {
			return fmt.Errorf("SecretsService - update - proto.Marshal: %w", err)
		}

		encData, err = s.key.Encrypt(rawData)
		if err != nil {
			return fmt.Errorf("SecretsService - update - uc.key.Encrypt(data): %w", err)
		}
	}

	encDescription, err := s.key.Encrypt([]byte(description))
	if err != nil {
		return fmt.Errorf("SecretsService - update - uc.key.Encrypt(description): %w", err)
	}

	if err = s.secretsRepo.Update(
		ctx,
		token,
		id,
		name,
		encDescription,
		noDescription,
		encData,
	); err != nil {
		return fmt.Errorf("SecretsService - update - uc.secretsRepo.Update: %w", err)
	}

	return nil
}

// EditBinary changes parameters of stored binary secret.
func (s *SecretsService) EditBinary(
	ctx context.Context,
	token string,
	id uuid.UUID,
	name, description string,
	noDescription bool,
	binary []byte,
) error {
	if len(binary) == 0 {
		return s.update(ctx, token, id, name, description, noDescription, nil)
	}

	_, msg, err := s.Get(ctx, token, id)
	if err != nil {
		return fmt.Errorf("SecretsService - EditBinary - uc.Get: %w", err)
	}

	data, ok := msg.(*p.Binary)
	if !ok {
		return fmt.Errorf("SecretsService - EditBinary - msg.(*goph.Binary): %w", ErrKindMismatch)
	}

	data.Binary = binary

	return s.update(ctx, token, id, name, description, noDescription, data)
}

// EditCard changes parameters of stored bank card.
func (s *SecretsService) EditCard(
	ctx context.Context,
	token string,
	id uuid.UUID,
	name, description string,
	noDescription bool,
	number, expiration, holder string,
	cvv int32,
) error {
	if number == "" && expiration == "" && holder == "" && cvv == 0 {
		return s.update(ctx, token, id, name, description, noDescription, nil)
	}

	_, msg, err := s.Get(ctx, token, id)
	if err != nil {
		return fmt.Errorf("SecretsService - EditCard - uc.Get: %w", err)
	}

	data, ok := msg.(*p.Card)
	if !ok {
		return fmt.Errorf("SecretsService - EditCard - msg.(*goph.Card): %w", ErrKindMismatch)
	}

	if number != "" {
		data.Number = number
	}

	if expiration != "" {
		data.Expiration = expiration
	}

	if holder != "" {
		data.Holder = holder
	}

	if cvv != 0 {
		data.Cvv = cvv
	}

	return s.update(ctx, token, id, name, description, noDescription, data)
}

// EditCreds changes parameters of stored credentials.
func (s *SecretsService) EditCreds(
	ctx context.Context,
	token string,
	id uuid.UUID,
	name, description string,
	noDescription bool,
	login, password string,
) error {
	if login == "" && password == "" {
		return s.update(ctx, token, id, name, description, noDescription, nil)
	}

	_, msg, err := s.Get(ctx, token, id)
	if err != nil {
		return fmt.Errorf("SecretsService - EditCreds - uc.Get: %w", err)
	}

	data, ok := msg.(*p.Credentials)
	if !ok {
		return fmt.Errorf("SecretsService - EditCreds - msg.(*goph.Credentials): %w", ErrKindMismatch)
	}

	if login != "" {
		data.Login = login
	}

	if password != "" {
		data.Password = password
	}

	return s.update(ctx, token, id, name, description, noDescription, data)
}

// EditText changes parameters of stored text secret.
func (s *SecretsService) EditText(
	ctx context.Context,
	token string,
	id uuid.UUID,
	name, description string,
	noDescription bool,
	text string,
) error {
	if text == "" {
		return s.update(ctx, token, id, name, description, noDescription, nil)
	}

	_, msg, err := s.Get(ctx, token, id)
	if err != nil {
		return fmt.Errorf("SecretsService - EditText - uc.Get: %w", err)
	}

	data, ok := msg.(*p.Text)
	if !ok {
		return fmt.Errorf("SecretsService - EditText - msg.(*goph.Text): %w", ErrKindMismatch)
	}

	data.Text = text

	return s.update(ctx, token, id, name, description, noDescription, data)
}

// Get retrieves full user's secret.
// All sensitive parts are decrypted.
func (s *SecretsService) Get(
	ctx context.Context,
	token string,
	id uuid.UUID,
) (*p.Secret, proto.Message, error) {
	secret, data, err := s.secretsRepo.Get(ctx, token, id)
	if err != nil {
		return nil, nil, fmt.Errorf("SecretsService - Get - uc.secretsRepo.Get: %w", err)
	}

	secret.Metadata, err = s.key.Decrypt(secret.GetMetadata())
	if err != nil {
		return nil, nil, fmt.Errorf("SecretsService - Get - uc.key.Decrypt(metadata): %w", err)
	}

	decryptedData, err := s.key.Decrypt(data)
	if err != nil {
		return nil, nil, fmt.Errorf("SecretsService - Get - uc.key.Decrypt(data): %w", err)
	}

	var msg proto.Message

	switch secret.GetKind() {
	case p.DataKind_BINARY:
		msg = &p.Binary{}

	case p.DataKind_CARD:
		msg = &p.Card{}

	case p.DataKind_CREDENTIALS:
		msg = &p.Credentials{}

	case p.DataKind_TEXT:
		msg = &p.Text{}
	}

	if err := proto.Unmarshal(decryptedData, msg); err != nil {
		return nil, nil, fmt.Errorf("SecretsService - Get - proto.Unmarshal: %w", err)
	}

	return secret, msg, nil
}

// Delete removes user's secret.
func (s *SecretsService) Delete(
	ctx context.Context,
	token string,
	id uuid.UUID,
) error {
	if err := s.secretsRepo.Delete(ctx, token, id); err != nil {
		return fmt.Errorf("SecretsService - Delete - uc.secretsRepo.Delete: %w", err)
	}

	return nil
}

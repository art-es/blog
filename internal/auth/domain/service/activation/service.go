//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package activation

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/art-es/blog/internal/auth/domain"
	"github.com/art-es/blog/internal/auth/dto"
)

type databus interface {
	ProduceActivationEmail(ctx context.Context, msg *dto.UserActivationEmailMessage) error
}

type Service struct {
	logger  *zap.Logger
	databus databus
}

func New(logger *zap.Logger, databus databus) *Service {
	return &Service{
		logger:  logger,
		databus: databus,
	}
}

func (s *Service) Activate(ctx context.Context, code string, tx domain.TxCommitter) error {
	uid, err := tx.ActivationCode().GetUserID(ctx, code)
	if err != nil {
		return fmt.Errorf("auth ID getting from repository error: %w", err)
	}
	if uid == 0 {
		return dto.ErrActivationCodeNotFound
	}

	ok, err := tx.User().Activate(ctx, uid)
	if err != nil {
		return fmt.Errorf("auth getting from repository error: %w", err)
	}
	if !ok {
		return dto.ErrUserNotFound
	}

	if err = tx.ActivationCode().RemoveCodes(ctx, uid); err != nil {
		return fmt.Errorf("codes removing from repository error: %w", err)
	}
	return nil
}

func (s *Service) SendCode(ctx context.Context, user *domain.User, tx domain.TxCommitter) error {
	// TODO: is need to add checking existence of code in repo?
	code := uuid.NewString()

	if err := tx.ActivationCode().Add(ctx, code, user.ID); err != nil {
		return fmt.Errorf("activation code adding to repository error: %w", err)
	}

	msg := &dto.UserActivationEmailMessage{
		Email: user.Email,
		Code:  code,
	}

	if err := s.databus.ProduceActivationEmail(ctx, msg); err != nil {
		s.logger.Error("produce message to databus error",
			zap.Error(err),
			zap.String("location", "auth/service/activation"))
	}

	return nil
}

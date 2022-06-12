package service

import (
	"context"
	"pos-microservices/cashier/auth"
	"pos-microservices/cashier/model"
	"pos-microservices/cashier/repository"
	"time"

	"github.com/golang-jwt/jwt"
)

type Service interface {
	CashierService
	AuthService
}

type service struct {
	cashierRepository repository.CashierRepository
	authService       auth.Authenticator
}

func NewService(r repository.CashierRepository, a auth.Authenticator) Service {
	return &service{
		cashierRepository: r,
		authService:       a,
	}
}

type CashierService interface {
	Signup(ctx context.Context, cashier *model.Cashier) (string, error)
	GetByID(ctx context.Context, id string) (*model.Cashier, error)
	GetByEmail(ctx context.Context, email string) (*model.Cashier, error)
	GetAll(ctx context.Context) ([]*model.Cashier, error)
	Update(ctx context.Context, cashier *model.Cashier) error
	Delete(ctx context.Context, id string) error
}

type AuthService interface {
	VerifyToken(tokenString string) (*jwt.Token, error)
	IssueToken(user string, duration time.Duration) (string, error)
}

func (s *service) VerifyToken(tokenString string) (*jwt.Token, error) {
	return s.authService.VerifyToken(tokenString)
}

func (s *service) IssueToken(user string, duration time.Duration) (string, error) {
	return s.authService.IssueToken(user, duration)
}

func (s *service) Signup(ctx context.Context, cashier *model.Cashier) (string, error) {
	cashier.Created = time.Now().Unix()

	return s.cashierRepository.Create(ctx, cashier)
}
func (s *service) GetByID(ctx context.Context, id string) (*model.Cashier, error) {
	return s.cashierRepository.GetByID(ctx, id)
}
func (s *service) GetByEmail(ctx context.Context, email string) (*model.Cashier, error) {
	panic("unimplemented")
}
func (s *service) GetAll(ctx context.Context) ([]*model.Cashier, error) {
	panic("unimplemented")
}
func (s *service) Update(ctx context.Context, cashier *model.Cashier) error {
	panic("unimplemented")
}
func (s *service) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

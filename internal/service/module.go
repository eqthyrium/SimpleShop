package service

import "SimpleShop/internal/domain"

type HttpModule interface {
	SignUp(nickname, email, password string) error
	LogIn(email, password string) (string, error)
	Homepage(userId int, searchingValue string) ([]domain.Product, map[string]bool, error)
	Purchase(userId, productId int) error
	Like(userId, productId int) error
	History(userId int) ([]domain.Product, []domain.Product, error)
	Recommendation(userId int) ([]domain.Product, []domain.Product, error)
}

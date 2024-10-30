package service

import "SimpleShop/internal/domain"

type HttpModule interface {
	SignUp(nickname, email, password string) error
	LogIn(email, password string) (string, error)
	Homepage(userId int) ([]domain.Product, error)
	Purchase(userId, productId int) error
	Like(userId, productId int) error
	//Login(email, password string) (err error)
	// Here we write what kind of services can be used in the http handler
}

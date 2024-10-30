package repository

import (
	"SimpleShop/internal/domain"
)

type DbModule interface {
	userRepository
	productRepository
	//PostRepository
	//CommentRepository
	//CategoryRepository
	//PostCategoryRepository
	//ReactionRepository
	//NotificationRepository

}

type userRepository interface {
	CreateUser(user *domain.User) error
	GetLastUserId() (int, error)
	//UpdateUser()error
	//DeleteUser(userId int) error
	//GetUserByID(userId int) (domain.User, error)
	//CheckUserByEmail(email string) (bool, error)
	GetUserByEmail(email string) (*domain.User, error)
}

type productRepository interface {
	RetrieveProducts(role string, userId int) ([]domain.Product, error)
}

//type PostRepository interface {
//}
//
//type CommentRepository interface {
//}
//
//type CategoryRepository interface {
//}
//
//type PostCategoryRepository interface {
//}
//
//type ReactionRepository interface {
//}
//
//type NotificationRepository interface {
//}

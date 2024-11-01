package repository

import (
	"SimpleShop/internal/domain"
)

type DbModule interface {
	userRepository
	productRepository
	relationship
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
	RetrievePurchasedProduct(userId int) ([]domain.Product, error)
	RetrieveLikedProduct(userId int) ([]domain.Product, error)
	RetrieveCollaborativeProduct(userId int) ([]domain.Product, error)
	RetrieveBehaviourBasedProduct(userId int) ([]domain.Product, error)
}

type relationship interface {
	Purchase(userId, productId int) error
	Like(userId, productId int) error
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

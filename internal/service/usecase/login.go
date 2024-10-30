package usecase

import (
	"SimpleShop/internal/domain"
	"SimpleShop/internal/service/session"
	"SimpleShop/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// ToDo
// 1. Check whether the client exists
// 2. Making token to that client

func (app *Application) LogIn(memberIdentity, password string) (string, error) {

	receivedUser, err := app.ServiceDB.GetUserByEmail(memberIdentity) // The handler side must check whether its error is ErrUserNotFound error, in order to be adjusted in giving back webpage
	if err != nil {
		return "", logger.ErrorWrapper("UseCase", "LogIn", "Failure in the getting the user by sending by email", err)
	}

	isTheSame := CheckPassword(receivedUser.Password, password)
	if !isTheSame {
		return "", domain.ErrUserNotFound

	}

	tokenSignature, err := session.CreateToken(receivedUser.UserId, receivedUser.Role)

	if err != nil {
		return "", logger.ErrorWrapper("UseCase", "LogIn", "Failure in the creating a token", err)
	}

	return tokenSignature, nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

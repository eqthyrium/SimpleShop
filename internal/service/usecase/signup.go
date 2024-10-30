package usecase

import (
	"SimpleShop/internal/domain"
	"SimpleShop/pkg/logger"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

// ToDo for SignUp:
// 1. Checking whether the entered input data is correct
// 2. Whether input email exist in the db
// 3. Hash the password
// 4. Insert into db (email, nickname, password(hashed), role(by default user))

func (app *Application) SignUp(nickname, email, password string) error {

	//checking whether the input data is correct
	correctnessData := app.isItCorrect(nickname, email, password)
	if !correctnessData {
		return logger.ErrorWrapper("UseCase", "SignUp", "There is an invalid entered credentials of the client to be signed up", domain.ErrInvalidCredential)
	}
	nickname, email = makeItLower(nickname, email)

	// check is there exist such email
	_, err := app.ServiceDB.GetUserByEmail(email)

	if !errors.Is(err, domain.ErrUserNotFound) {
		return logger.ErrorWrapper("UseCase", "SignUp", "The client entered such credential which is already in the data base", domain.ErrInvalidCredential)
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "SignUp", "Failed to hash password", err)
	}
	userId, err := app.ServiceDB.GetLastUserId()
	if err != nil {
		return logger.ErrorWrapper("UseCase", "SignUp", "Failed to get user id", err)
	}

	user := &domain.User{
		UserId:         userId + 1,
		Nickname:       nickname,
		MemberIdentity: email,
		Password:       hashedPassword,
		Role:           "User",
	}
	err = app.ServiceDB.CreateUser(user)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "SignUp", "Failed to create user", err)
	}

	return nil
}

func (app *Application) isItCorrect(nickname, email, password string) bool {

	answerNickname := nicknameCheck(nickname)
	answerEmail := emailCheck(email)
	answerPassword := passwordCheck(password)

	if !(answerEmail && answerPassword && answerNickname) {
		return false
	}

	return true

}

func nicknameCheck(nickname string) bool {
	nicknameRegex := `^[a-zA-Z0-9]([a-zA-Z0-9._-]{1,18}[a-zA-Z0-9])?$`
	re := regexp.MustCompile(nicknameRegex)
	return re.MatchString(nickname)
}
func emailCheck(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func passwordCheck(password string) bool {

	if len(password) < 8 || len(password) > 32 {
		return false
	}

	for i := 0; i < len(password); i++ {
		if (password[i] >= 'A' && password[i] <= 'Z') || (password[i] >= 'a' && password[i] <= 'z') {
			return true
		}
	}

	return false
}

func makeItLower(nickname, email string) (string, string) {
	return strings.ToLower(nickname), strings.ToLower(email)
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"SimpleShop/internal/domain"

	"github.com/google/uuid"
)

// ToDO
/*
1) Creation of the Token (done)
2) Verifying the signature of requested Token (done)
3) Extraction data from the Token (done)
4)
*/
type Token struct {
	UserId     int       `json:"userId"`
	UUID       string    `json:"uuid"`
	Role       string    `json:"role"`
	ExpireTime time.Time `json:"expireTimeActive"`
}

var mySecretKey string = "AddDeleteKey1618"

var MapUUID map[int]string = make(map[int]string)

func CreateToken(userId int, role string) (string, error) {
	token := Token{
		UserId:     userId,
		UUID:       uuid.New().String(),
		Role:       role,
		ExpireTime: time.Now().Add(60 * time.Minute),
	}

	MapUUID[token.UserId] = token.UUID

	tokenJson, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("Token-CreateSignedToken, token marshalling failed: %v", err)
	}

	signature := createSignature(tokenJson, mySecretKey) // hashsum

	signatureToken := base64.URLEncoding.EncodeToString(tokenJson) + "." + signature

	return signatureToken, nil
}

func VerifyToken(token string) error {

	passedToken := strings.Split(token, ".")

	if len(passedToken) != 2 {
		return fmt.Errorf("Token-ValidateToken: %w", domain.ErrInvalidToken)

	}
	playLoad := passedToken[0]
	decodedPlayLoad, err := base64.URLEncoding.DecodeString(playLoad)
	if err != nil {
		return fmt.Errorf("Token-ValidateToken: %w", err)

	}
	secondSignature := createSignature(decodedPlayLoad, mySecretKey)
	if hmac.Equal([]byte(secondSignature), []byte(passedToken[1])) {
		return nil

	}
	return fmt.Errorf("Token-ValidateToken: %w", domain.ErrInvalidToken)
}

func createSignature(tokenJson []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(tokenJson)
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func ExtractDataFromToken(token string) (*Token, error) {
	parts := strings.Split(token, ".")

	playLoad, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("Token-ExtractDataFromToken: %w", err)
	}

	var tokenObject Token = Token{}
	err = json.Unmarshal(playLoad, &tokenObject)
	if err != nil {
		return nil, fmt.Errorf("Token-ExtractDataFromToken: %w", err)
	}

	return &tokenObject, nil
}

// Extension of the token, when the particular token passed the thresholdtime, it should be added about 30 min to its lifespan
func ExtendTokenExistence(passedToken *Token) (string, error) {
	passedToken.ExpireTime = passedToken.ExpireTime.Add(45 * time.Minute)

	tokenJson, err := json.Marshal(passedToken)
	if err != nil {
		return "", fmt.Errorf("Token-CreateSignedToken, token marshalling failed: %v", err)
	}

	signature := createSignature(tokenJson, mySecretKey) // hashsum

	signatureToken := base64.URLEncoding.EncodeToString(tokenJson) + "." + signature

	return signatureToken, nil
}

func CheckTokenTime(token *Token) string {
	if token.ExpireTime.After(time.Now()) {
		// Token is still valid, check if it's nearing expiration
		if token.ExpireTime.Sub(time.Now()).Minutes() <= 15.00 {
			return "Extend-Token"
		}
		return "Valid-Token"
	}
	return "Expired-Token"

}

package auth

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"stockx-backend/db"
	"stockx-backend/db/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Credentials struct {
	Username string `json:"username"` //can be also email
	Password string `json:"password"`
}

const (
	VERKEY = "jdnfksdmfksd"
)

type Token struct {
	Token string `json:"token"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func RegisterUser(register Register) (bool, error) {
	var validUsername = regexp.MustCompile("^[a-zA-Z0-9]*[-]?[a-zA-Z0-9]*$")

	if !validUsername.MatchString(register.Username) {
		return false, errors.New("invalid username. Should contain 5-20 characters")
	}

	hashedPass, err := hashPassword(register.Password)
	if err != nil {
		return false, err
	}

	newUser := models.User{
		Username: register.Username,
		Email:    register.Email,
		Password: hashedPass,
	}

	err = db.PutItemInTable(newUser, "User")
	if err != nil {
		return false, err
	}

	// email.SendConfirmRegistrationEmail(register.Email, conf.Conf.Email)

	return true, nil
}

func LogIn(login Credentials) (Token, error) {
	var username string

	user, err := db.GetUserFromTable(login.Username)

	if err != nil {
		return Token{}, err
	}

	if !comparePasswords(user.Password, []byte(login.Password)) {
		return Token{}, errors.New("could not login - incorrect password")
	}

	token, err := CreateToken(username)
	if err != nil {
		return Token{}, errors.New("could not generate token")
	}

	returnToken := Token{
		Token: token,
	}

	return returnToken, nil
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	if err := bcrypt.CompareHashAndPassword(byteHash, plainPwd); err != nil {
		return false
	}

	return true
}

func CreateToken(username string) (string, error) {
	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["username"] = username

	atClaims["exp"] = time.Now().Add(time.Minute * 300).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(VERKEY))

	if err != nil {
		return "", err
	}

	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}

	if err := token.Claims.Valid(); err != nil {
		return err
	}

	return nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(VERKEY), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("x-auth-token")
	return bearToken
}

func GetUsernameFromToken(r *http.Request) (string, error) {
	tokenString := ExtractToken(r)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(VERKEY), nil
	})

	if err != nil {
		return "", err
	}

	return claims["username"].(string), nil
}

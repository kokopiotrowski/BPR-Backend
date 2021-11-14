package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
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

func RegisterUser(db *sql.DB, register Register) (bool, error) {
	var validUsername = regexp.MustCompile("^[a-zA-Z0-9]*[-]?[a-zA-Z0-9]*$")

	if !validUsername.MatchString(register.Username) {
		return false, errors.New("invalid username. Should contain 5-20 characters")
	}

	hashedPass, err := hashPassword(register.Password)
	if err != nil {
		return false, err
	}

	query := "INSERT INTO Users(Username, Pass) VALUES (?, ?)"

	tx, err := db.Begin()

	if err != nil {
		return false, err
	}

	_, err = tx.Exec(query, register.Username, hashedPass)

	if err != nil {
		return false, err
	}

	if err = tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

func LogIn(db *sql.DB, login Credentials) (Token, error) {
	var userID uint64

	var username, hashedPassDB string

	err := db.QueryRow("SELECT UserId, Username, Pass FROM Users WHERE Username=?", login.Username).Scan(&userID, &username, &hashedPassDB)

	if err != nil {
		return Token{}, errors.New("could not login - didn't find user with this username")
	}

	if !comparePasswords(hashedPassDB, []byte(login.Password)) {
		return Token{}, errors.New("could not login - incorrect password")
	}

	token, err := CreateToken(userID, username)
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

func CreateToken(userID uint64, username string) (string, error) {
	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID

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

func GetUserIDFromToken(r *http.Request) (int64, error) {
	tokenString := ExtractToken(r)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(VERKEY), nil
	})

	if err != nil {
		return 0, err
	}

	return int64(claims["user_id"].(float64)), nil
}

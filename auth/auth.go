package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"stockx-backend/conf"
	"stockx-backend/db"
	"stockx-backend/db/models"
	"stockx-backend/email"
	"stockx-backend/reserr"
	"stockx-backend/util"
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
	Email    string `json:"email"` //can be also email
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"authorization"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func RegisterUser(register Register) (bool, error) {
	var validUsername = regexp.MustCompile("^[a-zA-Z0-9]*[-]?[a-zA-Z0-9]*$")

	if !validUsername.MatchString(register.Username) {
		return false, reserr.BadRequest("info", errors.New("username constains not accepted symbols"), "Username constains not accepted symbols")
	}

	if _, err := mail.ParseAddress(register.Email); err != nil {
		return false, reserr.BadRequest("info", err, "Invalid email!")
	}

	hashedPass, err := HashPassword(register.Password)
	if err != nil {
		return false, err
	}

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return false, err
	}

	dt := time.Now().In(loc)

	newUser := models.User{
		Email:       register.Email,
		Username:    register.Username,
		Password:    hashedPass,
		DateCreated: dt.Format("02-01-2006 15:04:05"),
		IsAdmin:     false,
		IsConfirmed: false,
	}

	err = db.PutUserInTheTable(newUser)
	if err != nil {
		return false, reserr.Internal("error", err, "failed to save newly registered user in database")
	}

	err = db.PutTradesInTheTable(newUser.Email, models.Trades{
		Email:         newUser.Email,
		Credits:       models.DefaultAmountOfCreditsForNewUsers,
		BoughtStocks:  []models.BoughtStock{},
		SoldStocks:    []models.SoldStock{},
		ShortStocks:   []models.ShortStock{},
		BoughtToCover: []models.BoughtToCover{},
		HoldLong:      []models.HoldLong{},
		HoldShort:     []models.HoldShort{},
	})
	if err != nil {
		return false, reserr.Internal("error", err, "failed to add new trades for newly registered user in database")
	}

	err = db.PutStatisticsInTheTable(newUser.Email, models.Statistics{
		Email: newUser.Email,
	})
	if err != nil {
		return false, reserr.Internal("error", err, "failed to begin tracking statistics for newly registered user in database")
	}

	registeredUsers, err := db.GetListOfRegisteredUsers()
	if err != nil {
		return false, reserr.Internal("error", err, "Failed to add user to list of users")
	}

	registeredUsers.Users = append(registeredUsers.Users, newUser.Email)

	err = db.PutListOfRegisteredUsersInTheTable(registeredUsers)
	if err != nil {
		return false, reserr.Internal("error", err, "Failed to add user to list of users")
	}

	confirmToken, err := generateConfirmationToken(newUser.Email)
	if err != nil {
		return false, reserr.Internal("error", err, "Failed to generate confirmation account link for user")
	}

	go func() {
		err = email.SendConfirmRegistrationEmail(newUser.Email, newUser.Username, confirmToken)
	}()

	return true, err
}

func LogIn(login Credentials) (Token, error) {
	user, err := db.GetUserFromTable(login.Email)
	if err != nil {
		return Token{}, err
	}

	if !user.IsConfirmed {
		return Token{}, reserr.Forbidden("error", errors.New("failed to login - user account not confirmed"), "User registration has not been confirmed yet")
	}

	if !ComparePasswords(user.Password, []byte(login.Password)) {
		return Token{}, reserr.Forbidden("error", errors.New("failed to login - wrong password"), "Incorrect password")
	}

	token, err := CreateToken(user.Email, user.IsAdmin)
	if err != nil {
		return Token{}, reserr.Internal("error", err, "Failed to Login - try again later")
	}

	returnToken := Token{
		Token: token,
	}

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return Token{}, reserr.Internal("error", err, "Failed to Login - try again later")
	}

	dt := time.Now().In(loc)

	user.DateLastAccessed = dt.Format("01-02-2006 15:04:05")

	err = db.PutUserInTheTable(user)
	if err != nil {
		return Token{}, reserr.Internal("error", err, "Failed to Login - try again later")
	}

	return returnToken, nil
}

func GetUserEmailFromRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	err := TokenValid(r)
	if err != nil {
		return "", err
	}

	u, err := GetEmailFromToken(r)
	if err != nil {
		return "", err
	}

	return u, nil
}

func CheckIfAuthorized(w http.ResponseWriter, r *http.Request, email *string) bool { //email, error
	err := TokenValid(r)
	if err != nil {
		util.RespondWithJSON(w, r, http.StatusUnauthorized, "Could not retrieve user id from token", nil)
		return false
	}

	u, err := GetEmailFromToken(r)
	if err != nil {
		util.RespondWithJSON(w, r, http.StatusUnauthorized, "Could not retrieve user id from token", nil)
		return false
	}

	*email = u

	return true
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	if err := bcrypt.CompareHashAndPassword(byteHash, plainPwd); err != nil {
		return false
	}

	return true
}

func CreateToken(email string, isAdmin bool) (string, error) {
	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["email"] = email
	atClaims["admin"] = isAdmin

	atClaims["exp"] = time.Now().Add(time.Minute * 300).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(conf.Conf.Server.EncodingVerKey))

	if err != nil {
		return "", err
	}

	return token, nil
}

func generateConfirmationToken(email string) (string, error) {
	var err error

	atClaims := jwt.MapClaims{}
	atClaims["email"] = email

	atClaims["exp"] = time.Now().Add(time.Minute * 300).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(conf.Conf.Server.EncodingVerKey))

	if err != nil {
		return "", err
	}

	return token, nil
}

func ExtractEmailFromConfirmationToken(token string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Conf.Server.EncodingVerKey), nil
	})

	if err != nil {
		return "", err
	}

	return claims["email"].(string), nil
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
		return []byte(conf.Conf.Server.EncodingVerKey), nil
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

func GetEmailFromToken(r *http.Request) (string, error) {
	tokenString := ExtractToken(r)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Conf.Server.EncodingVerKey), nil
	})

	if err != nil {
		return "", err
	}

	return claims["email"].(string), nil
}

func CheckIfAdmin(r *http.Request) (bool, error) {
	tokenString := ExtractToken(r)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Conf.Server.EncodingVerKey), nil
	})

	if err != nil {
		return false, err
	}

	return claims["admin"].(bool), nil
}

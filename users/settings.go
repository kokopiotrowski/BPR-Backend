package users

import (
	"errors"
	"stockx-backend/auth"
	"stockx-backend/db"
	"stockx-backend/reserr"
)

type UserSettings struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Joined   string `json:"joined"` //date that user created account
}

type UpdatedPassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

var confirmationPage = `<center><h2>Account confirmed, go to our page and login <a href="http://bpr-frontend.s3-website.eu-north-1.amazonaws.com/login">StockX</a></h2></center>`

func DeleteUser(email string) error {
	err := db.DeleteUser(email)
	if err != nil {
		return err
	}

	err = db.DeleteStatistics(email)
	if err != nil {
		return err
	}

	err = db.DeleteTrades(email)
	if err != nil {
		return err
	}

	listOfUsers, err := db.GetListOfRegisteredUsers()
	if err != nil {
		return err
	}

	for i, registeredUserEmail := range listOfUsers.Users {
		if registeredUserEmail == email {
			listOfUsers.Users = append(listOfUsers.Users[:i], listOfUsers.Users[i+1:]...)
			break
		}
	}

	err = db.PutListOfRegisteredUsersInTheTable(listOfUsers)

	return err
}

func GetUser(email string) (UserSettings, error) {
	user, err := db.GetUserFromTable(email)
	if err != nil {
		return UserSettings{}, err
	}

	return UserSettings{Email: user.Email, Username: user.Username, Joined: user.DateCreated}, nil
}

func ChangePassword(email string, model UpdatedPassword) error {
	user, err := db.GetUserFromTable(email)
	if err != nil {
		return err
	}

	if !auth.ComparePasswords(user.Password, []byte(model.OldPassword)) {
		return errors.New("invalid password")
	}

	newHashedPass, err := auth.HashPassword(model.NewPassword)
	if err != nil {
		return err
	}

	err = db.UpdateUsersPassword(email, newHashedPass)

	return err
}

func ConfirmUserAccount(token string) (string, error) {
	email, err := auth.ExtractEmailFromConfirmationToken(token)
	if err != nil {
		return "", reserr.BadRequest("error", err, "Failed to confirm user")
	}

	user, err := db.GetUserFromTable(email)
	if err != nil {
		return "", reserr.BadRequest("error", err, "Failed to confirm user")
	}

	user.IsConfirmed = true

	err = db.PutUserInTheTable(user)
	if err != nil {
		return "", reserr.BadRequest("error", err, "Failed to confirm user")
	}

	return confirmationPage, nil
}

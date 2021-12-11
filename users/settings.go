package users

import (
	"errors"
	"stockx-backend/auth"
	"stockx-backend/db"
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

func DeleteUser(email string) error {
	err := db.DeleteUser(email)
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

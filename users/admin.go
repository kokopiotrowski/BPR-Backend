package users

import (
	"stockx-backend/db"
)

type UserInformationForAdmin struct {
	Email            string `json:"email"`
	Username         string `json:"username"`
	IsAdmin          bool   `json:"isAdmin"`
	DateCreated      string `json:"dateCreated"`
	DateUpdated      string `json:"dateUpdated"`
	DateLastAccessed string `json:"dateLastAccessed"`
	IsConfirmed      bool   `json:"isConfirmed"`
}

func GetListOfUserData() ([]UserInformationForAdmin, error) {
	listOfRegisteredUsers, err := db.GetListOfRegisteredUsers()
	if err != nil {
		return []UserInformationForAdmin{}, err
	}

	users := []UserInformationForAdmin{}

	for _, email := range listOfRegisteredUsers.Users {
		u, err := db.GetUserFromTable(email)
		if err != nil {
			return []UserInformationForAdmin{}, err
		}

		returnUser := UserInformationForAdmin{
			Email:            u.Email,
			Username:         u.Username,
			IsAdmin:          u.IsAdmin,
			DateCreated:      u.DateCreated,
			DateUpdated:      u.DateUpdated,
			DateLastAccessed: u.DateLastAccessed,
			IsConfirmed:      u.IsConfirmed,
		}

		users = append(users, returnUser)
	}

	return users, nil
}

package models

type User struct {
	Email            string `json:"email"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	IsAdmin          bool   `json:"isAdmin"`
	DateCreated      string `json:"dateCreated"`
	DateUpdated      string `json:"dateUpdated"`
	DateLastAccessed string `json:"dateLastAccessed"`
	IsConfirmed      bool   `json:"isConfirmed"`
}

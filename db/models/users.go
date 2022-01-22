package models

var (
	DefaultAmountOfCreditsForNewUsers = float32(100000)
)

type ListOfRegisteredUsers struct {
	Region string   `json:"region"` // future possibility to separate registered users by region
	Users  []string `json:"users"`
}

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

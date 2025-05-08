package models

type UserTokenInfo struct {
	ID   string
	Role string
}

type User struct {
	Email        string `json:"email"`
	StorageUsed  int64  `json:"storageUsed"`
	StorageLimit int64  `json:"storageLimit"`
	IsAdmin      bool   `json:"isAdmin"`
}

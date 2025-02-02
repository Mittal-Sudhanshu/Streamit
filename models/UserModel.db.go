// package models

// import "go.mongodb.org/mongo-driver/mongo"

//	type User struct {
//		Id ObjectId `bson:"_id" json:"id"`
//	}
package models

type GoogleResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Verified bool   `json:"verified_email"`
	Picture  string `json:"picture"`
}

type User struct {
	ID      string `bson:"_id" json:"id"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
	IsAdmin bool   `json:"isAdmin"`
}


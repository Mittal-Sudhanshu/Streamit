// package models

// import "go.mongodb.org/mongo-driver/mongo"

//	type User struct {
//		Id ObjectId `bson:"_id" json:"id"`
//	}
package models

import "time"

type User struct {
	Id          string    `bson:"_id" json:"id"`
	Email       string    `bson:"email" json:"email"`
	Picture_url string    `bson:"picture_url" json:"picture_url"`
	Name        string    `bson:"name" json:"name"`
	Admin       bool      `bson:"admin" json:"admin"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}


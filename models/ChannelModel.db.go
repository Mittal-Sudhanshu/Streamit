package models

type Channel struct {
	ID          string `json:"id" bson:"_id"`
	ChannelName string `json:"channelName"`
	UserId      string `json:"userId"`
}

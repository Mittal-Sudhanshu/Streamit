package models

import "time"

type Privacy int

const (
	Private Privacy = iota
	Unlisted
	Public
)

type Stream struct {
	StreamId      string    `json:"streamId"`
	StreamKey     string    `json:"streamKey"`
	UserId        string    `json:"userId"`
	ChannelName   string    `json:"channelName"`
	StreamTitle   string    `json:"streamTitle"`
	StreamDesc    string    `json:"streamDesc"`
	StartTime     time.Time `json:"startTime"`
	ScheduledTime time.Time `json:"scheduledTime"`
	IsLive        bool      `json:"isLive"`
	IsScheduled   bool      `json:"isScheduled"`
	Privacy       Privacy   `json:"privacy"`
}

type StreamMetaData struct {
	Stream      Stream `json:"stream"`
	Views       int    `json:"views"`
	FlvFilePath string `json:"flvFilePath"`
	LocalIp     string `json:"localIp"`
}

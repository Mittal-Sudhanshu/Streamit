package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"streamit/models"
	"streamit/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Channel struct {
	ChannelName string "json:channelName"
	ID          string "bson:_id"
	UserId      string "json:userId"
}
type StreamBody struct {
	StreamTitle string "json:streamTitle"
	StreamDesc  string "json:streamDesc"
	Privacy     string "json:privacy"
	Scheduled   bool   "json:scheduled"
	ScheduledAt string "json:scheduledAt"
}

func generateRandomString(length int) string {
	// Create a byte slice of the desired length
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err) // Handle error
	}

	// Encode the random bytes to a URL-safe base64 string and return the desired substring
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func CreateNewStreamKey(c *fiber.Ctx) error {
	var newStreamKey string = generateRandomString(5)
	var channel Channel
	var streamData models.Stream
	if err := c.BodyParser(&streamData); err != nil {
		// log.Println("Body: ", body.ChannelName)
	}
	err := utils.ChannelDb.FindOne(c.Context(), bson.M{"userId": c.Locals("userID")}).Decode(&channel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Channel not found",
			})
		}
		log.Println("MongoDB error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	stream := models.Stream{
		ChannelName: channel.ChannelName,
		StreamKey:   newStreamKey,
		UserId:      c.Locals("userID").(string),
		StreamTitle: streamData.StreamTitle,
		StreamDesc:  streamData.StreamDesc,
		IsLive:      false,
		IsScheduled: false,
		Privacy:     models.Public,
	}
	log.Print(stream)
	utils.StreamDb.InsertOne(c.Context(), bson.M{"channelName": channel.ChannelName, "streamKey": newStreamKey, "userId": c.Locals("userID"), "streamTitle": streamData.StreamTitle, "streamDesc": streamData.StreamDesc, "isLive": false, "isScheduled": false, "privacy": models.Public})
	return c.JSON(fiber.Map{"url": "rtmp://localhost:1935/" + channel.ChannelName, "streamKey": newStreamKey})
}

func CreateNewChannel(c *fiber.Ctx) error {
	var body Channel

	if err := c.BodyParser(&body); err != nil {
		// log.Println("Body: ", body.ChannelName)
	}
	findUserWithChannel, err := utils.ChannelDb.Find(c.Context(), bson.M{"userId": c.Locals("userID")})
	var ChannelsWithUser []Channel
	findUserWithChannel.All(c.Context(), &ChannelsWithUser)
	if err != nil {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "User Already owns a channel"})
	}
	if len(ChannelsWithUser) != 0 {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "User Already owns a channel"})
	}
	findChannelName, err := utils.ChannelDb.Find(c.Context(), bson.M{"channelName": body.ChannelName})
	var Channels []Channel
	findChannelName.All(c.Context(), &Channels)

	if err != nil {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "Channel Name already exists"})
	}
	if len(Channels) != 0 {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "Channel Name already exists"})
	}

	createdChannel, err := utils.ChannelDb.InsertOne(c.Context(), bson.M{"channelName": body.ChannelName, "userId": c.Locals("userID")})
	if err != nil {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "Channel Name already exists"})
	}
	return c.JSON(fiber.Map{"message": "Channel Created", "channel": createdChannel})
}

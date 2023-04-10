package controllers

import (
	"context"
	"module/configs"
	"module/models"
	"module/responses"
	"module/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var FilmCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "Films")

func CreateFilm(c *fiber.Ctx) error {
	ctx, canceled := context.WithTimeout(context.Background(), time.Second*10)
	auth, _ := JwtFromHeader(c, fiber.HeaderAuthorization)
	var filmP models.Film
	var user models.User
	defer canceled()
	claims, err := utils.EncodeAccsesToken(auth)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()}})
	}
	if bodyErr := c.BodyParser(&filmP); bodyErr != nil {
		c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": bodyErr.Error()}})
	}

	if validErr := validate.Struct(&filmP); validErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": validErr.Error()},
		})
	}

	objId, _ := primitive.ObjectIDFromHex(claims.Id)

	usrErr := UserCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if usrErr != nil {
		c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": usrErr.Error()},
		})
	}
	if !user.IsAdmin {
		c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": "not accept permission"},
		})
	}
	newFilm := models.Film{
		Id:          primitive.NewObjectID(),
		Title:       filmP.Title,
		CountTicket: filmP.CountTicket,
		TimeStart:   filmP.TimeStart,
		TimeEnd:     filmP.TimeEnd}

	result, mongoErr := TicketCollection.InsertOne(ctx, newFilm)

	if mongoErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": mongoErr.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": result}})
}

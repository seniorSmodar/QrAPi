package controllers

import (
	"context"
	"module/configs"
	"module/models"
	"module/responses"
	"module/utils"
	"net/http"
	"time"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// "context"
// "net/http"
// "strconv"
// "time"
// "workspace/configs"
// "workspace/models"
// "workspace/responses"
// "workspace/utils"

// "github.com/gofiber/fiber/v2"
// "go.mongodb.org/mongo-driver/bson"
// "go.mongodb.org/mongo-driver/bson/primitive"
// "go.mongodb.org/mongo-driver/mongo"
var TicketCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "Tickets")

func createTicket(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	auth, _ := JwtFromHeader(c, fiber.HeaderAuthorization)
	var user models.User
	var ticketP models.Ticket
	defer cancel()

	claims, err := utils.EncodeAccsesToken(auth)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.Response{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}
	if bodyErr := c.BodyParser(&ticketP); bodyErr != nil {
		c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": bodyErr.Error()}})
	}

	if validErr := validate.Struct(&ticketP); validErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": validErr.Error()},
		})
	}

	objId, _ := primitive.ObjectIDFromHex(claims.Id)

	userErr := UserCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if userErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": userErr.Error()}})
	}

	if !user.IsAdmin {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": "not accept permission"},
		})
	}

	var film models.Film

	filmErr := FilmCollection.FindOne(ctx, bson.M{"_id": ticketP.FilmId}).Decode(film)
	if filmErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": filmErr.Error()}})
	}

	newTicket := models.Ticket{
		Id:     primitive.NewObjectID(),
		UserId: user.Id,
		FilmId: film.Id,
		Row:    ticketP.Row,
		Place:  ticketP.Place,
		Price:  ticketP.Price,
		Date:   time.Now(),
	}
	result, err := TicketCollection.InsertOne(ctx, newTicket)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{
		Status:  http.StatusOK,
		Message: "error",
		Data:    &fiber.Map{"data": result}})

}

func DeleteTicket(c *fiber.Ctx)error{
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*10)
	auth,_:=JwtFromHeader(c,fiber.HeaderAuthorization)
	ticketId:=c.Params("ticketId")
	var user models.User
	defer cancel()
	claims,err:=utils.EncodeAccsesToken(auth)
	if err !=nil{
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":err.Error()}})
	}
	objId,_:=primitive.ObjectIDFromHex(claims.Id)
	userErr:=UserCollection.FindOne(ctx,bson.M{"_id":objId}).Decode(&user); if userErr!=nil{
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":userErr.Error()}})
	}

	if(!user.IsAdmin){
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"error":"not accept permission"}})
	}

	obj2Id,_:=primitive.ObjectIDFromHex(ticketId)
	result,ticketErr:=TicketCollection.DeleteOne(ctx,bson.M{"_id":obj2Id})
	if ticketErr!=nil{
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":ticketErr.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data: &fiber.Map{"data":result}})


}


func BuyTicket(c *fiber.Ctx)error{
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*10)
	auth,_:=JwtFromHeader(c,fiber.HeaderAuthorization)
	filmId:=c.Params("filmId")
	countBuy:=c.Params("countBuy")
	defer cancel()
	claims,err:=utils.EncodeAccsesToken(auth)
	if err!=nil{
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":err.Error()}})
	}
	objId,_:=primitive.ObjectIDFromHex(claims.Id)
	var user models.User
	userErr:=UserCollection.FindOne(ctx,bson.M{"_id":objId}).Decode(&user);if userErr!=nil{
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":userErr.Error()}})
	}
	obj2Id,_:=primitive.ObjectIDFromHex(filmId)
	var film models.Film
	filmErr:=FilmCollection.FindOne(ctx,bson.M{"_id":obj2Id}).Decode(&film); if filmErr!=nil{
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":filmErr.Error()}})
	}
	 countBuyInt, err := strconv.Atoi(countBuy); if  err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":"count Ticket must be int"},
		})
	}
	
	if film.CountTicket-int64(countBuyInt)<0{
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":"ticket of less than zero"},
		})
	}
	film.CountTicket-=int64(countBuyInt)
	
	update := bson.M{"countTicket": film.CountTicket}

    result, mongoErr := TicketCollection.UpdateOne(ctx, bson.M{"_id": obj2Id}, bson.M{"$set": update})
	if(mongoErr!=nil){
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":mongoErr.Error()}})

	}
	return c.Status(http.StatusBadRequest).JSON(responses.Response{
		Status: fiber.StatusOK,
		Message: "success",
		Data: &fiber.Map{"data":result}})
}

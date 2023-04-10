package models
import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct{
	Id primitive.ObjectID `bson:"_id"`
	FilmId primitive.ObjectID `bson:"filmId"`
	UserId primitive.ObjectID `bson:"userId"`
	Price   float32  `json:"price" validate:"required"`
	Place  int16 `json:"place" validate:"required"`
	Row int16 `json:"row" validate:"required"`
	Date time.Time `json:"date" validate:"required"`


}

 
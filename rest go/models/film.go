package models
import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Film struct{
	Id primitive.ObjectID `bson:"_id"`
	Title string `json:"title" validate:"required"`
	TimeStart time.Time `json:"timeStart" validate:"required"`
	TimeEnd time.Time `json:"timeEnd" validate:"required"`
	CountTicket int64 `json:"countTicket"`
	

}
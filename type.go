package pasetobackendpasabar

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
	Email    string `bson:"email,omitempty" json:"email,omitempty"`
	Token    string `json:"token,omitempty" bson:"token,omitempty"`
	Private  string `json:"private,omitempty" bson:"private,omitempty"`
	Public   string `json:"public,omitempty" bson:"public,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Response struct {
	Status  bool        `json:"status" bson:"status"`
	Message string      `json:"message" bson:"message"`
	Data    interface{} `json:"data" bson:"data"`
}

type Catalog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" `
	Nomorid     int                `json:"nomorid" bson:"nomorid"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Image       string             `json:"image" bson:"image"`
	Status      bool               `json:"status" bson:"status"`
}

type About struct {
	ID          int    `json:"id" bson:"id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Image       string `json:"image" bson:"image"`
	Status      bool   `json:"status" bson:"status"`
}

type Tour struct {
	ID          int       `json:"id" bson:"id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Cari        string    `json:"cari" bson:"cari"`
	Tanggal     string    `json:"tanggal" bson:"tanggal"`
	Image       string    `json:"image" bson:"image"`
	Harga       int       `json:"harga" bson:"harga"`
	Catalog     []Catalog `json:"catalog" bson:"catalog"`
	Rating      string    `json:"rating" bson:"rating"`
	Status      bool      `json:"status" bson:"status"`
}

type HotelResto struct {
	ID          int       `json:"id" bson:"id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Cari        string    `json:"cari" bson:"cari"`
	Tanggal     string    `json:"tanggal" bson:"tanggal"`
	Image       string    `json:"image" bson:"image"`
	Harga       int       `json:"harga" bson:"harga"`
	Catalog     []Catalog `json:"catalog" bson:"catalog"`
	Rating      string    `json:"rating" bson:"rating"`
	Status      bool      `json:"status" bson:"status"`
}

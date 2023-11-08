package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email    string             `bson:"email,omitempty" json:"email,omitempty"`
	Username string             `bson:"username,omitempty" json:"username,omitempty"`
	Password string             `bson:"password,omitempty" json:"password,omitempty"`
	Role     string             `bson:"role,omitempty" json:"role,omitempty"`
}

type Credential struct {
	Status  bool   `bson:"status" json:"status"`
	Token   string `bson:"token,omitempty" json:"token,omitempty"`
	Message string `bson:"message,omitempty" json:"message,omitempty"`
}

type Catalog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Image       string             `json:"image" bson:"image"`
	IsDone      bool               `bson:"isdone,omitempty" json:"isdone,omitempty"`
}

type CatalogList struct {
	Items []Catalog `bson:"items,omitempty" json:"items,omitempty"`
}

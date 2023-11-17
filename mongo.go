package pasetobackendpasabar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetConnectionMongo(MongoString, dbname string) *mongo.Database {
	MongoInfo := atdb.DBInfo{
		DBString: os.Getenv(MongoString),
		DBName:   dbname,
	}
	conn := atdb.MongoConnect(MongoInfo)
	return conn
}

func SetConnection(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func CreateUser(mongoconn *mongo.Database, collection string, userdata User) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPassword(userdata.Password)
	if err != nil {
		return err
	}
	privateKey, publicKey := watoken.GenerateKey()
	userid := userdata.Username
	tokenstring, err := watoken.Encode(userid, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokenstring)
	// decode token to get userid
	useridstring := watoken.DecodeGetId(publicKey, tokenstring)
	if useridstring == "" {
		fmt.Println("expire token")
	}
	fmt.Println(useridstring)
	userdata.Private = privateKey
	userdata.Public = publicKey
	userdata.Password = hashedPassword

	// Insert the user data into the database
	return atdb.InsertOneDoc(mongoconn, collection, userdata)
}

// func GetNameAndPassowrd(mongoconn *mongo.Database, collection string) []User {
// 	user := atdb.GetAllDoc[[]User](mongoconn, collection)
// 	return user
// }

func CreateNewUserRole(mongoconn *mongo.Database, collection string, userdata User) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPassword(userdata.Password)
	if err != nil {
		return err
	}
	userdata.Password = hashedPassword

	// Insert the user data into the database
	return atdb.InsertOneDoc(mongoconn, collection, userdata)
}

func DeleteUser(mongoconn *mongo.Database, collection string, userdata User) interface{} {
	filter := bson.M{"username": userdata.Username}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func ReplaceOneDoc(mongoconn *mongo.Database, collection string, filter bson.M, userdata User) interface{} {
	return atdb.ReplaceOneDoc(mongoconn, collection, filter, userdata)
}

func FindUser(mongoconn *mongo.Database, collection string, userdata User) User {
	filter := bson.M{"username": userdata.Username}
	return atdb.GetOneDoc[User](mongoconn, collection, filter)
}

func IsPasswordValid(mongoconn *mongo.Database, collection string, userdata User) bool {
	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mongoconn, collection, filter)
	return CheckPasswordHash(userdata.Password, res.Password)
}

func CreateUserAndAddToken(privateKeyEnv string, mongoconn *mongo.Database, collection string, userdata User) error {
	// Hash the password before storing it
	hashedPassword, err := HashPassword(userdata.Password)
	if err != nil {
		return err
	}
	userdata.Password = hashedPassword

	// Create a token for the user
	tokenstring, err := watoken.Encode(userdata.Username, os.Getenv(privateKeyEnv))
	if err != nil {
		return err
	}

	userdata.Token = tokenstring

	// Insert the user data into the MongoDB collection
	if err := atdb.InsertOneDoc(mongoconn, collection, userdata.Username); err != nil {
		return nil // Mengembalikan kesalahan yang dikembalikan oleh atdb.InsertOneDoc
	}

	// Return nil to indicate success
	return nil
}

func InsertUserdata(MongoConn *mongo.Database, username, role, password string) (InsertedID interface{}) {
	req := new(User)
	req.Username = username
	req.Password = password
	req.Role = role
	return InsertOneDoc(MongoConn, "user", req)
}

func InsertOneDoc(db *mongo.Database, collection string, doc interface{}) (insertedID interface{}) {
	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	return insertResult.InsertedID
}

func InsertUser(db *mongo.Database, collection string, userdata User) string {
	hash, _ := HashPassword(userdata.Password)
	userdata.Password = hash
	atdb.InsertOneDoc(db, collection, userdata)
	return "Username : " + userdata.Username + "\nPassword : " + userdata.Password
}

// catalog
func CreateNewCatalog(mongoconn *mongo.Database, collection string, catalogdata Catalog) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, catalogdata)
}

// catalog function
func CreateCatalog(mongoconn *mongo.Database, collection string, catalogdata Catalog) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, catalogdata)
}

func DeleteCatalog(mongoconn *mongo.Database, collection string, catalogdata Catalog) interface{} {
	filter := bson.M{"nomorid": catalogdata.Nomorid}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedCatalog(mongoconn *mongo.Database, collection string, filter bson.M, catalogdata Catalog) interface{} {
	filter = bson.M{"nomorid": catalogdata.Nomorid}
	return atdb.ReplaceOneDoc(mongoconn, collection, filter, catalogdata)
}

func GetAllCatalog(mongoconn *mongo.Database, collection string) []Catalog {
	catalog := atdb.GetAllDoc[[]Catalog](mongoconn, collection)
	return catalog
}

func GetAllCatalogID(mongoconn *mongo.Database, collection string, catalogdata Catalog) Catalog {
	filter := bson.M{
		"nomorid":     catalogdata.Nomorid,
		"title":       catalogdata.Title,
		"description": catalogdata.Description,
		"image":       catalogdata.Image,
	}
	catalogID := atdb.GetOneDoc[Catalog](mongoconn, collection, filter)
	return catalogID
}

// about function

func CreateAbout(mongoconn *mongo.Database, collection string, aboutdata About) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, aboutdata)
}

func DeleteAbout(mongoconn *mongo.Database, collection string, aboutdata About) interface{} {
	filter := bson.M{"id": aboutdata.ID}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedAbout(mongoconn *mongo.Database, collection string, filter bson.M, aboutdata About) interface{} {
	filter = bson.M{"id": aboutdata.ID}
	return atdb.ReplaceOneDoc(mongoconn, collection, filter, aboutdata)
}

func GetAllAbout(mongoconn *mongo.Database, collection string) []About {
	about := atdb.GetAllDoc[[]About](mongoconn, collection)
	return about
}

func GetIDAbout(mongoconn *mongo.Database, collection string, aboutdata About) About {
	filter := bson.M{"id": aboutdata.ID}
	return atdb.GetOneDoc[About](mongoconn, collection, filter)
}

// <--- ini tour--->

// tour post
func GCFCreateTour(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datatour Tour
	err := json.NewDecoder(r.Body).Decode(&datatour)
	if err != nil {
		return err.Error()
	}

	if err := CreateTour(mconn, collectionname, datatour); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Create Tour", datatour))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Create Tour", datatour))
	}
}

// delete tour
func GCFDeleteTour(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datatour Tour
	err := json.NewDecoder(r.Body).Decode(&datatour)
	if err != nil {
		return err.Error()
	}

	if err := DeleteTour(mconn, collectionname, datatour); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Delete Tour", datatour))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Delete Tour", datatour))
	}
}

// update tour
func GCFUpdateTour(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datatour Tour
	err := json.NewDecoder(r.Body).Decode(&datatour)
	if err != nil {
		return err.Error()
	}

	if err := UpdatedTour(mconn, collectionname, bson.M{"id": datatour.ID}, datatour); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Update Tour", datatour))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Update Tour", datatour))
	}
}

// get all tour
func GCFGetAllTour(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datatour := GetAllTour(mconn, collectionname)
	if datatour != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Tour", datatour))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Tour", datatour))
	}
}

// get all tour by id
func GCFGetAllTourID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datatour Tour
	err := json.NewDecoder(r.Body).Decode(&datatour)
	if err != nil {
		return err.Error()
	}

	tour := GetAllTourId(mconn, collectionname, datatour)
	if tour != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Tour", tour))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Tour", tour))
	}
}

// <--- ini hotelresto--->

// hotelresto post
func GCFCreateHotelResto(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datahotelresto HotelResto
	err := json.NewDecoder(r.Body).Decode(&datahotelresto)
	if err != nil {
		return err.Error()
	}

	if err := CreateHotelResto(mconn, collectionname, datahotelresto); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Create HotelResto", datahotelresto))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Create HotelResto", datahotelresto))
	}
}

// delete hotelresto
func GCFDeleteHotelResto(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datahotelresto HotelResto
	err := json.NewDecoder(r.Body).Decode(&datahotelresto)
	if err != nil {
		return err.Error()
	}

	if err := DeleteHotelResto(mconn, collectionname, datahotelresto); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Delete HotelResto", datahotelresto))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Delete HotelResto", datahotelresto))
	}
}

// update hotelresto
func GCFUpdateHotelResto(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datahotelresto HotelResto
	err := json.NewDecoder(r.Body).Decode(&datahotelresto)
	if err != nil {
		return err.Error()
	}

	if err := UpdatedHotelResto(mconn, collectionname, bson.M{"id": datahotelresto.ID}, datahotelresto); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Update HotelResto", datahotelresto))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Update HotelResto", datahotelresto))
	}
}

// get all hotelresto
func GCFGetAllHotelResto(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datahotelresto := GetAllHotelResto(mconn, collectionname)
	if datahotelresto != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All HotelResto", datahotelresto))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All HotelResto", datahotelresto))
	}
}

// get all hotelresto by id
func GCFGetAllHotelRestoID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datahotelresto HotelResto
	err := json.NewDecoder(r.Body).Decode(&datahotelresto)
	if err != nil {
		return err.Error()
	}

	hotelresto := GetAllHotelRestoId(mconn, collectionname, datahotelresto)
	if hotelresto != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All HotelResto", hotelresto))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All HotelResto", hotelresto))
	}
}

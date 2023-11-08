package modul

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/badoux/checkmail"
	model "github.com/pasabar/pasetobackendpasabar/model"
)

func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

func InsertOneDoc(db *mongo.Database, col string, docs interface{}) (insertedID primitive.ObjectID, err error) {
	cols := db.Collection(col)
	result, err := cols.InsertOne(context.Background(), docs)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return
}

// user
func Register(db *mongo.Database, col string, userdata model.User) error {
	if userdata.Username == "" || userdata.Password == "" || userdata.Email == "" {
		return fmt.Errorf("Data tidak lengkap")
	}

	// Periksa apakah email valid
	if err := checkmail.ValidateFormat(userdata.Email); err != nil {
		return fmt.Errorf("Email tidak valid")
	}

	// Periksa apakah email dan username sudah terdaftar
	userExists, _ := GetUserFromEmail(db, col, userdata.Email)
	if userExists.Email != "" {
		return fmt.Errorf("Email sudah terdaftar")
	}
	userExists, _ = GetUserFromUsername(db, col, userdata.Username)
	if userExists.Username != "" {
		return fmt.Errorf("Username sudah terdaftar")
	}

	// Periksa apakah password memenuhi syarat
	if len(userdata.Password) < 6 {
		return fmt.Errorf("Password minimal 6 karakter")
	}
	if strings.Contains(userdata.Password, " ") {
		return fmt.Errorf("Password tidak boleh mengandung spasi")
	}

	// Periksa apakah username memenuhi syarat
	if strings.Contains(userdata.Username, " ") {
		return fmt.Errorf("Username tidak boleh mengandung spasi")
	}

	// Simpan pengguna ke basis data
	hash, _ := HashPassword(userdata.Password)
	user := bson.M{
		"_id":      primitive.NewObjectID(),
		"email":    userdata.Email,
		"username": userdata.Username,
		"password": hash,
		"role":     "user",
	}
	_, err := InsertOneDoc(db, col, user)
	if err != nil {
		return fmt.Errorf("SignUp: %v", err)
	}
	return nil
}

func LogIn(db *mongo.Database, col string, userdata model.User) (user model.User, status bool, err error) {
	if userdata.Username == "" || userdata.Password == "" {
		err = fmt.Errorf("Data tidak lengkap")
		return user, false, err
	}

	// Periksa apakah pengguna dengan username tertentu ada
	userExists, _ := GetUserFromUsername(db, col, userdata.Username)
	if userExists.Username == "" {
		err = fmt.Errorf("Username tidak ditemukan")
		return user, false, err
	}

	// Periksa apakah kata sandi benar
	if !CheckPasswordHash(userdata.Password, userExists.Password) {
		err = fmt.Errorf("Password salah")
		return user, false, err
	}
	return userExists, true, nil
}

func UpdateUser(db *mongo.Database, col string, userdata model.User) (user model.User, status bool, err error) {
	if userdata.Username == "" || userdata.Email == "" {
		err = fmt.Errorf("Data tidak boleh kosong")
		return user, false, err
	}

	// Simpan pengguna ke basis data
	existingUser, err := GetUserFromID(db, col, userdata.ID)
	if err != nil {
		return user, false, err
	}

	// Periksa apakah data yang akan diupdate sama dengan data yang sudah ada
	if userdata.Username == existingUser.Username && userdata.Email == existingUser.Email {
		err = fmt.Errorf("Data yang ingin diupdate tidak boleh sama")
		return user, false, err
	}

	checkmail.ValidateFormat(userdata.Email)
	if err != nil {
		err = fmt.Errorf("Email tidak valid")
		return user, false, err
	}

	// Periksa apakah username memenuhi syarat
	if strings.Contains(userdata.Username, " ") {
		err = fmt.Errorf("Username tidak boleh mengandung spasi")
		return user, false, err
	}

	// Simpan pengguna ke basis data
	hash, _ := HashPassword(userdata.Password)
	filter := bson.M{"_id": userdata.ID}
	update := bson.M{
		"$set": bson.M{
			"email":    userdata.Email,
			"username": userdata.Username,
			"password": hash,
			"role":     "user",
		},
	}
	cols := db.Collection(col)
	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return user, false, err
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("Data tidak berhasil diupdate")
		return user, false, err
	}
	return user, true, nil
}

func ChangePassword(db *mongo.Database, col string, userdata model.User) (user model.User, status bool, err error) {
	// Periksa apakah pengguna dengan username tertentu ada
	userExists, err := GetUserFromUsername(db, col, userdata.Username)
	if err != nil {
		return user, false, err
	}

	// Periksa apakah password memenuhi syarat
	if userdata.Password == "" {
		err = fmt.Errorf("Password tidak boleh kosong")
		return user, false, err
	}
	if len(userdata.Password) < 6 {
		err = fmt.Errorf("Password minimal 6 karakter")
		return user, false, err
	}
	if strings.Contains(userdata.Password, " ") {
		err = fmt.Errorf("Password tidak boleh mengandung spasi")
		return user, false, err
	}

	// Periksa apakah password sama dengan password lama
	if CheckPasswordHash(userdata.Password, userExists.Password) {
		err = fmt.Errorf("Password tidak boleh sama")
		return user, false, err
	}

	// Simpan pengguna ke basis data
	hash, _ := HashPassword(userdata.Password)
	userExists.Password = hash
	filter := bson.M{"username": userdata.Username}
	update := bson.M{
		"$set": bson.M{
			"email":    userdata.Email,
			"username": userdata.Username,
			"password": userExists.Password,
			"role":     "user",
		},
	}
	cols := db.Collection(col)
	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return user, false, err
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("PAssword tidak berhasil diupdate")
		return user, false, err
	}
	return user, true, nil
}

func DeleteUser(db *mongo.Database, col string, username string) error {
	cols := db.Collection(col)
	filter := bson.M{"username": username}
	result, err := cols.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for Username %s: %s", username, err.Error())
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no data found for Username %s", username)
	}
	return nil
}

func GetUserFromID(db *mongo.Database, col string, _id primitive.ObjectID) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}
	err = cols.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("GetUserFromID: %v\n", err)
	}
	return user, nil
}

func GetUserFromUsername(db *mongo.Database, col string, username string) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"username": username}
	err = cols.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("GetUserFromUsername: %v\n", err)
		return user, err
	}
	return user, nil
}

func GetUserFromEmail(db *mongo.Database, col string, email string) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"email": email}
	err = cols.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("GetUserFromEmail: %v\n", err)
	}
	return user, nil
}

func GetAllUser(db *mongo.Database, col string) (userlist []model.User) {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &userlist)
	if err != nil {
		fmt.Println(err)
	}
	return userlist
}

// catalog
func InsertCatalog(db *mongo.Database, col string, catalog model.Catalog) (insertedID primitive.ObjectID, err error) {
	insertedID, err = InsertOneDoc(db, col, catalog)
	if err != nil {
		fmt.Printf("InsertCatalog: %v\n", err)
	}
	return insertedID, err
}

func GetCatalogFromID(db *mongo.Database, col string, id primitive.ObjectID) (catalog model.Catalog) {
	cols := db.Collection(col)
	filter := bson.M{"_id": id}
	err := cols.FindOne(context.Background(), filter).Decode(&catalog)
	if err != nil {
		fmt.Printf("GetCatalogFromID: %v\n", err)
	}
	return catalog
}

func GetCatalogList(db *mongo.Database, col string) (cataloglist model.CatalogList) {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.Background(), &cataloglist.Items)
	if err != nil {
		fmt.Println(err)
	}
	return cataloglist
}

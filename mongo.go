package pasetobackendpasabar

import (
	"context"
	"fmt"
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

func GetNameAndPassowrd(mongoconn *mongo.Database, collection string) []User {
	user := atdb.GetAllDoc[[]User](mongoconn, collection)
	return user
}

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

func CreateUserAndAddedToeken(PASETOPRIVATEKEYENV string, mongoconn *mongo.Database, collection string, userdata User) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPassword(userdata.Password)
	if err != nil {
		return err
	}
	userdata.Password = hashedPassword

	// Insert the user data into the database
	atdb.InsertOneDoc(mongoconn, collection, userdata)

	// Create a token for the user
	tokenstring, err := watoken.Encode(userdata.Username, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		return err
	}
	userdata.Token = tokenstring

	// Update the user data in the database
	return atdb.ReplaceOneDoc(mongoconn, collection, bson.M{"username": userdata.Username}, userdata)
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

// func AuthenticateUserAndGenerateToken(privateKeyEnv string, mongoconn *mongo.Database, collection string, userdata User) (string, error) {
// 	// Cari pengguna berdasarkan nama pengguna
// 	username := userdata.Username
// 	password := userdata.Password
// 	userdata, err := FindUserByUsername(mongoconn, collection, username)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Memeriksa kata sandi
// 	if !CheckPasswordHash(password, userdata.Password) {
// 		return "", errors.New("Password salah") // Gantilah pesan kesalahan sesuai kebutuhan Anda
// 	}

// 	// Generate token untuk otentikasi
// 	tokenstring, err := watoken.Encode(username, os.Getenv(privateKeyEnv))
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenstring, nil
// }

// func FindUserByUsername(mongoconn *mongo.Database, collection string, username string) (User, error) {
// 	var user User
// 	filter := bson.M{"username": username}
// 	err := mongoconn.Collection(collection).FindOne(context.TODO(), filter).Decode(&user)
// 	if err != nil {
// 		return User{}, err
// 	}
// 	return user, nil
// }

// // create login using Private
// func CreateLogin(mongoconn *mongo.Database, collection string, userdata User) interface{} {
// 	// Hash the password before storing it
// 	hashedPassword, err := HashPassword(userdata.Password)
// 	if err != nil {
// 		return err
// 	}
// 	userdata.Password = hashedPassword
// 	// Create a token for the user
// 	tokenstring, err := watoken.Encode(userdata.Username, userdata.Private)
// 	if err != nil {
// 		return err
// 	}
// 	userdata.Token = tokenstring

// 	// Insert the user data into the database
// 	return atdb.InsertOneDoc(mongoconn, collection, userdata)
// }

// func MongoCreateConnection(MongoString, dbname string) *mongo.Database {
// 	MongoInfo := atdb.DBInfo{
// 		DBString: os.Getenv(MongoString),
// 		DBName:   dbname,
// 	}
// 	conn := atdb.MongoConnect(MongoInfo)
// 	return conn
// }

// func InsertOneDoc(db *mongo.Database, collection string, doc interface{}) (insertedID interface{}) {
// 	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), doc)
// 	if err != nil {
// 		fmt.Printf("InsertOneDoc: %v\n", err)
// 	}
// 	return insertResult.InsertedID
// }

// func GetAllUser(MongoConn *mongo.Database, colname string) []User {
// 	data := atdb.GetAllDoc[[]User](MongoConn, colname)
// 	return data
// }

// func GetOneUser(MongoConn *mongo.Database, colname string, userdata User) User {
// 	filter := bson.M{"username": userdata.Username}
// 	data := atdb.GetOneDoc[User](MongoConn, colname, filter)
// 	return data
// }

// func PasswordValidator(MongoConn *mongo.Database, colname string, userdata User) bool {
// 	filter := bson.M{"username": userdata.Username}
// 	data := atdb.GetOneDoc[User](MongoConn, colname, filter)
// 	hashChecker := CompareHashPass(userdata.Password, data.Password)
// 	return hashChecker
// }

// func InsertUserdata(MongoConn *mongo.Database, username, role, password string) (InsertedID interface{}) {
// 	req := new(User)
// 	req.Username = username
// 	req.Password = password
// 	req.Role = role
// 	return InsertOneDoc(MongoConn, "user", req)
// }

// func CompareUsername(MongoConn *mongo.Database, Colname, username string) bool {
// 	filter := bson.M{"username": username}
// 	err := atdb.GetOneDoc[User](MongoConn, Colname, filter)
// 	users := err.Username
// 	if users == "" {
// 		return false
// 	}
// 	return true
// }

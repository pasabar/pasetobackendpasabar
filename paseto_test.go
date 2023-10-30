package pasetobackendpasabar

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateNewUserRole(t *testing.T) {
	var userdata User
	userdata.Username = "pasabar"
	userdata.Password = "lodons"
	userdata.Role = "admin"
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	CreateNewUserRole(mconn, "user", userdata)
}

func TestDeleteUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	var userdata User
	userdata.Username = "lolz"
	DeleteUser(mconn, "user", userdata)
}

func CreateNewUserToken(t *testing.T) {
	var userdata User
	userdata.Username = "pasabar"
	userdata.Password = "lodons"
	userdata.Role = "admin"

	// Create a MongoDB connection
	mconn := SetConnection("MONGOSTRING", "pasabar13")

	// Call the function to create a user and generate a token
	err := CreateUserAndAddToken("your_private_key_env", mconn, "user", userdata)

	if err != nil {
		t.Errorf("Error creating user and token: %v", err)
	}
}

func TestGFCPostHandlerUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	var userdata User
	userdata.Username = "pasabar"
	userdata.Password = "lodons"
	userdata.Role = "admin"
	CreateNewUserRole(mconn, "user", userdata)
}

func TestGeneratePasswordHash(t *testing.T) {
	password := "brazilia"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)
	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("musa", privateKey)
	fmt.Println(hasil, err)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	var userdata User
	userdata.Username = "musa"
	userdata.Password = "brazilia"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	var userdata User
	userdata.Username = "musa"
	userdata.Password = "brazilia"

	anu := IsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}

func TestUserFix(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	var userdata User
	userdata.Username = "pasabar"
	userdata.Password = "lodons"
	userdata.Role = "admin"
	CreateUser(mconn, "user", userdata)
}

// var privatekey = "privatekey"
// var publickeyb = "publickey"
// var encode = "encode"

// func TestGenerateKeyPASETO(t *testing.T) {
// 	privateKey, publicKey := watoken.GenerateKey()
// 	fmt.Println(privateKey)
// 	fmt.Println(publicKey)
// 	hasil, err := watoken.Encode("pasabar", privateKey)
// 	fmt.Println(hasil, err)
// }

// func TestHashPass(t *testing.T) {
// 	password := "pasabarcoba"

// 	Hashedpass, err := HashPass(password)
// 	fmt.Println("error : ", err)
// 	fmt.Println("Hash : ", Hashedpass)
// }

// func TestHashFunc(t *testing.T) {
// 	conn := MongoCreateConnection("MONGOSTRING", "pasabar3")
// 	userdata := new(User)
// 	userdata.Username = "pasabar"
// 	userdata.Password = "pasabarcoba"

// 	data := GetOneUser(conn, "user", User{
// 		Username: userdata.Username,
// 		Password: userdata.Password,
// 	})
// 	fmt.Printf("%+v", data)
// 	fmt.Println(" ")
// 	hashpass, _ := HashPass(userdata.Password)
// 	fmt.Println("Hasil hash : ", hashpass)
// 	compared := CompareHashPass(userdata.Password, data.Password)
// 	fmt.Println("result : ", compared)
// }

// func TestTokenEncoder(t *testing.T) {
// 	conn := MongoCreateConnection("MONGOSTRING", "pasabar3")
// 	privateKey, publicKey := watoken.GenerateKey()
// 	userdata := new(User)
// 	userdata.Username = "pasabar"
// 	userdata.Password = "pasabarcoba"

// 	data := GetOneUser(conn, "user", User{
// 		Username: userdata.Username,
// 		Password: userdata.Password,
// 	})
// 	fmt.Println("Private Key : ", privateKey)
// 	fmt.Println("Public Key : ", publicKey)
// 	fmt.Printf("%+v", data)
// 	fmt.Println(" ")

// 	encode := TokenEncoder(data.Username, privateKey)
// 	fmt.Printf("%+v", encode)
// }

// func TestInsertUserdata(t *testing.T) {
// 	conn := MongoCreateConnection("MONGOSTRING", "pasabar3")
// 	password, err := HashPass("pasabar13")
// 	fmt.Println("err", err)
// 	data := InsertUserdata(conn, "pasabar", "role", password)
// 	fmt.Println(data)
// }

// func TestDecodeToken(t *testing.T) {
// 	deco := watoken.DecodeGetId("public",
// 		"token")
// 	fmt.Println(deco)
// }

// func TestCompareUsername(t *testing.T) {
// 	conn := MongoCreateConnection("MONGOSTRING", "pasabar3")
// 	deco := watoken.DecodeGetId("public",
// 		"token")
// 	compare := CompareUsername(conn, "user", deco)
// 	fmt.Println(compare)
// }

// func TestEncodeWithRole(t *testing.T) {
// 	privateKey, publicKey := watoken.GenerateKey()
// 	role := "admin"
// 	username := "pasabar"
// 	encoder, err := EncodeWithRole(role, username, privateKey)

// 	fmt.Println(" error :", err)
// 	fmt.Println("Private :", privateKey)
// 	fmt.Println("Public :", publicKey)
// 	fmt.Println("encode: ", encoder)

// }

// func TestDecoder2(t *testing.T) {
// 	pay, err := Decoder(publickeyb, encode)
// 	user, _ := DecodeGetUser(publickeyb, encode)
// 	role, _ := DecodeGetRole(publickeyb, encode)
// 	use, ro := DecodeGetRoleandUser(publickeyb, encode)
// 	fmt.Println("user :", user)
// 	fmt.Println("role :", role)
// 	fmt.Println("user and role :", use, ro)
// 	fmt.Println("err : ", err)
// 	fmt.Println("payload : ", pay)
// }

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

// func TestDeleteUser(t *testing.T) {
// 	mconn := SetConnection("MONGOSTRING", "pasabar13")
// 	var userdata User
// 	userdata.Username = "lolz"
// 	DeleteUser(mconn, "user", userdata)
// }

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

func TestCatalog(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	var catalogdata Catalog
	catalogdata.Nomorid = 1
	catalogdata.Title = "garuts"
	catalogdata.Description = "membahana"
	catalogdata.Image = "https://images3.alphacoders.com/165/thumb-1920-165265.jpg"
	CreateNewCatalog(mconn, "catalog", catalogdata)
}

func TestAllCatalog(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	catalog := GetAllCatalog(mconn, "catalog")
	fmt.Println(catalog)
}

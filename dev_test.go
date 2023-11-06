package pasetobackendpasabar

import (
	"fmt"
	"testing"

	model "github.com/pasabar/pasetobackendpasabar/model"
	modul "github.com/pasabar/pasetobackendpasabar/modul"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mconn = SetConnection("MONGOSTRING", "pasabar13")

// user
func TestRegister(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	var userdata User
	userdata.Email = "pasabarsky@gmail.com"
	userdata.Username = "pasabarsky"
	userdata.Role = "admin"
	userdata.Password = "lodonsky"

	nama := InsertUser(mconn, "user", userdata)
	fmt.Println(nama)
}

// test login
// func TestLogIn(t *testing.T) {
// 	var userdata model.User
// 	userdata.Username = "dapskuy"
// 	userdata.Password = "kepoah"
// 	user, status, err := modul.LogIn(mconn, "user", userdata)
// 	fmt.Println("Status", status)
// 	if err != nil {
// 		t.Errorf("Error logging in user: %v", err)
// 	} else {
// 		fmt.Println("Login success", user)
// 	}
// }

// test change password
func TestChangePassword(t *testing.T) {
	username := "ujang"
	oldpassword := "ujang123"
	newpassword := "ujanggeming"

	var userdata model.User
	userdata.Username = username
	userdata.Password = newpassword

	userdata, status, err := modul.ChangePassword(mconn, "user", username, oldpassword, newpassword)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error changing password: %v", err)
	} else {
		fmt.Println("Password change success for user", userdata)
	}
}

// test delete user
func TestDeleteUser(t *testing.T) {
	username := "kamir"

	err := modul.DeleteUser(mconn, "user", username)
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	} else {
		fmt.Println("Delete user success")
	}

	_, err = modul.GetUserFromUsername(mconn, "user", username)
	if err == nil {
		fmt.Println("Data masih ada")
	}
}

func TestGetUserFromID(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("")
	anu, _ := modul.GetUserFromID(mconn, "user", id)
	fmt.Println(anu)
}

func TestGetUserFromUsername(t *testing.T) {
	anu, err := modul.GetUserFromUsername(mconn, "user", "pasabarsky")
	if err != nil {
		t.Errorf("Error getting user: %v", err)
		return
	}
	fmt.Println(anu)
}

func TestGetAllUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	anu := modul.GetAllUser(mconn, "user")
	fmt.Println(anu)
}

// catalog
func TestInsertCatalog(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	var catalogdata model.Catalog
	catalogdata.Catalog = "Tasik"
	catalogdata.Description = "Ada iqbal disana jadi mengmusa"
	catalogdata.IsDone = true

	nama, err := modul.InsertCatalog(mconn, "catalog", catalogdata)
	if err != nil {
		t.Errorf("Error inserting catalog: %v", err)
	}
	fmt.Println(nama)
}

func TestGetCatalogFromID(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	id, _ := primitive.ObjectIDFromHex("6548c4a7b0a03450a264258c")
	anu := modul.GetCatalogFromID(mconn, "catalog", id)
	fmt.Println(anu)
}

func TestGetCatalogList(t *testing.T) {
	anu := modul.GetCatalogList(mconn, "user")
	fmt.Println(anu)
}

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
	var data model.User
	data.ID = primitive.NewObjectID()
	data.Email = "pasabar13@gmail.com"
	data.Username = "pasabar09123"
	data.Role = "user"
	data.Password = "secret"

	err := modul.Register(mconn, "user", data)
	if err != nil {
		t.Errorf("Error registering user: %v", err)
	} else {
		fmt.Println("Register success", data)
	}
}

func TestLogIn(t *testing.T) {
	var data model.User
	data.Username = "iyas"
	data.Password = "secret"

	user, status, err := modul.LogIn(mconn, "user", data)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error logging in user: %v", err)
	} else {
		fmt.Println("Login success", user)
	}
}

// test change password
func TestChangePassword(t *testing.T) {
	var data model.User
	data.Email = "iyas@gmail.com" // email tidak diubah
	data.Username = "iyas"        // username tidak diubah
	data.Role = "user"            // role tidak diubah

	data.Password = "ganteng"

	// username := "ryaaspo123"

	_, status, err := modul.ChangePassword(mconn, "user", data)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error updateting document: %v", err)
	} else {
		fmt.Println("Password berhasil diubah dengan username:", data.Username)
	}
}

func TestUpdateUser(t *testing.T) {
	var data model.User
	data.Email = "iyas@gmail.com"
	data.Username = "iyas"
	data.Role = "user"

	data.Password = "secret" // password tidak diubah

	id, err := primitive.ObjectIDFromHex("654a8d045fed94712c1f3a89")
	data.ID = id
	if err != nil {
		fmt.Printf("Data tidak berhasil diubah")
	} else {

		_, status, err := modul.UpdateUser(mconn, "user", data)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error updateting document: %v", err)
		} else {
			fmt.Printf("Data berhasil diubah untuk id: %s\n", id)
		}
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
	id, _ := primitive.ObjectIDFromHex("654a8cd3884072873a89a682")
	anu, _ := modul.GetUserFromID(mconn, "user", id)
	fmt.Println(anu)
}

func TestGetUserFromUsername(t *testing.T) {
	anu, err := modul.GetUserFromUsername(mconn, "user", "pasabar09123")
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
	catalogdata.Title = "Pangandaran"
	catalogdata.Description = "sangat indah banget"
	catalogdata.Image = "https://media.suara.com/pictures/653x366/2022/04/30/75648-taman-langit-pangalengan.webp"
	catalogdata.IsDone = true

	nama, err := modul.InsertCatalog(mconn, "catalog", catalogdata)
	if err != nil {
		t.Errorf("Error inserting catalog: %v", err)
	}
	fmt.Println(nama)
}

func TestGetCatalogFromID(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pasabar13")
	id, _ := primitive.ObjectIDFromHex("654a8f229fb816979403be96")
	anu := modul.GetCatalogFromID(mconn, "catalog", id)
	fmt.Println(anu)
}

func TestGetCatalogList(t *testing.T) {
	anu := modul.GetCatalogList(mconn, "catalog")
	fmt.Println(anu)
}

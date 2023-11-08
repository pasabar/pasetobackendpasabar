package modul

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	model "github.com/pasabar/pasetobackendpasabar/model"
	"github.com/whatsauth/watoken"
)

func GCFHandler(MONGOCONNSTRINGENV, dbname, col string) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	data := GetAllUser(mconn, col)
	return GCFReturnStruct(data)
}

func GCFHandlerGetUser(MONGOCONNSTRINGENV, dbname, col string, username string) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	data, err := GetUserFromUsername(mconn, col, username)
	if err != nil {
		return GCFReturnStruct(err)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerRegister(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	}
	err = Register(mconn, collectionname, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Register success" + datauser.Username

	return GCFReturnStruct(Response)
}

func GCFHandlerLogIn(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	}
	user, status, err := LogIn(mconn, collectionname, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Response.Message = "Gagal Encode Token :" + err.Error()
	} else {
		Response.Message = "Login success" + user.Username + strconv.FormatBool(status)
		Response.Token = tokenstring
	}

	return GCFReturnStruct(Response)
}

func GCFHandlerUpdateUser(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "Bad Request: error parsing application/json: " + err.Error()
	}
	user, status, err := UpdateUser(mconn, collectionname, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Update success " + user.Username + " " + strconv.FormatBool(status)

	return GCFReturnStruct(Response)
}

func GCFHandlerChangePassword(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	}
	user, status, err := ChangePassword(mconn, collectionname, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Password change success for user" + user.Username + strconv.FormatBool(status)

	return GCFReturnStruct(Response)
}

func GCFHandlerDeleteUser(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	}
	err = DeleteUser(mconn, collectionname, datauser.Username)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Delete user success"

	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

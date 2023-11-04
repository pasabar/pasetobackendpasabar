package pasetobackendpasabar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func GCFFindUserByID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		return err.Error()
	}
	user := FindUser(mconn, collectionname, datauser)
	return GCFReturnStruct(user)
}

func GCFFindUserByName(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		return err.Error()
	}

	// Jika username kosong, maka respon "false" dan data tidak ada
	if datauser.Username == "" {
		return "false"
	}

	// Jika ada username, mencari data pengguna
	user := FindUser(mconn, collectionname, datauser)

	// Jika data pengguna ditemukan, mengembalikan data pengguna dalam format yang sesuai
	if user != (User{}) {
		return GCFReturnStruct(user)
	}

	// Jika tidak ada data pengguna yang ditemukan, mengembalikan "false" dan data tidak ada
	return "false"
}

func GCFDeleteHandler(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		return err.Error()
	}
	DeleteUser(mconn, collectionname, datauser)
	return GCFReturnStruct(datauser)
}

func GCFUpdateHandler(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		return err.Error()
	}
	ReplaceOneDoc(mconn, collectionname, bson.M{"username": datauser.Username}, datauser)
	return GCFReturnStruct(datauser)
}

func GCFCreateHandler(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		return err.Error()
	}

	// Hash the password before storing it
	hashedPassword, hashErr := HashPassword(datauser.Password)
	if hashErr != nil {
		return hashErr.Error()
	}
	datauser.Password = hashedPassword

	createErr := CreateNewUserRole(mconn, collectionname, datauser)
	fmt.Println(createErr)

	return GCFReturnStruct(datauser)
}

func GFCPostHandlerUser(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false

	// Mendapatkan data yang diterima dari permintaan HTTP POST
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		// Menggunakan variabel MONGOCONNSTRINGENV untuk string koneksi MongoDB
		mongoConnStringEnv := MONGOCONNSTRINGENV

		mconn := SetConnection(mongoConnStringEnv, dbname)

		// Lakukan pemeriksaan kata sandi menggunakan bcrypt
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			Response.Message = "Selamat Datang"
		} else {
			Response.Message = "Password Salah"
		}
	}

	// Mengirimkan respons sebagai JSON
	responseJSON, _ := json.Marshal(Response)
	return string(responseJSON)
}

func GCFPostHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang"
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func GCFLoginTest(username, password, MONGOCONNSTRINGENV, dbname, collectionname string) bool {
	// Membuat koneksi ke MongoDB
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	// Mencari data pengguna berdasarkan username
	filter := bson.M{"username": username}
	collection := collectionname
	res := atdb.GetOneDoc[User](mconn, collection, filter)

	// Memeriksa apakah pengguna ditemukan dalam database
	if res == (User{}) {
		return false
	}

	// Memeriksa apakah kata sandi cocok
	return CheckPasswordHash(password, res.Password)
}

func Login(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp Credential
	mconn := SetConnection(MongoEnv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, Colname, datauser) {
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(Privatekey))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				resp.Status = true
				resp.Message = "Selamat Datang"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "Password Salah"
		}
	}
	return GCFReturnStruct(resp)
}

func ReturnStringStruct(Data any) string {
	jsonee, _ := json.Marshal(Data)
	return string(jsonee)
}

func Register(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	userdata := new(User)
	resp.Status = false
	conn := GetConnectionMongo(Mongoenv, dbname)
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		hash, err := HashPassword(userdata.Password)
		if err != nil {
			resp.Message = "Gagal Hash Password" + err.Error()
		}
		InsertUserdata(conn, userdata.Username, userdata.Role, hash)
		resp.Message = "Berhasil Input data"
	}
	response := ReturnStringStruct(resp)
	return response
}

// add encrypt password to database and tokenstring
// func GCFCreateHandler(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {

// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var datauser User
// 	err := json.NewDecoder(r.Body).Decode(&datauser)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	CreateNewUserRole(mconn, collectionname, datauser)
// 	return GCFReturnStruct(datauser)
// }

// func GCFCreateHandlerTokenPaseto(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var datauser User
// 	err := json.NewDecoder(r.Body).Decode(&datauser)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	hashedPassword, hashErr := HashPassword(datauser.Password)
// 	if hashErr != nil {
// 		return hashErr.Error()
// 	}
// 	datauser.Password = hashedPassword
// 	CreateNewUserRole(mconn, collectionname, datauser)
// 	tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
// 	if err != nil {
// 		return err.Error()
// 	}
// 	datauser.Token = tokenstring
// 	return GCFReturnStruct(datauser)
// }

// func GCFCreateAccountAndToken(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var datauser User
// 	err := json.NewDecoder(r.Body).Decode(&datauser)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	hashedPassword, hashErr := HashPassword(datauser.Password)
// 	if hashErr != nil {
// 		return hashErr.Error()
// 	}
// 	datauser.Password = hashedPassword
// 	CreateUserAndAddedToeken(PASETOPRIVATEKEYENV, mconn, collectionname, datauser)
// 	return GCFReturnStruct(datauser)
// }

// func GCFCreateRegister(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var userdata User
// 	err := json.NewDecoder(r.Body).Decode(&userdata)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	CreateUser(mconn, collectionname, userdata)
// 	return GCFReturnStruct(userdata)
// }

// func GCFCreateTokenAndSaveToDB(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) (string, error) {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

// 	// Inisialisasi variabel datauser
// 	var datauser User

// 	// Membaca data JSON dari permintaan HTTP ke dalam datauser
// 	if err := json.NewDecoder(r.Body).Decode(&datauser); err != nil {
// 		return "", err // Mengembalikan kesalahan langsung
// 	}

// 	// Generate a token for the user
// 	tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
// 	if err != nil {
// 		return "", err // Mengembalikan kesalahan langsung
// 	}
// 	datauser.Token = tokenstring

// 	// Simpan pengguna ke dalam basis data
// 	if err := atdb.InsertOneDoc(mconn, collectionname, datauser); err != nil {
// 		return tokenstring, nil // Mengembalikan kesalahan langsung
// 	}

// 	return tokenstring, nil // Mengembalikan token dan nil untuk kesalahan jika sukses
// }

// func GCFLoginAfterCreate(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var userdata User
// 	err := json.NewDecoder(r.Body).Decode(&userdata)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	if IsPasswordValid(mconn, collectionname, userdata) {
// 		tokenstring, err := watoken.Encode(userdata.Username, os.Getenv("PASETOPRIVATEKEYENV"))
// 		if err != nil {
// 			return err.Error()
// 		}
// 		userdata.Token = tokenstring
// 		return GCFReturnStruct(userdata)
// 	} else {
// 		return "Password Salah"
// 	}
// }

// func GCFLoginAfterCreater(MONGOCONNSTRINGENV, dbname, collectionname, privateKeyEnv string, r *http.Request) (string, error) {
// 	// Ambil data pengguna dari request, misalnya dari body JSON atau form data.
// 	var userdata User
// 	// Implement the logic to extract user data from the request (r) here.

// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

// 	// Lakukan otentikasi pengguna yang baru saja dibuat.
// 	token, err := AuthenticateUserAndGenerateToken(privateKeyEnv, mconn, collectionname, userdata)
// 	if err != nil {
// 		return "", err
// 	}
// 	return token, nil
// }

// func GCFLoginAfterCreatee(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var userdata User
// 	err := json.NewDecoder(r.Body).Decode(&userdata)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	if IsPasswordValid(mconn, collectionname, userdata) {
// 		// Password is valid, return a success message or some other response.
// 		return "Login successful"

// 	} else {
// 		// Password is not valid, return an error message.
// 		return "Password Salah"
// 	}
// }

// func GCFLoginAfterCreateee(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var userdata User
// 	err := json.NewDecoder(r.Body).Decode(&userdata)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	if IsPasswordValid(mconn, collectionname, userdata) {
// 		// Password is valid, construct and return the GCFReturnStruct.
// 		response := CreateResponse(true, "Berhasil Login", userdata)
// 		return GCFReturnStruct(response) // Return GCFReturnStruct directly
// 	} else {
// 		// Password is not valid, return an error message.
// 		return "Password Salah"
// 	}
// }
// func GCFLoginAfterCreateeee(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var userdata User
// 	err := json.NewDecoder(r.Body).Decode(&userdata)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	if IsPasswordValid(mconn, collectionname, userdata) {
// 		// Password is valid, return a success message or some other response.
// 		return GCFReturnStruct(userdata)
// 	} else {
// 		// Password is not valid, return an error message.
// 		return "Password Salah"
// 	}
// }

// func GCFLoginFixx(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var userdata User
// 	err := json.NewDecoder(r.Body).Decode(&userdata)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	if IsPasswordValid(mconn, collectionname, userdata) {
// 		// Password is valid, construct and return the GCFReturnStruct.
// 		userMap := map[string]interface{}{
// 			"Username": userdata.Username,
// 			"Password": userdata.Password,
// 			"Private":  userdata.Private,
// 			"Publick":  userdata.Public,
// 		}
// 		response := CreateResponse(true, "Berhasil Login", userMap)
// 		return GCFReturnStruct(response) // Return GCFReturnStruct directly
// 	} else {
// 		// Password is not valid, return an error message.
// 		return "Password Salah"
// 	}
// }

// func Login(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
// 	var resp Credential
// 	mconn := MongoCreateConnection(MongoEnv, dbname)
// 	var datauser peda.User
// 	err := json.NewDecoder(r.Body).Decode(&datauser)
// 	if err != nil {
// 		resp.Message = "error parsing application/json: " + err.Error()
// 	} else {
// 		if peda.IsPasswordValid(mconn, Colname, datauser) {
// 			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(Privatekey))
// 			if err != nil {
// 				resp.Message = "Gagal Encode Token : " + err.Error()
// 			} else {
// 				resp.Status = true
// 				resp.Message = "Selamat Datang"
// 				resp.Token = tokenstring
// 			}
// 		} else {
// 			resp.Message = "Password Salah"
// 		}
// 	}
// 	return ReturnStringStruct(resp)
// }

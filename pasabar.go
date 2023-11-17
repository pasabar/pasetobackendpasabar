package pasetobackendpasabar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

// <--- ini catalog --->

// catalog post
func GCFCreateCatalog(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datacatalog Catalog
	err := json.NewDecoder(r.Body).Decode(&datacatalog)
	if err != nil {
		return err.Error()
	}
	if err := CreateCatalog(mconn, collectionname, datacatalog); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Create Catalog", datacatalog))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Create Catalog", datacatalog))
	}
}

// delete catalog
func GCFDeleteCatalog(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datacatalog Catalog
	err := json.NewDecoder(r.Body).Decode(&datacatalog)
	if err != nil {
		return err.Error()
	}

	if err := DeleteCatalog(mconn, collectionname, datacatalog); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Delete Catalog", datacatalog))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Delete Catalog", datacatalog))
	}
}

// update catalog
func GCFUpdateCatalog(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datacatalog Catalog
	err := json.NewDecoder(r.Body).Decode(&datacatalog)
	if err != nil {
		return err.Error()
	}

	if err := UpdatedCatalog(mconn, collectionname, bson.M{"id": datacatalog.ID}, datacatalog); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Update Catalog", datacatalog))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Update Catalog", datacatalog))
	}
}

// get all catalog
func GCFGetAllCatalog(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datacatalog := GetAllCatalog(mconn, collectionname)
	if datacatalog != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Catalog", datacatalog))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Catalog", datacatalog))
	}
}

// get all catalog by id
func GCFGetAllCatalogID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datacatalog Catalog
	err := json.NewDecoder(r.Body).Decode(&datacatalog)
	if err != nil {
		return err.Error()
	}

	catalog := GetAllCatalogID(mconn, collectionname, datacatalog)
	if catalog != (Catalog{}) {
		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Catalog", datacatalog))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Catalog", datacatalog))
	}
}

// <--- ini about --->

// about post
func GCFCreateAbout(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var dataabout About
	err := json.NewDecoder(r.Body).Decode(&dataabout)
	if err != nil {
		return err.Error()
	}

	if err := CreateAbout(mconn, collectionname, dataabout); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Create About", dataabout))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Create About", dataabout))
	}
}

// delete about
func GCFDeleteAbout(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var dataabout About
	err := json.NewDecoder(r.Body).Decode(&dataabout)
	if err != nil {
		return err.Error()
	}

	if err := DeleteAbout(mconn, collectionname, dataabout); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Delete About", dataabout))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Delete About", dataabout))
	}
}

// update about
func GCFUpdateAbout(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var dataabout About
	err := json.NewDecoder(r.Body).Decode(&dataabout)
	if err != nil {
		return err.Error()
	}

	if err := UpdatedAbout(mconn, collectionname, bson.M{"id": dataabout.ID}, dataabout); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Update About", dataabout))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Update About", dataabout))
	}
}

// get all about
func GCFGetAllAbout(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	dataabout := GetAllAbout(mconn, collectionname)
	if dataabout != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All About", dataabout))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All About", dataabout))
	}
}

// event tour function

func CreateTour(mongoconn *mongo.Database, collection string, tourdata Tour) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, tourdata)
}

func DeleteTour(mongoconn *mongo.Database, collection string, tourdata Tour) interface{} {
	filter := bson.M{"id": tourdata.ID}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedTour(mongoconn *mongo.Database, collection string, filter bson.M, tourdata Tour) interface{} {
	filter = bson.M{"id": tourdata.ID}
	return atdb.ReplaceOneDoc(mongoconn, collection, filter, tourdata)
}

func GetAllTour(mongoconn *mongo.Database, collection string) []Tour {
	tour := atdb.GetAllDoc[[]Tour](mongoconn, collection)
	return tour
}

func GetAllTourId(mongoconn *mongo.Database, collection string, tourdata Tour) []Tour {
	filter := bson.M{"id": tourdata.ID}
	tour := atdb.GetOneDoc[[]Tour](mongoconn, collection, filter)
	return tour
}

// event hotelresto function

func CreateHotelResto(mongoconn *mongo.Database, collection string, hotelrestodata HotelResto) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, hotelrestodata)
}

func DeleteHotelResto(mongoconn *mongo.Database, collection string, hotelrestodata HotelResto) interface{} {
	filter := bson.M{"id": hotelrestodata.ID}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedHotelResto(mongoconn *mongo.Database, collection string, filter bson.M, hotelrestodata HotelResto) interface{} {
	filter = bson.M{"id": hotelrestodata.ID}
	return atdb.ReplaceOneDoc(mongoconn, collection, filter, hotelrestodata)
}

func GetAllHotelResto(mongoconn *mongo.Database, collection string) []HotelResto {
	hotelresto := atdb.GetAllDoc[[]HotelResto](mongoconn, collection)
	return hotelresto
}

func GetAllHotelRestoId(mongoconn *mongo.Database, collection string, hotelrestodata HotelResto) []HotelResto {
	filter := bson.M{"id": hotelrestodata.ID}
	hotelresto := atdb.GetOneDoc[[]HotelResto](mongoconn, collection, filter)
	return hotelresto
}

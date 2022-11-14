package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/averageNetAdmin/go_test_task/internal/claim"
	"github.com/averageNetAdmin/go_test_task/internal/mydb"
	"github.com/averageNetAdmin/go_test_task/internal/response"
	"github.com/golang-jwt/jwt"
)

const logFile = "./go_test_task.log"
const PORT = "8080"

var errors = make(chan error, 10)

func registrate(w http.ResponseWriter, r *http.Request) {
	var body = make([]byte, 1024)
	n, err := r.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		errors <- err
		return
	}
	var user = make(map[string]interface{})
	err = json.Unmarshal(body[:n], &user)
	if err != nil {
		errors <- err
		return
	}
	name := user["name"].(string)
	password := user["password"].(string)

	db, _ := mydb.Init()
	err = db.AddUser(name, password)
	if err != nil {
		errors <- err
		w.Write((&response.Response{
			Error: err.Error(),
			Data:  "",
		}).ToByte())
	}
	w.Write((&response.Response{
		Error: "",
		Data:  "",
	}).ToByte())
}

func login(w http.ResponseWriter, r *http.Request) {
	var body = make([]byte, 1024)
	n, err := r.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		errors <- err
		return
	}
	var user = make(map[string]interface{})
	err = json.Unmarshal(body[:n], &user)
	if err != nil {
		errors <- err
		w.Write((&response.Response{
			Error: err.Error(),
			Data:  "",
		}).ToByte())
		return
	}
	name := user["name"].(string)
	password := user["password"].(string)

	db, _ := mydb.Init()
	isExist := db.CheckPassword(name, password)
	if !isExist {
		w.Write((&response.Response{
			Error: "username or password are invalid",
			Data:  "",
		}).ToByte())
		return
	}

	claim := &claim.UserClaim{
		StandardClaims: jwt.StandardClaims{},
		Username:       name,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	hasher := sha1.New()
	hasher.Write([]byte(name + password))
	sha := fmt.Sprintf("%x", hasher.Sum(nil))

	token, err := t.SignedString([]byte(sha))
	if err != nil {
		w.Write((&response.Response{
			Error: "invalid data",
			Data:  "",
		}).ToByte())
		errors <- err
		return
	}
	db.SetToken(name, sha)
	w.Write((&response.Response{
		Error: "",
		Data:  token,
	}).ToByte())
}

func getTokenFromRequest(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("empty token")
	}
	authParts := strings.Split(auth, " ")
	if len(authParts) < 2 || authParts[0] == "" {
		return "", fmt.Errorf("empty token")
	}
	return authParts[1], nil
}

func chechToken(t string) bool {
	token, err := jwt.ParseWithClaims(t, &claim.UserClaim{}, func(t *jwt.Token) (interface{}, error) {
		db, _ := mydb.Init()
		claim := t.Claims.(*claim.UserClaim)
		key, err := db.GetToken(claim.Username)
		return []byte(key), err
	})
	if err != nil || !token.Valid {
		return false
	}
	return true
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	t, err := getTokenFromRequest(r)
	if err != nil {
		w.Write((&response.Response{
			Error: "please login or registrate",
			Data:  "",
		}).ToByte())
		return
	}
	fmt.Println(t)
	if !chechToken(t) {
		w.Write((&response.Response{
			Error: "invalid token",
			Data:  "",
		}).ToByte())
		return
	}
	// handle request
	w.Write((&response.Response{
		Error: "",
		Data:  "you are entered",
	}).ToByte())

}

func logout(w http.ResponseWriter, r *http.Request) {
	var body = make([]byte, 1024)
	n, err := r.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		w.Write((&response.Response{
			Error: err.Error(),
			Data:  "",
		}).ToByte())
		errors <- err
		return
	}
	var user = make(map[string]interface{})
	err = json.Unmarshal(body[:n], &user)
	if err != nil {
		errors <- err
		w.Write((&response.Response{
			Error: err.Error(),
			Data:  "",
		}).ToByte())
		return
	}
	name := user["name"].(string)

	db, _ := mydb.Init()
	db.SetToken(name, "")
	w.Write((&response.Response{
		Error: "",
		Data:  "you are logged out",
	}).ToByte())
}

func errorHandler() {
	f, err := os.Create(logFile)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		err := <-errors
		f.Write([]byte(err.Error() + "\n"))
	}
}

func main() {
	go errorHandler()

	http.HandleFunc("/api", tokenHandler)
	http.HandleFunc("/registrate", registrate)
	http.HandleFunc("/log_out", logout)
	http.HandleFunc("/log_in", login)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.ListenAndServe(":"+PORT, nil)
}

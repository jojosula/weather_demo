package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func AllHandlerSetup(r *mux.Router) {
	// login
	{
		r.HandleFunc("/v1/login", HandlerLogin).
			Methods(http.MethodPost)
		r.HandleFunc("/v1/logout", HandlerLogout).
			Methods(http.MethodPost)
	}
	// weather
	{
		r.HandleFunc("/v1/weather", AuthMiddleware(HandlerGetWeather)).
			Methods(http.MethodGet)
	}
}

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	var loginParam LoginRequest
	resp := Response{
		Code:    -1,
		Message: "login failed",
	}
	defer func() {
		b, _ := json.Marshal(resp)
		fmt.Fprintf(w, "%s", string(b))
	}()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginParam)
	if err != nil {
		panic(err)
	}
	log.Printf("login parameter:%v", loginParam)

	user := GetUser(loginParam.Username, loginParam.Password)
	if len(user.Username) > 0 {
		usertoken := GetUserToken(user.ID)
		// create when token not exist
		if len(usertoken.AccessToken) == 0 {
			usertoken = CreateUserToken(user.ID)
		}
		resp = Response{
			Code:    0,
			Message: "login success",
			Result:  usertoken,
		}
	}
}

func HandlerLogout(w http.ResponseWriter, r *http.Request) {
	var resp Response
	defer func() {
		b, _ := json.Marshal(resp)
		fmt.Fprintf(w, "%s", string(b))
	}()

	accessToken := r.Header.Get("Authorization")
	log.Printf("logout parameter:%v", accessToken)

	DeleteUserToken(accessToken)
	resp = Response{
		Code:    0,
		Message: "logout success",
	}
}

func HandlerGetWeather(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Code:    -1,
		Message: "failed",
	}
	defer func() {
		b, _ := json.Marshal(resp)
		fmt.Fprintf(w, "%s", string(b))
	}()

	keys := r.URL.Query()
	location := keys.Get("location")
	log.Printf("get weather parameter:%v", location)

	data := GetWeatherData(location)
	if len(data.Data) > 0 {
		var result interface{}
		err := json.Unmarshal([]byte(data.Data), &result)
		if err != nil {
			log.Fatal(err)
			return
		}

		resp = Response{
			Code:    0,
			Message: "success",
			Result:  result,
		}
	}
}

package main

import (
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type UserToken struct {
	UserID      int       `json:"user_id"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
}

type WeatherData struct {
	Location  string    `json:"location"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result,omitempty"`
}

type WeatherAPIResponse struct {
	Success string      `json:"success"`
	Result  interface{} `json:"result"`
	Records interface{} `json:"records"`
}

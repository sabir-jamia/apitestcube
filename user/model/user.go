package model

import "time"

type User struct  {
	Id                int       `json:"id"`
	Email             string    `json:"email"`
	Username          string    `json:"username"`
	Password          string    `json:"password"`
	RegistrationDate time.Time `json:"registration_date"`
	ClientId         int       `json:"client_id"`
	Status            bool      `json:"status"`
}

type ClientUser struct {
	Id          int       `json:"id"`
	TestcubeId int       `json:"testcue_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email       string    `json:"email"`
	CreatedOn  time.Time `json:"created_on"`
	UpdatedOn  time.Time `json:"updated_on"`
	Status      bool      `json:"status"`
	Language    string    `json:"language"`
}

func (ClientUser) TableName() string {
	return "users"
}

type UserProfile struct {
	Id        int    `json:"id"`
	UserId    int    `json:"testcube_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Contact   int    `json:"email"`
	Photo     string `json:"photo"`
	Language  string `json:"language"`
	Theme     int    `json:"theme"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}
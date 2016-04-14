package model

import "time"

type User struct  {
	Id               int       `json:"id"`
	Email            string    `json:"email"`
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	RegistrationDate time.Time `json:"registrationDate"`
	ClientId         int       `json:"clientId"`
	Status           bool      `json:"status"`
}

type ClientUser struct {
	Id         int       `json:"id"`
	TestcubeId int       `json:"testcubeId"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Email      string    `json:"email"`
	CreatedOn  time.Time `json:"createdOn"`
	UpdatedOn  time.Time `json:"updatedOn"`
	Status     bool      `json:"status"`
	Language   string    `json:"language"`
}

func (ClientUser) TableName() string {
	return "users"
}

type UserProfile struct {
	Id        int    `json:"id"`
	UserId    int    `json:"testcubeId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Contact   int    `json:"email"`
	Photo     string `json:"photo"`
	Language  string `json:"language"`
	Theme     int    `json:"theme"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}
package model

import "time"

type User struct  {
	Id                int       `json:"id"`
	Email             string    `json:"email"`
	Username          string    `json:"username"`
	Password          string    `json:"password"`
	Registration_date time.Time `json:"registration_date"`
	Client_id         int       `json:"client_id"`
	Status            int       `json:"status"`
}

type ClientUser struct {
	Id          int       `json:"id"`
	Testcube_id int       `json:"testcue_id"`
	First_name  string    `json:"first_name"`
	Last_name   string    `json:"last_name"`
	Email       string    `json:"email"`
	Created_on  time.Time `json:"created_on"`
	Updated_on  time.Time `json:"updated_on"`
	Status      int       `json:"status"`
	Language    string    `json:"language"`
}

func (ClientUser) TableName() string {
	return "users"
}

type UserProfile struct {
	Id         int    `json:"id"`
	User_id    int    `json:"testcue_id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Contact    int    `json:"email"`
	Photo      string `json:"photo"`
	Language   string `json:"language"`
	Theme      int    `json:"theme"`
}


func (UserProfile) TableName() string {
	return "user_profile"
}


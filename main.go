package main

import (
	"github.com/thejackrabbit/aqua"
	_"github.com/thejackrabbit/aero/conf"
	_"github.com/thejackrabbit/aero/db/orm"
	User "github.com/apitestcube/user/service"
)

func main() {
	server := aqua.NewRestServer()
	server.AddService(&User.UserService{})
	server.Run()
}
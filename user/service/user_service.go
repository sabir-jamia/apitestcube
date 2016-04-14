package UserService

import (
	"github.com/thejackrabbit/aqua"
	"github.com/apitestcube/user/model"
	"github.com/apitestcube/utils"
)

type UserService struct {
	aqua.RestService          `prefix:"api" root:"/" version:"1"`
	create aqua.PostApi       `url:"users"`
	delete aqua.DeleteApi     `url:"users/{id}"`
	deleteList aqua.DeleteApi `url:"users"`
	fetch aqua.GetApi         `url:"users/{id}"`
	fetchAll aqua.GetApi      `url:"users"`
	patch aqua.PatchApi       `url:"users/{id}"`
	update aqua.PutApi        `url:"users/{id}"`
	replaceList aqua.PutApi   `url:"users"`
}

type Response struct {
	Status   string       `json:"status"`
	Code     int          `json:"code"`
	Data     interface{}  `json:"data"`
	Message string       `json:"message"`
}

var UserTable = model.UserTable{}

/**
 * Create a resource
 *
 * @param aqua.Jar j
 * @return Response anonymous
*/
func (u *UserService) Create(j aqua.Jar) Response {
	j.LoadVars()
	data := model.User{}
	clientData := model.ClientUser{}
	userProfileData := model.UserProfile{}
	utils.GetPostData(j.Body, &data)
	utils.GetPostData(j.Body, &clientData)
	utils.GetPostData(j.Body, &userProfileData)
	err := UserTable.SaveUser(&data, &clientData, &userProfileData)
	if err == nil {
		response := map[string]interface{}{ "clientId" : data.ClientId}
		return Response{Status:"success", Code:200, Message:"", Data:response}
	} else {
		return Response{Status:"error", Code:405, Message:err.Error(), Data:""}
	}
}

/**
 * Delete a resource
 *
 * @param  int id
 * @return Response anonymous
*/
func (u *UserService) Delete(id int) Response {
	return Response{Status:"error", Code:405, Message:"The DELETE method has not been defined for individual resources", Data:""}
}

/**
 * Delete a collection, or members of a collection
 *
 * @param  aqua.Jar j
 * @return Response anonymous
*/
func(u *UserService) DeleteList(j aqua.Jar) Response {
	j.LoadVars()
	data := model.User{}
	utils.GetPostData(j.Body, &data)
	return Response{Status:"error", Code:405, Message:"The DELETE method has not been defined for collections", Data:""}
}

/**
 * Fetch a resource
 *
 * @param  int id
 * @return Response anonymous
*/
func (u *UserService) Fetch(id int) Response {
	user, err := UserTable.GetUserById(id)
	if err == nil {
		return Response{Status:"success", Code:200, Message:"", Data:user}
	}
	return Response{Status:"error", Code:405, Message:err.Error(), Data:""}
	//return Response{Status:"error", Code:405, Messaage:"The GET method has not been defined for individual resources", Data:""}
}

/**
 * Fetch all or a subset of resources
 *
 * @param  null
 * @return Response anonymous
*/
func (u *UserService) FetchAll(j aqua.Jar) Response{
	j.LoadVars()
	queryParams := j.QueryVars
	userName := queryParams["username"]
	email := queryParams["email"]
	if  userName != "" {
		row, err := UserTable.GetUserByUserName(userName)
		if err == nil {
			return Response{Status:"success", Code:200, Data:row}
		} else {
			return Response{Status:"error", Code:405, Data:"", Message:err.Error()}
		}
	}

	if email != "" {
		row, err := UserTable.GetUserByEmail(email)
		if err == nil {
			return Response{Status:"success", Code:200, Data:row}
		} else {
			return Response{Status:"error", Code:405, Data:"", Message:err.Error()}
		}
	}
	return Response{Status:"success", Code:200, Data:UserTable.ListUser()}
	//return Response{Status:"error", Code:405, Messaage:"The GET method has not been defined for collections", Data:""}
}

/**
 * Patch (partial in-place update) a resource
 *
 * @param  int id
 * @param  aqua.Jar j
 * @return Response anonymous
*/
func (UserService) Patch(id int, j aqua.Jar) Response{
	j.LoadVars()
	data := model.User{Id:id}
	utils.GetPostData(j.Body, &data)
	err := UserTable.UpdateStatus(id, &data)
	if err ==nil {
		return Response{Status:"status", Code:200, Message:"", Data:"status updated successfully"}
	}
	return Response{Status:"error", Code:405, Message:"status not updated successfully", Data:""}
	//return Response{Status:"error", Code:405, Messaage:"The PATCH method has not been defined for individual resources", Data:""}
}

/**
 * Update a resource
 *
 * @param  int id
 * @param aqua.Jar j
 * @return Response anonymous
*/
func (u *UserService) Update (id int, j aqua.Jar) Response {
	j.LoadVars()
	data := model.User{}
	utils.GetPostData(j.Body, &data)
	return Response{Status:"error", Code:405, Message:"The PUT method has not been defined for individual resources", Data:""}
}

/**
 * Replace a collection or members of a collection
 *
 * @param  aqua.Jar j
 * @return Response anonymous
*/
func (u *UserService)ReplaceList(j aqua.Jar) Response {
	data := model.User{}
	utils.GetPostData(j.Body, &data)
	return Response{Status:"error", Code:405, Message:"The PUT method has not been defined for collections", Data:""}
}
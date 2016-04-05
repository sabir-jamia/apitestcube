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
	Messaage string       `json:"message"`
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
		response := map[string]interface{}{ "clientId" : data.Client_id}
		return Response{Status:"success", Code:200, Messaage:"", Data:response}
	} else {
		return Response{Status:"error", Code:405, Messaage:"Some error occured", Data:""}
	}
}

/**
 * Delete a resource
 *
 * @param  int id
 * @return Response anonymous
*/
func (u *UserService) Delete(id int) Response {
	return Response{Status:"error", Code:405, Messaage:"The DELETE method has not been defined for individual resources", Data:""}
}

/**
 * Delete a collection, or members of a collection
 *
 * @param  aqua.Jar j
 * @return Response anonymous
*/
func(u *UserService) DeleteList(j aqua.Jar) Response {
	data := model.User{}
	utils.GetPostData(j.Body, &data)
	return Response{Status:"error", Code:405, Messaage:"The DELETE method has not been defined for collections", Data:""}
}

/**
 * Fetch a resource
 *
 * @param  int id
 * @return Response anonymous
*/
func (u *UserService) Fetch(id int) Response {
	return Response{Status:"error", Code:405, Messaage:"The GET method has not been defined for individual resources", Data:""}
}

/**
 * Fetch all or a subset of resources
 *
 * @param  null
 * @return Response anonymous
*/
func (user *UserService) FetchAll() Response{
	//add params
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
	data := model.User{}
	utils.GetPostData(j.Body, &data)
	return Response{Status:"error", Code:405, Messaage:"The PATCH method has not been defined for individual resources", Data:""}
}

/**
 * Update a resource
 *
 * @param  int id
 * @param aqua.Jar j
 * @return Response anonymous
*/
func (u *UserService) Update (id int, j aqua.Jar) Response {
	data := model.User{}
	utils.GetPostData(j.Body, &data)
	return Response{Status:"error", Code:405, Messaage:"The PUT method has not been defined for individual resources", Data:""}
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
	return Response{Status:"error", Code:405, Messaage:"The PUT method has not been defined for collections", Data:""}
}
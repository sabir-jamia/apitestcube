package model

import (
	"github.com/thejackrabbit/aero/db/orm"
	"time"
	"strings"
	"os/exec"
	"os"
	"strconv"
	"math/rand"
	"errors"
)

type UserTable struct {

}

func (u *UserTable) ListUser() []User{
	db := orm.Get(true)
	var users []User
	db.Find(&users)
	return users;
}

func (u *UserTable) SaveUser (user *User, clientUser *ClientUser, userProfile *UserProfile) error{
	db := orm.Get(true)
	trans := db.Begin()
	user.RegistrationDate = time.Now()
	if trans.NewRecord(user) {
		err := trans.Create(&user).Error
		if err == nil {
			user.ClientId = user.Id;
			trans.Save(&user)
			userProfile.UserId = user.Id
			userProfile.Language = "en"
			trans.NewRecord(&userProfile)
			err := trans.Create(&userProfile).Error
			curDirectory, _ := os.Getwd()
			curDirectory = "/var/www/ZendTestcube/module/User/src/User/Model/"
			if err == nil {
				clientDb := "clientdb0" + strconv.Itoa(user.ClientId)
				clientDbUser := "client0" + strconv.Itoa(user.ClientId)
				clientDbPass := "tolexo"
				trans.Exec("CREATE SCHEMA IF NOT EXISTS "+ clientDb +" DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;");
				trans.Exec("USE "+ clientDb +";");
				trans.Exec("CREATE USER '"+ clientDbUser + "'@'localhost' IDENTIFIED BY '"+ clientDbPass +"';");
				trans.Exec("GRANT SELECT, INSERT, UPDATE, DELETE ON "+ clientDb +".* TO '"+ clientDbUser +"'@'localhost';");
				dbUser := "root"
				dbPass := "tolexo"
				sqlConn := "mysql -u"+ dbUser +" -p"+ dbPass +" -hlocalhost "+ clientDb + " < " + curDirectory + "clientdb2.sql";
				exec.Command("sh", "-c", sqlConn).Run()

				clientUser.FirstName = strings.Title(strings.ToLower(clientUser.FirstName))
				clientUser.CreatedOn = time.Now()
				clientUser.UpdatedOn = time.Now()
				clientUser.TestcubeId = user.ClientId
				clientUser.Language = "en"
				if clientUser.LastName != "" {
					clientUser.LastName = strings.Title(strings.ToLower(clientUser.LastName))
				}
				trans.NewRecord(&clientUser)
				trans.Create(&clientUser)
				trans.Exec("USE testcubedb;");
				trans.Commit()
			} else {
				trans.Rollback()
				return err
			}
		} else {
			trans.Rollback()
			return err
		}
	} else {
		trans.Rollback()
		return nil
	}
	return nil
}

func (u *UserTable) GetUserById(id int) (*User, error){
	db := orm.Get(true)
	user := User{}
	err := db.Where(&User{Id : id, Status: true}).Find(&user).Error
	return &user, err
}

func (u *UserTable) GetUserByUserName(username string) (*User, error) {
	db := orm.Get(true)
	user := User{}
	err := db.Where(&User{Username:username}).Find(&user).Error
	return &user, err
}

func (u *UserTable) UpdateStatus(id int, user *User) error{
	db := orm.Get(true)
	err := db.Model(User{}).Where("id = ?", user.Id).Update("status", user.Status).Error
	return err
}

func (u *UserTable) GetNewPasswordByEmail(email string) (bool, error){
	db := orm.Get(true)
	characters := "abcdefghijklmnopqrstuvwxyz0123456789";
	length := len(characters);
	var str []byte;
	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++  {
		str = append(str, characters[rand.Intn(length - 1)]);
	}
	user := User{}
	err := db.Where("email = ?", email).First(&user).Error
	if err == nil {
		if user.Id == 0 {
			return false, errors.New("User does not exists")
		}
		user.Password = string(str)
		err := db.Save(&user).Error
		if err == nil {
			return true, nil
		}
	}
	return false, err
}

func (u *UserTable) IsEmailExist(email string)  error{
	db := orm.Get(true)
	user := User{}
	err := db.Where("email = ?", email).First(&user).Error
	if err == nil {
		if user.Id == 0 {
			return errors.New("User does not exist with given email")
		}
		return  nil;
	}
	return  err;
}

func (u *UserTable) CheckUserExists(username string) (err error){
	if strings.ContainsAny(username, "@") {
		err = u.IsEmailExist(username)
	} else {
		_, err = u.GetUserByUserName(username)
	}
	return err
}

func (u *UserTable) UserProfile(userId int) (map[string]interface{}, error) {
	var clientId, contact int
	var email, username, firstName, lastName, photo string
	var registrationDate time.Time
	db := orm.Get(true)
	sql := "SELECT u.client_id, u.email, u.username, u.registration_date, "+
	      "up.first_name, up.last_name, up.photo, up.contact " +
	      "FROM users u LEFT JOIN user_profile up "+
	      "ON u.client_id = up.user_id "+
	      "WHERE u.id = ?"
	row := db.Raw(sql, userId).Row()
	row.Scan(&clientId, &email, &username, &registrationDate, &firstName, &lastName, &photo, &contact)
	if clientId != 0 {
		return map[string]interface{}{
			"clientId" : clientId, "username" :username,
			"contact" : contact, "firstName" : firstName,
			"lastName" : lastName, "photo" : photo,
			"registrationDate" :registrationDate, "email" : email,
		}, nil
	}
	return map[string]interface{}{}, errors.New("User does not exist")
}

func (u *UserTable) ClientGeneralProfileUpdate(userProfile *UserProfile) error {
	db := orm.Get(true)
	err := db.Model(UserProfile{}).Where("user_id = ?", userProfile.UserId).
	Updates(UserProfile{FirstName:userProfile.FirstName, LastName:userProfile.LastName}).Error
	return err
}

func (u *UserTable) SuperGeneralProfileUpdate(userProfile *UserProfile) error {
	db := orm.Get(true)
	err := db.Model(UserProfile{}).Where("user_id", userProfile.UserId).
	Updates(UserProfile{FirstName:userProfile.FirstName,
		LastName:userProfile.LastName,
		Contact:userProfile.Contact,
		Photo:userProfile.Photo,
	}).Error
	return err;
}

func (u *UserTable) UpdateProfilePassword(passwordData map [string]string) {
	db := orm.Get(true)
	db.Model(User{}).Where("client_id = ?", passwordData["clientId"]).
	Updates(User{Password:passwordData["newPassword"]})
	//pass word should be stored in md5
}

func (u *UserTable) ClientGeneralProfileSettings(profile *UserProfile) {
	db := orm.Get(true)
	db.Model(UserProfile{}).Where("user_id = ?", profile.UserId).
	Updates(UserProfile{
		Language:profile.LastName,
		Theme:profile.Theme,
	})
}
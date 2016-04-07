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
	user.Registration_date = time.Now()
	if trans.NewRecord(user) {
		err := trans.Create(&user).Error
		if err == nil {
			user.Client_id = user.Id;
			trans.Save(&user)
			userProfile.User_id = user.Id
			userProfile.Language = "en"
			trans.NewRecord(&userProfile)
			err := trans.Create(&userProfile).Error
			curDirectory, _ := os.Getwd()
			curDirectory = "/var/www/ZendTestcube/module/User/src/User/Model/"
			if err == nil {
				clientDb := "clientdb0" + strconv.Itoa(user.Client_id)
				clientDbUser := "client0" + strconv.Itoa(user.Client_id)
				clientDbPass := "tolexo"
				trans.Exec("CREATE SCHEMA IF NOT EXISTS "+ clientDb +" DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;");
				trans.Exec("USE "+ clientDb +";");
				trans.Exec("CREATE USER '"+ clientDbUser + "'@'localhost' IDENTIFIED BY '"+ clientDbPass +"';");
				trans.Exec("GRANT SELECT, INSERT, UPDATE, DELETE ON "+ clientDb +".* TO '"+ clientDbUser +"'@'localhost';");
				dbUser := "root"
				dbPass := "tolexo"
				sqlConn := "mysql -u"+ dbUser +" -p"+ dbPass +" -hlocalhost "+ clientDb + " < " + curDirectory + "clientdb2.sql";
				exec.Command("sh", "-c", sqlConn).Run()

				clientUser.First_name = strings.Title(strings.ToLower(clientUser.First_name))
				clientUser.Created_on = time.Now()
				clientUser.Updated_on = time.Now()
				clientUser.Testcube_id = user.Client_id
				clientUser.Language = "en"
				if clientUser.Last_name != "" {
					clientUser.Last_name = strings.Title(strings.ToLower(clientUser.Last_name))
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
	err := db.Where(&User{Status:true, Username:username}).Find(&user).Error
	return &user, err
}

func (u *UserTable) UpdateStatus(id int, user *User) error{
	db := orm.Get(true)
	sql := "UPDATE users SET status = ? WHERE id = ?";
	err := db.Exec(sql, user.Status, user.Id).Error
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
		err = db.Table("users").Where("email = ?", email).Update("password", string(str)).Error
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

func (u *UserTable) ClientGeneralProfileUpdate(userProfile UserProfile) {
	db := orm.Get(true)
	update := db.Table("user_profile").Where("testcube_id = ?", userProfile.User_id)
	update.Updates(UserProfile{First_name:userProfile.First_name, Last_name:userProfile.Last_name})
}



/*class UserTable
{

protected $tableGateway;

protected $resultSetPrototype;

private $_dbConfig;

private $_adapter;

private $_connection;

public function __construct(TableGateway $tableGateway)
{
$this->tableGateway = $tableGateway;
$this->_adapter = $this->tableGateway->getAdapter();
$this->_connection = $this->_adapter->getDriver()->getConnection();
$this->resultSetPrototype = new ResultSet();
}

public function superGeneralProfileUpdate($profileData)
{
$data = array(
'profileFirstName' => $profileData['profileFirstName'],
'profileLastName' => $profileData['profileLastName'],
'profileContact' => $profileData['profileContact'],
'profilePic' => $profileData['profilePic'],
'random' => $profileData['random']
);

$sql = new Sql($this->tableGateway->getAdapter());

$update = $sql->update();
$update->table('user_profile');

if ($data['random'] == $data['profilePic']) {
$update->set(array(
'first_name' => $profileData['profileFirstName'],
'last_name' => $profileData['profileLastName'],
'contact' => $profileData['profileContact']
));
} else {
$update->set(array(
'first_name' => $profileData['profileFirstName'],
'last_name' => $profileData['profileLastName'],
'contact' => $profileData['profileContact'],
'photo' => $profileData['profilePic']
));
}

$update->where(array(
'user_id' => $profileData['clientId']
));
// echo $update->getSqlString();

$statement = $sql->prepareStatementForSqlObject($update);

try {
$result = $statement->execute(); // works fine
} catch (\Exception $e) {
die('Error: ' . $e->getMessage());
}

return $result;
}

public function updateProfilePassword($passwordData)
{
$data = array(
'clientId' => $passwordData['clientId'],
'newPassword' => $passwordData['newPassword']
);

$sql = new Sql($this->tableGateway->getAdapter());

$update = $sql->update();
$update->table('users');
$update->set(array(
'password' => md5($passwordData['newPassword'])
));
$update->where(array(
'client_id' => $passwordData['clientId']
));
// echo $update->getSqlString();

$statement = $sql->prepareStatementForSqlObject($update);
try {
$result = $statement->execute(); // works fine
} catch (\Exception $e) {
die('Error: ' . $e->getMessage());
}

return $result;
}

public function clientGeneralProfileSettings($profileData)
{
$data = array(
'clientId' => $profileData['clientId'],
'selectedLanguage' => $profileData['selectedLanguage'],
'themeColor' => $profileData['themeColor']
);

$sql = new Sql($this->tableGateway->getAdapter());

$update = $sql->update();
$update->table('user_profile');
$update->set(array(
'language' => $profileData['selectedLanguage']
));
$update->set(array(
'theme' => $profileData['themeColor']
));
$update->where(array(
'user_id' => $profileData['clientId']
));
// echo $update->getSqlString();

$statement = $sql->prepareStatementForSqlObject($update);
try {
$result = $statement->execute(); // works fine
$res = 1;
} catch (\Exception $e) {
die('Error: ' . $e->getMessage());
$res = 0;
}

return $res;
}

/*
 * public function clientGeneralProfileSettings($profileData){
 * $data = array(
 * 'clientId' => $profileData['clientId'],
 * 'selectedLanguage' => $profileData['selectedLanguage'],
 * 'themeColor' => $profileData['themeColor']
 * );
 * $theme = $profileData['themeColor'];
 * $language = $profileData['selectedLanguage'];
 * $clientId = $profileData['clientId'];
 *
 *
 *
 * $adapter1 = $this->dbAdapter12;
 * $themeQuery = "UPDATE (update user_profile SET theme = $theme, language = '$language' where user_id = $clientId;) as themeselected";
 * $themeData = $adapter1->query($themeQuery, Adapter::QUERY_MODE_EXECUTE)->toArray();;
 * echo $themeData;
 * die();
 *
 *
 *
 * return $themeData;
 *
 *
 *
 *
 *
 * /*$sql = new Sql ( $this->tableGateway->getAdapter() );
 *
 * $update = $sql->update();
 * $update->table('user_profile');
 * $update->set(array('language' => $profileData['selectedLanguage']));
 * $update->set(array('theme' => $profileData['themeColor']));
 * $update->where(array('client_id' => $profileData['clientId']));
 * //echo $update->getSqlString();
 *
 * $statement = $sql->prepareStatementForSqlObject($update);
 * try {
 * $result = $statement->execute(); // works fine
 * } catch (\Exception $e) {
 * die('Error: ' . $e->getMessage());
 * }
 *
 * return $result;
 *
 * }
}*/
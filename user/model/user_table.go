package model

import (
	"github.com/thejackrabbit/aero/db/orm"
	"time"
	"strings"
	"os/exec"
	"os"
	"strconv"
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
	//'password' => md5($register->password),
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

public function setDbCredentails($dbConfig)
{
$this->_dbConfig = $dbConfig;
}

public function updatestatus($id)
{
$data = array(
'status' => '1'
);
if ($this->getUserById($id)) {
$this->tableGateway->update($data, array(
'id' => $id
));
return 1;
} else {
return 0;
}
}

public function getUserById($id)
{
$rowset = $this->tableGateway->select(array(
'id' => $id,
'status' => '0'
));
$row = $rowset->current();

if (! $row) {
/* throw new \Exception("Could not found row $username");
}
return $row;
}

public function getUserByUserName($username, $usernameType)
{
$rowset = $this->tableGateway->select(array(
$usernameType => $username,
'status' => '1'
));
$row = $rowset->current();

if (! $row) {
return 'notconfirmed';
// throw new \Exception("Could not found row $username");
}
return $row;
}

public function getPasswordByEmail($email)
{
$characters = 'abcdefghijklmnopqrstuvwxyz0123456789';
$length = strlen($characters);
$string = '';
for ($i = 0; $i < $length - 30; $i ++) {
$string .= $characters[rand(0, $length - 1)];
}

$sql = new Sql($this->tableGateway->getAdapter());

$update = $sql->update();
$update->table('users');
$update->set(array(
'password' => md5($string)
));
$update->where(array(
'email' => $email
));

// echo $update->getSqlString();

$statement = $sql->prepareStatementForSqlObject($update);
try {
$result = $statement->execute(); // works fine
} catch (\Exception $e) {
die('Error: ' . $e->getMessage());
}
// die();
return $string;
}

public function isEmailexist($txtVal)
{
$rowset = $this->tableGateway->select(array(
'email' => $txtVal
));
$row = $rowset->current();

if ($row) {
return true;
} else {
return false;
}
}

public function checkUserExists($user)
{
$pos = strpos($user, '@');
if ($pos === false) {
$rowset = $this->tableGateway->select(array(
'username' => $user
));
$row = $rowset->current();
} else {
$rowset = $this->tableGateway->select(array(
'email' => $user
));
$row = $rowset->current();
}

if ($row) {
return true;
} else {
return false;
}
}

public function userProfile($userid)
{
$user_id = (int) $userid;

$where = new Where();
$where->equalTo('usersTabl.id', $user_id);

$sql = new Sql($this->tableGateway->getAdapter());
$select = $sql->select()
->from(array(
'usersTabl' => 'users'
))
->columns(array(
'clientId' => 'client_id',
'email' => 'email',
'username' => 'username',
'regDate' => 'registration_date'
))
->join(array(
'userProfileTabl' => 'user_profile'
), 'usersTabl.client_id = userProfileTabl.user_id', array(
'firstName' => 'first_name',
'lastName' => 'last_name',
'photo' => 'photo',
'contact' => 'contact'
))
->where($where);

// echo $select->getSqlString();
$statement = $sql->prepareStatementForSqlObject($select);
$result = $this->resultSetPrototype->initialize($statement->execute())
->toArray();
// \Zend\Debug\Debug::dump($result);
return $result;
}

public function clientGeneralProfileUpdate($profileData)
{
$data = array(
'profileFirstName' => $profileData['profileFirstName'],
'profileLastName' => $profileData['profileLastName']
);

$sql = new Sql($this->tableGateway->getAdapter());

$update = $sql->update();
$update->table('users');
$update->set(array(
'first_name' => $profileData['profileFirstName'],
'last_name' => $profileData['profileLastName']
));
$update->where(array(
'testcube_id' => $profileData['clientId']
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

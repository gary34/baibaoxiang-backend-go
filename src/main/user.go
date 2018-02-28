package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"utils"
)

const (
	AnonymousUser = "anonymous"
	NormalUser    = "normal"
	AdminUser     = "admin"
)

type User struct {
	CommonModle
	UID      string `gorm:"column:uid"`
	Name     string
	Role     string
	Scopes   string
	Password string   `json:"-"`
	Salt     string   `json:"-"`
	Baobeis  []Baobei `gorm:"many2many:user_baobeis;"`
}

func createUser(uid, name, role, scopes, password string) (*User, error) {

	salt := ""
	pass := ""
	if password != "" {
		salt = newSalt()
		pass = encodingPassword(password, salt)
	}

	user := &User{
		UID:      uid,
		Name:     name,
		Role:     role,
		Scopes:   scopes,
		Password: pass,
		Salt:     salt,
	}
	conn := GetDB()
	if err := conn.Create(user).Error; err != nil {
		gLogger.ELog("create user error %s", err.Error())
		return nil, err
	}
	return user, nil
}

func createAnonymousUser(uid string) (*User, error) {
	return createUser(
		uid,
		uid,
		AnonymousUser,
		"",
		"")
}

func newSalt() string {
	return utils.RandomString(6)
}

func encodingPassword(password, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(password))
	expectedMAC := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(expectedMAC)
}

func setUserPasswordHandler(context *gin.Context) {
	// context.Params.
	uid := context.PostForm("uid")
	// gLogger.ILog("uid is %s", uid)
	password := context.PostForm("password")
	user, err := getUserByUID(uid)
	err = user.setPassword(password)
	// err = rmUserBaobei(uid, uint(iBaobei))
	renderJSON(nil, err, context)
}
func (user *User) checkPassword(password string) bool {
	if user.Password == "" {
		return true
	}
	return encodingPassword(password, user.Salt) == user.Password
}

func (user *User) setPassword(password string) error {
	user.Salt = newSalt()
	user.Password = encodingPassword(password, user.Salt)
	if err := GetDB().Save(user).Error; err != nil {
		gLogger.ELog("update  user password error %s", err.Error())
		return err
	}
	return nil
}

func getUserFromContext(context *gin.Context) (*User, error) {
	c, ok := context.Get("user_context")
	var err error

	if !ok {
		err = fmt.Errorf("user need login")
		renderJSON(nil, err, context)
		return nil, err
	}
	u, ok := c.(*User)
	if !ok {
		err = fmt.Errorf("invalied user data")
		renderJSON(nil, err, context)
		return nil, err
	}
	return u, nil
}
func getUserFavoriesHandler(context *gin.Context) {
	user, err := getUserFromContext(context)
	if err != nil {
		context.Abort()
		return
	}
	list, err := user.getFavories()
	renderJSON(list, err, context)
}

func userAuthenticateHandler(context *gin.Context) {
	username := context.Query("username")
	password := context.Query("password")
	user, err := getUserByName(username)
	if err != nil {
		renderJSON(nil, err, context)
		return
	}
	if user == nil {
		renderJSON(nil, fmt.Errorf("user not exist"), context)
		return
	}
	if !user.checkPassword(password) {
		renderJSON(nil, fmt.Errorf("user passsword worrry"), context)
		return
	}
	ss, err := user.genJWT()
	renderJSON(struct{ JWT string }{ss}, err, context)
}

type UserClaims struct {
	Name   string `json:"name,omitempty"`
	UID    string `json:"uid,omitempty"`
	Scopes string `json:"scopes,omitempty"`
	Role   string `json:"role,omitempty"`
	jwt.StandardClaims
}

const jwtSignKey = `key`

func (user *User) genJWT() (string, error) {

	claims := &UserClaims{
		Name:   user.Name,
		UID:    user.UID,
		Scopes: user.Scopes,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	ss, err := utils.GenJwt(claims, jwtSignKey)
	if err != nil {
		gLogger.ELog("signed jwt token error %s", err.Error())
	}
	return ss, err
}

func (user *User) getFavories() (list []Baobei, err error) {
	conn := GetDB()
	if err = conn.Model(user).Related(&user.Baobeis, "Baobeis").Error; err != nil {
		gLogger.ELog("query user baobeis error %s", err.Error)
		return
	}
	for i := 0; i < len(user.Baobeis); i++ {
		b := &user.Baobeis[i]
		b.InFavor = true
	}
	return user.Baobeis, err
}

func getUserByUID(uid string) (user *User, err error) {
	conn := GetDB()
	// var u User
	user = &User{UID: uid}
	if err = conn.Where("uid = ?", uid).First(user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			gLogger.ELog("query user error %s,uid %s", err.Error(), uid)

		} else {
			//create a anonymous
			return createAnonymousUser(uid)
			// err = nil
		}
		return
	}
	return
}

func getUserByName(name string) (user *User, err error) {
	conn := GetDB()
	// var u User
	user = &User{Name: name}
	if err = conn.Where("name = ?", name).First(user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			gLogger.ELog("query user error %s,uid %s", err.Error(), name)
		} else {
			//create a anonymous
			// return createAnonymousUser(uid)
			err = nil
		}
		return
	}
	return
}

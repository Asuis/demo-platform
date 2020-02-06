package user

import (
	"crypto/md5"
	"demo-plaform/model/db"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Login struct {
	Account string `form:"account" json:"account" binding:"required"`
	Passwd string `form:"passwd" json:"passwd" xml:"passwd"  binding:"required"`
}

type LoginResult struct {
	Token string
}

type Register struct {
	Account string `form:"account" json:"account" binding:"required"`
	Type string
	Passwd string `form:"passwd" json:"passwd" xml:"passwd"  binding:"required"`
	IP string
}

func SignIn(login *Login) (string, error) {

	var usr = &db.User{
		LoginName:     login.Account,
		LoginType: db.LoginNotype,
	}

	has,err := db.GetByUser(usr)

	if err != nil {
		return "", err
	}

	if !has {
		return "", errors.New("user is not exist")
	}

	hasher := md5.New()

	salt := usr.Salt
	hasher.Write([]byte(fmt.Sprintf("%s_%d", login.Passwd, salt)))
	if hex.EncodeToString(hasher.Sum(nil)) == usr.Passwd {
		token,err := SignedToken(usr)
		return token ,err
	}
	return "", err
}

func SignUp(register *Register) (string, error) {

	// 是否存在
	has, err := db.GetByUser(&db.User{
		LoginName:      register.Account,
		LoginType: db.LoginNotype,
	})
	if err != nil {
		return "", err
	}
	if has {
		return "", errors.New("the user is already exist")
	}

	now := time.Now()
	hasher := md5.New()
	salt := rand.Uint64()
	hasher.Write([]byte(fmt.Sprintf("%s_%d", register.Passwd, salt)))

	usr := db.User{
		Name:            register.Account,
		Passwd:          hex.EncodeToString(hasher.Sum(nil)),
		Email:           register.Account,
		LoginType:       db.LoginNotype,
		LoginSource:     0,
		LoginName:       register.Account,
		Type:            db.UserTypeIndividual,
		OwnedOrgs:       nil,
		Orgs:            nil,
		Repos:           nil,
		Roles:           nil,
		DockerContainer: nil,
		DockerImage:     nil,
		Location:        "",
		Website:         "",
		Rands:           "",
		Salt:            salt,
		Created:         now,
		CreatedUnix:     now.Unix(),
		Updated:         now,
		UpdatedUnix:     now.Unix(),
	}

	
	err = db.CreateUser(&usr)
	if err != nil {
		return "", err
	}

	token, err := SignedToken(&usr)

	if err != nil {
		return "", err
	}

	return token, nil
}


func updateAvatar(userID int64, avatar string) error {
	err := db.UpdateUser(&db.User{
		Id: userID,
		Avatar: avatar,
	})
	if err != nil {
		return err
	}
	return nil
}

func update()  {}

func exit(token string) error {

	sign, err := ParseToken(token)
	if err != nil {
		return err
	}
	err = db.UpdateUser(&db.User{
		Id: sign.Ac,
		Salt: rand.Uint64(),
	})
	if err != nil {
		return err
	}
	return nil
}

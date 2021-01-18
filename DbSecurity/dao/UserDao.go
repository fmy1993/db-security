package dao

import (
	"log"

	"DbSecurity/model"
	"DbSecurity/tool"
)

type UserDao struct {
	*tool.Orm
}

func NewUserDao() *UserDao {
	return &UserDao{tool.DbEngine}
}

//根据用户号查询用户
func (ud *UserDao) QueryUserByPhone(phone string) *model.User {
	var user model.User
	if _, err := ud.Where("phone = ?", phone).Get(&user); err != nil {
		log.Fatal(err.Error())
	}
	return &user
}

//插入用户
func (ud *UserDao) InsertUser(user model.User) int64 {
	result, err := ud.InsertOne(&user)
	if err != nil {
		log.Fatal(err.Error())
		return 0
	}
	return result
}

//检查指纹是否已经存在
func (ud *UserDao) IsFingerPrintExist(fingerPrint string) *model.User {
	var user model.User
	_, err := ud.Where("finger_print = ?", fingerPrint).Get(&user)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &user
}

//修改密码
func (ud *UserDao) RevisePwd(id int64, user *model.User) error {
	_, err := ud.Exec("update user set password = ? where id = ?", user.Password, id)
	if err != nil {
		return err
	}
	return nil
}

//获取所有用户
func (ud *UserDao) QueryAllUser() (*[]model.User, error) {
	var users []model.User
	err := ud.Find(&users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}
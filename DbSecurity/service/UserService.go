package service

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	"DbSecurity/dao"
	"DbSecurity/model"
	"DbSecurity/param"
	"DbSecurity/tool"
)

type UserService struct {
}

//获取盐值
func (us *UserService) Salt() string {
	chars := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"
	bytes := []byte(chars)
	var salt []byte
	for i := 0; i < 6; i++ {
		salt = append(salt, bytes[rand.Intn(len(bytes))])
	}
	return string(salt)
}

//注册用户
func (us *UserService) Register(registerParam param.RegisterParam) (*model.User, int64) {
	userDao := dao.NewUserDao()
	user := userDao.QueryUserByPhone(registerParam.Phone)
	if user.Id != 0 {
		return user, user.Id
	} else {
		user := model.User{}
		var us UserService
		salt := us.Salt()
		user.Phone = registerParam.Phone
		user.Password = salt + ":" + tool.EncoderSha256(salt+registerParam.Password)
		user.DateJoined = time.Now().Unix()
		user.FingerPrint = us.getFingerPrint(registerParam.Phone)
		user.IsSuperUser = 0
		user.Id = userDao.InsertUser(user)
		return &user, 0
	}
}

//生成指纹
func (us *UserService) getFingerPrint(phone string) string {
	fingerPrint := ""
	for {
		fingerPrint = ""
		var lsBit []int
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < 100; i++ {
			lsBit = append(lsBit, rand.Intn(2))
		}
		q, _ := strconv.Atoi(string(phone[2]))
		if q == 0 {
			q += 1
		}
		var lsNum []int
		add := 0
		for j := 0; j < len(phone); j++ {
			temp, _ := strconv.Atoi(string(phone[j]))
			add += temp
		}
		for k := 0; k < 20; k++ {
			lsNum = append(lsNum, int(math.Pow(float64(k), float64(q)))%add)
		}
		for p := 0; p < len(lsNum); p++ {
			fingerPrint += strconv.Itoa(lsBit[lsNum[p]])
		}
		userDao := dao.NewUserDao()
		if userDao.IsFingerPrintExist(fingerPrint).Id == 0 {
			break
		}
	}
	return fingerPrint
}

//根据用户名查询用户
func (us *UserService) GetUser(phone string) *model.User {
	userDao := dao.NewUserDao()
	var user *model.User
	user = userDao.QueryUserByPhone(phone)
	return user
}

//修改密码
func (us *UserService) RevisePassword(id int64, user *model.User) error {
	userDao := dao.NewUserDao()
	err := userDao.RevisePwd(id, user)
	if err != nil {
		return err
	}
	return nil
}

//得到user表的所有数据
func (us *UserService) GetUsers() (*[]model.User, error) {
	var userDao = dao.NewUserDao()
	return userDao.QueryAllUser()
}


package service

import (
	"db-security-backend/param/userParam"
	"db-security-backend/util"
	"math"
	"math/rand"
	"strconv"
	"time"

	"db-security-backend/dao"
	"db-security-backend/model"
)

type UserService struct {
}

// Salt 获取盐值
func (us *UserService) Salt() string {
	chars := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"
	bytes := []byte(chars)
	var salt []byte
	for i := 0; i < 6; i++ {
		salt = append(salt, bytes[rand.Intn(len(bytes))])
	}
	return string(salt)
}

// Register 注册用户
func (us *UserService) Register(registerParam userParam.RegisterParam) (*model.User, int64) {
	userDao := dao.NewUserDao()
	user := userDao.QueryUserByPhone(registerParam.Phone)
	if user.Id != 0 {
		return user, user.Id
	} else {
		user := model.User{}
		salt := us.Salt()
		user.Phone = registerParam.Phone
		user.Password = salt + ":" + util.EncoderSha256(salt+registerParam.Password)
		user.DateJoined = time.Now().Format("2006-01-02 15:04:05")
		user.FingerPrint = us.getFingerPrint(registerParam.Phone)
		user.IsSuperUser = 0
		user.Id = userDao.InsertUser(user)
		return &user, 0
	}
}

// 生成指纹
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

// GetUser 根据手机号查询用户
func (us *UserService) GetUser(phone string) *model.User {
	userDao := dao.NewUserDao()
	var user *model.User
	user = userDao.QueryUserByPhone(phone)
	return user
}

// GetUsersByFingerprint 根据fingerPrint查询相关用户
func (us *UserService) GetUsersByFingerprint(fingerprint string) (*model.User, int) {
	var userDao = dao.NewUserDao()
	users, _ := userDao.QueryAllUser()
	var min = 20
	var user model.User
	for i := range *users {
		var count = 0
		for j := 0; j < 20; j++ {
			if (*users)[i].FingerPrint[j] != fingerprint[j] {
				count++
			}
		}
		if count < min {
			min = count
			user = (*users)[i]
		}
	}
	return &user, min
}

// GetUsersByFingerprint2 根据fingerPrint查询相关用户
func (us *UserService) GetUsersByFingerprint2(fingerprint string) map[string]int {
	var userDao = dao.NewUserDao()
	users, _ := userDao.QueryAllUser()
	var res = make(map[string]int)
	for i := range *users {
		var count = 0
		for j := 0; j < 20; j++ {
			if (*users)[i].FingerPrint[j] != fingerprint[j] {
				count++
			}
		}
		res[(*users)[i].Phone] = count
	}
	return res
}

// RevisePassword 修改密码
func (us *UserService) RevisePassword(id int64, user *model.User) error {
	userDao := dao.NewUserDao()
	err := userDao.RevisePwd(id, user)
	if err != nil {
		return err
	}
	return nil
}

// GetUsers 得到user表的所有数据
func (us *UserService) GetUsers() (*[]model.User, error) {
	var userDao = dao.NewUserDao()
	return userDao.QueryAllUser()
}

package pwd

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// HashPassword 密码哈希
func HashPassword(password string) string {
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
	}
	return string(pwdBytes)
}

// CheckPasswordHash 验证密码哈希
func CheckPasswordHash(DbPwd string, UsrPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(DbPwd), []byte(UsrPwd)) //验证
	if err != nil {
		return false
	}
	return true
}

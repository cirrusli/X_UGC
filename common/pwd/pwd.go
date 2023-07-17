package pwd

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Encryption 密码加密
func Encryption(password string) string {
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
	}
	return string(pwdBytes)
}

// Decryption 解密
func Decryption(password1 string, password2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2)) //验证（对比） password1数据库中的 password2用户输入的
	if err != nil {
		return false
	}
	return true
}

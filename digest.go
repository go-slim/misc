package misc

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func MD5(str string) string {
	hash := md5.Sum([]byte(str))
	//md5str := hex.EncodeToString(hash)
	md5str := fmt.Sprintf("%x", hash) //将[]byte转成16进制
	return strings.ToUpper(md5str)
}

func Sha1(str string) string {
	hash := sha1.Sum([]byte(str))
	//return hex.EncodeToString(hash)
	return fmt.Sprintf("%x", hash)
}

func Sha256(str string) string {
	hash := sha256.Sum256([]byte(str))
	//return hex.EncodeToString(hash)
	return fmt.Sprintf("%x", hash)
}

func PasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

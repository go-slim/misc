package misc

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// MD5 计算字符串的 MD5 摘要，并以大写十六进制字符串返回。
func MD5(str string) string {
	hash := md5.Sum([]byte(str))
	//md5str := hex.EncodeToString(hash)
	md5str := fmt.Sprintf("%x", hash) //将[]byte转成16进制
	return strings.ToUpper(md5str)
}

// Sha1 计算字符串的 SHA-1 摘要，并以十六进制字符串返回（小写）。
func Sha1(str string) string {
	hash := sha1.Sum([]byte(str))
	//return hex.EncodeToString(hash)
	return fmt.Sprintf("%x", hash)
}

// Sha256 计算字符串的 SHA-256 摘要，并以十六进制字符串返回（小写）。
func Sha256(str string) string {
	hash := sha256.Sum256([]byte(str))
	//return hex.EncodeToString(hash)
	return fmt.Sprintf("%x", hash)
}

// PasswordHash 使用 bcrypt 生成口令哈希。
// 返回的哈希可用于后续的 PasswordVerify 校验。
func PasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// PasswordVerify 使用 bcrypt 对明文口令与哈希进行校验。
// 返回 true 表示匹配，false 表示不匹配。
func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

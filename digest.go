package misc

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"math"
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

var base62Digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Encode62 将10进制转换成62进制
//
// Encode62(15758306521) // HCSFOD
func Encode62(num int64) string {
	var bytes []byte
	for num > 0 {
		bytes = append(bytes, base62Digits[num%62])
		num = num / 62
	}
	for left, right := 0, len(bytes)-1; left < right; left, right = left+1, right-1 {
		bytes[left], bytes[right] = bytes[right], bytes[left]
	}
	return string(bytes)
}

// Decode62 将62进制转换成10进制
// Decode62("HCSFOD") // 15758306521
func Decode62(str string) int64 {
	var num int64
	n := len(str)
	for i := 0; i < n; i++ {
		pos := strings.IndexByte(base62Digits, str[i])
		num += int64(math.Pow(62, float64(n-i-1)) * float64(pos))
	}
	return num
}

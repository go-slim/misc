package misc

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// MD5 calculates the MD5 digest of a string and returns it as an uppercase hexadecimal string.
func MD5(str string) string {
	hash := md5.Sum([]byte(str))
	//md5str := hex.EncodeToString(hash)
	md5str := fmt.Sprintf("%x", hash) // convert []byte to hexadecimal
	return strings.ToUpper(md5str)
}

// Sha1 calculates the SHA-1 digest of a string and returns it as a hexadecimal string (lowercase).
func Sha1(str string) string {
	hash := sha1.Sum([]byte(str))
	//return hex.EncodeToString(hash)
	return fmt.Sprintf("%x", hash)
}

// Sha256 calculates the SHA-256 digest of a string and returns it as a hexadecimal string (lowercase).
func Sha256(str string) string {
	hash := sha256.Sum256([]byte(str))
	//return hex.EncodeToString(hash)
	return fmt.Sprintf("%x", hash)
}

// PasswordHash generates password hash using bcrypt.
// The returned hash can be used for subsequent PasswordVerify validation.
func PasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// PasswordVerify uses bcrypt to validate plaintext password against hash.
// Returns true for match, false for mismatch.
func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

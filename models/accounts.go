package models

import (
	"github.com/dgrijalva/jwt-go"
	u "blackbox-kaizen/utils"
	"strings"
	"github.com/jinzhu/gorm"
	"os"
	"golang.org/x/crypto/bcrypt"
)

// Token JWT claims struct
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// Account struct represents the user account
type Account struct {
	gorm.Model
	Email string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token";sql:"-"`
}

// Validate incoming user details
func (account *Account) Validate() (map[string] interface{}, bool) {
	
	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is not valid"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is not valid"), false
	}

	// Email must be unique
	temp := &Account{}

	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}

	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user"), false
	}

	return u.Message(false, "Requirement passed"), true
}

// Create user account
func (account *Account) Create() (map[string] interface{}) {
	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString
	
	account.Password = "" // delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}


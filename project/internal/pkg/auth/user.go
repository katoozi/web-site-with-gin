package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var userSQLQuery = `
	CREATE TABLE IF NOT EXISTS "user" (
		"id" serial not null PRIMARY KEY,
		"first_name" varchar(30),
		"last_name" varchar(150),
		"password" varchar(130) not null,
		"last_login" timestamptz default now(),
		"date_joined" timestamptz default now(),
		"username" varchar(150) unique not null,
		"email" varchar(254),
		"is_active" boolean not null default 'true',
		"is_staff" boolean not null default 'false',
		"is_superuser" boolean not null default 'false'
	);`

// User will have the user table schema
type User struct {
	ID         int       `db:"id" sqltools:"id"`
	FirstName  string    `db:"first_name" sqltools:"first_name"`
	LastName   string    `db:"last_name" sqltools:"last_name"`
	Password   string    `db:"password" sqltools:"password"`
	LastLogin  time.Time `db:"last_login" sqltools:"last_login"`
	DateJoined time.Time `db:"date_joined" sqltools:"date_joined"`
	Username   string    `db:"username" sqltools:"username"`
	IsActive   bool      `db:"is_active" sqltools:"is_active" default:"1"`
	IsStaff    bool      `db:"is_staff" sqltools:"is_staff" default:"1"`
	IsSperuser bool      `db:"is_superuser" sqltools:"is_superuser" default:"1"`
	Email      string    `db:"email" sqltools:"email"`
}

// NewUser is the User type factory function
func NewUser(firstName, lastName, username, email, password string) *User {
	hashpassword, err := generate(password)
	if err != nil {
		log.Fatalf("Error while encrypting password: %v", err)
	}
	return &User{
		FirstName:  firstName,
		LastName:   lastName,
		Username:   username,
		Email:      email,
		Password:   hashpassword,
		IsActive:   true,
		IsStaff:    false,
		IsSperuser: false,
	}
}

// GenerateInsertQuery will generate sql insert query for postgresql
func (u *User) GenerateInsertQuery() string {
	schema := `INSERT INTO "user" (first_name,last_name,email,username,password,is_active,is_staff,is_superuser) VALUES ('%s','%s','%s','%s','%s','%t','%t','%t');`
	return fmt.Sprintf(
		schema,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Username,
		u.Password,
		u.IsActive,
		u.IsStaff,
		u.IsSperuser,
	)
}

func generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

//Compare string with generated hash
func (u *User) Compare(s string) error {
	incoming := []byte(s)
	return bcrypt.CompareHashAndPassword([]byte(u.Password), incoming)
}

// GetUser will fetch user from db by username
func GetUser(username string, dbCon *sqlx.DB) *User {
	query := fmt.Sprintf(`SELECT * from "user" WHERE username='%s';`, username)
	userObj := User{}
	err := dbCon.Get(&userObj, query)
	if err != nil {
		// log.Printf("Error while get user from db:%v", err)
		return nil
	}
	return &userObj
}

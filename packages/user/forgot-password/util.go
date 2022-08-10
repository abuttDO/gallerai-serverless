package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var repo Repository

func makeResponse(statusCode int, body interface{}, err error) *Response {
	bodyString, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("error marshalling body: %s\n", err)
		return &Response{}
	}
	var response Response
	response.StatusCode = statusCode
	response.Headers = map[string]string{
		"Content-Type": "application/json",
	}
	response.Body = string(bodyString)
	return &response
}

func initDatabase() (db *gorm.DB) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=require", os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_PASSWORD"))
	if os.Getenv("DATABASE_TLS") == "false" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_PASSWORD"))
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func validateToken(tokenString string) (user string, valid bool, err error) {
	// Convert os.Getenv("HMAC_SECRET") to []byte
	secret := []byte(os.Getenv("HMAC_SECRET"))

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
		// Cast interface to string
		user = claims["user"].(string)
		return user, true, nil
	}
	return "", false, err
}

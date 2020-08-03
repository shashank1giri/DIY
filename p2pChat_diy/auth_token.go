// This file consists of utility functions for JWT token authentication

package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)
const (
	secret     = "q80a7uni=#_uo&ezdd76+^rto=i_n%2$@^rk4(@3kekuu=@912"
	// Todo :- can use git secrets for pushing the password after encrypting
	UserId     = "user_id"
	ExpiryTime = "exp_time"
)

func getToken(user_id int, duration time.Duration)(string,error){
	claims:=jwt.MapClaims{}
	claims[UserId] = user_id
	claims[ExpiryTime] = time.Now().Add(duration).Unix()
	log.Println(user_id,claims["exp"])
	token , err:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte("key"))
	return token,err
}

/* verifies whether a jwt token string is valid or not.
	exp_time > current time
	and was signed with the same secret key. */
func verifyToken(tokenStr string) (bool,int) {
	userId := -1
	token,err:=jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _,ok:=token.Method.(*jwt.SigningMethodHMAC);!ok{
			return nil,fmt.Errorf("invalid algo used in token: %v",token.Header["alg"])
		}
		return []byte(secret),nil
	})
	if err != nil{
		logrus.WithFields(logrus.Fields{
			"error":err,
			"token":tokenStr,
		}).Error("Error while parsing jwt token")

		return false,userId
	}
	if claims,ok:=token.Claims.(jwt.MapClaims);ok && token.Valid{
		expiryTime,ok := claims[ExpiryTime].(float64)
		user,ok1 := claims[UserId].(float64)
		log.Print(time.Now().Unix())
		if !ok || !ok1 || int64(expiryTime) < time.Now().Unix(){
			return false,userId
		}
		return true,int(user)

	}else{
		return false,userId
	}

}
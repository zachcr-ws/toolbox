package auth

import (
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"strconv"
	"time"
)

type AuthClient struct {
	Username string
	Password string
	ExpDay   int
}

func New(username, password string, day int) *AuthClient {
	return &AuthClient{
		Username: username,
		Password: password,
		ExpDay:   day,
	}
}

func (this *AuthClient) Key(l int) []byte {
	stMap := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(stMap)
	resul := []byte{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		resul = append(resul, bytes[r.Intn(len(bytes))])
	}
	return resul
}

func (this *AuthClient) Verify(token string, key []byte) error {
	tokenObj, err := jwt.Parse(token, func(tokenObj *jwt.Token) (interface{}, error) {
		if _, ok := tokenObj.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tokenObj.Header["alg"])
		}

		return key, nil
	})

	return err
}

func (this *AuthClient) Produce() (token string, key string, err error) {

	jwtObj := jwt.New(jwt.SigningMethodHS256)
	jwtObj.Claims["exp"] = time.Now().Add(time.Hour * 24 * this.ExpDay)
	jwtObj.Claims["aeris"] = strconv.Itoa(int(time.Now().UnixNano())) + this.Username

	key = this.Key(16)
	token, err = jwtObj.SignedString(key)
	return
}

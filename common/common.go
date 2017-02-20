package common

import (
	"io"
	"io/ioutil"
	"math/rand"
    "time"
    "strconv"
)

func CopyBody(body io.ReadCloser) []byte {
	if body == nil {
		return []byte{}
	}

	safe := &io.LimitedReader{R: body, N: 1 << 26} // 64M max memory
	rb, _ := ioutil.ReadAll(safe)
	body.Close()

	return rb
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func RandomNumber(count int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result string = ""
	for i := 0; i < count; i++ {
		result += strconv.Itoa(r.Intn(10))
	}
	return result
}

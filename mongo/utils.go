package mongo

import (
	"crypto/md5"
	"encoding/hex"
	"gopkg.in/mgo.v2"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var Seesion map[string]*mgo.Session

func MgoGeneratorId(r string) string {
	ctx := md5.New()
	st := strconv.FormatInt(time.Now().UnixNano(), 16)
	return hex.EncodeToString(ctx.Sum([]byte(st + r)))
}

func MgoGeneratorRandom(leng int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < leng; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetSession(hosts string, user string, password string, db string) (*mgo.Session, error) {
	if Seesion == nil {
		Seesion = make(map[string]*mgo.Session)
	} else {
		s := Seesion[db]
		if s != nil {
			if s.Ping() == nil {
				return s.Clone(), nil
			}
		}
	}

	var err error
	dial := new(mgo.DialInfo)
	dial.Addrs = strings.Split(hosts, ",")
	dial.Username = user
	dial.Password = password
	dial.Database = db
	ss, err := mgo.DialWithInfo(dial)
	if err != nil {
		return nil, err
	}
	Seesion[db] = ss
	return ss.Clone(), nil
}

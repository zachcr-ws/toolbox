package mongo

import (
	"crypto/md5"
	"encoding/hex"
	"gopkg.in/mgo.v2"
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

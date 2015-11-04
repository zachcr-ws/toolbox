package mysql

import (
	"database/sql"
	"github.com/astaxie/beedb"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var MasterDB *sql.DB
var SlaveDb *sql.DB
var Config MysqlConfig

func InitMysqlConfig(master, slave, user, password, name string) {
	Config.MasterAddress = master
	Config.SlaveAddress = slave
	Config.User = user
	Config.Password = password
	Config.DbName = name

	MasterDB = ConnectMysql(true)
	SlaveDb = ConnectMysql(false)
}

type MysqlConfig struct {
	MasterAddress string
	SlaveAddress  string
	User          string
	Password      string
	DbName        string
}

type MysqlQuery struct {
	Table   string
	OrderBy string
	Offset  int
	Size    int
	Fields  string
	GroupBy string
	Where   string
}

func ConnectMysql(master bool) *sql.DB {
	addr := Config.MasterAddress
	user := Config.User
	name := Config.DbName
	pswd := Config.Password
	if !master {
		addr = Config.SlaveAddress
	}

	db, err := sql.Open("mysql", user+":"+pswd+"@tcp("+addr+")/"+name+"?charset=utf8")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func (q *MysqlQuery) Exec(newOrm bool, query string, args ...interface{}) (sql.Result, error) {
	db := SlaveDb
	if newOrm {
		db = ConnectMysql(false)
	}
	orm := beedb.New(db)
	return orm.Exec(query, args...)
}

func (q *MysqlQuery) FindOne(result interface{}, newOrm bool) error {
	db := SlaveDb
	if newOrm {
		db = ConnectMysql(false)
	}

	orm := beedb.New(db)
	return orm.SetTable(q.Table).Where(q.Where).OrderBy(q.OrderBy).Limit(q.Offset, q.Size).Select(q.Fields).Find(result)
}

func (q *MysqlQuery) FindAll(result interface{}, newOrm bool) error {
	db := SlaveDb
	if newOrm {
		db = ConnectMysql(false)
	}

	orm := beedb.New(db)
	return orm.SetTable(q.Table).Where(q.Where).OrderBy(q.OrderBy).Limit(q.Offset, q.Size).Select(q.Fields).FindAll(result)
}

func (q *MysqlQuery) Upsert(data interface{}, newOrm bool) error {
	db := MasterDB
	if newOrm {
		db = ConnectMysql(true)
	}

	orm := beedb.New(db)
	return orm.SetTable(q.Table).Save(data)
}

func (q *MysqlQuery) SaveBatch(data []interface{}) error {
	orm := beedb.New(MasterDB)

	for i := 0; i < len(data); i++ {
		err := orm.SetTable(q.Table).Save(data[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (q *MysqlQuery) Delete(data interface{}, newOrm bool) (int64, error) {
	db := MasterDB
	if newOrm {
		db = ConnectMysql(true)
	}
	orm := beedb.New(db)

	err := orm.SetTable(q.Table).Where(q.Where).OrderBy(q.OrderBy).Limit(q.Offset, q.Size).Find(data)
	if err != nil {
		return int64(0), err
	}
	return orm.SetTable(q.Table).Delete(data)
}

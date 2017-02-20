package mysql_v2

import (
	"database/sql"
	"github.com/astaxie/beedb"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func InitMysqlConfig(master, slave, user, password, name string, idle, open int) (*sql.DB, *sql.DB) {

	config := &MysqlConfig{
		MasterAddress: master,
		SlaveAddress:  slave,
		User:          user,
		Password:      password,
		DbName:        name,
		PoolIdle:      idle,
		PoolOpen:      open,
	}

	MasterDB := ConnectMysql(config, true)
	SlaveDb := ConnectMysql(config, false)

	return MasterDB, SlaveDb
}

type MysqlConfig struct {
	MasterAddress string
	SlaveAddress  string
	User          string
	Password      string
	DbName        string
	PoolIdle      int
	PoolOpen      int
}

type MysqlQuery struct {
	MasterDB *sql.DB
	SlaveDb  *sql.DB
	Table    string
	OrderBy  string
	Offset   int
	Size     int
	Fields   string
	GroupBy  string
	Where    string
}

func ConnectMysql(config *MysqlConfig, master bool) *sql.DB {
	addr := config.MasterAddress
	user := config.User
	name := config.DbName
	pswd := config.Password
	if !master {
		addr = config.SlaveAddress
	}

	db, err := sql.Open("mysql", user+":"+pswd+"@tcp("+addr+")/"+name+"?charset=utf8")
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxIdleConns(config.PoolIdle)
	db.SetMaxOpenConns(config.PoolOpen)
	return db
}

func (q *MysqlQuery) Exec(query string, args ...interface{}) (sql.Result, error) {
	db := q.SlaveDb
	orm := beedb.New(db)
	return orm.Exec(query, args...)
}

func (q *MysqlQuery) FindOne(result interface{}) error {
	db := q.SlaveDb
	orm := beedb.New(db)
	if q.Fields == "" {
		q.Fields = "*"
	}
	return orm.SetTable(q.Table).Where(q.Where).OrderBy(q.OrderBy).Limit(q.Size, q.Offset).Select(q.Fields).Find(result)
}

func (q *MysqlQuery) FindAll(result interface{}) error {
	db := q.SlaveDb
	orm := beedb.New(db)
	if q.Fields == "" {
		q.Fields = "*"
	}
	return orm.SetTable(q.Table).Where(q.Where).OrderBy(q.OrderBy).Limit(q.Size, q.Offset).Select(q.Fields).FindAll(result)
}

func (q *MysqlQuery) Upsert(data interface{}) error {
	db := q.MasterDB
	orm := beedb.New(db)
	return orm.SetTable(q.Table).Save(data)
}

func (q *MysqlQuery) Delete() (int64, error) {
	db := q.MasterDB
	orm := beedb.New(db)
	return orm.SetTable(q.Table).Where(q.Where).DeleteRow()
}

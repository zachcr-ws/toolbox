package mysql

import (
	"database/sql"
	"github.com/astaxie/beedb"
	"github.com/astaxie/beego/orm"
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
	dbtype := "master"
	if !master {
		addr = Config.SlaveAddress
		dbtype = "slave"
	}
	conn := user + ":" + pswd + "@tcp(" + addr + ")/" + name + "?charset=utf8"
	err := orm.RegisterDataBase(dbtype, "mysql", conn)
	if err != nil {
		log.Fatalln(err)
	}
	db, err := orm.GetDB(dbtype)
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxIdleConns(1000)
	db.SetMaxOpenConns(2000)
	return db
}

func SetMasterConns(idle, open int) {
	MasterDB.SetMaxIdleConns(idle)
	MasterDB.SetMaxOpenConns(open)
}

func SetSlaveConns(idle, open int) {
	SlaveDb.SetMaxIdleConns(idle)
	SlaveDb.SetMaxOpenConns(open)
}

func (q *MysqlQuery) Exec(newOrm bool, query string, args ...interface{}) (sql.Result, error) {
	db := SlaveDb
	if newOrm {
		db = ConnectMysql(false)
		defer db.Close()
	}
	dbmodel := beedb.New(db)
	return dbmodel.Exec(query, args...)
}

func (q *MysqlQuery) FindOne(result interface{}, newOrm bool) error {
	db := SlaveDb
	if newOrm {
		db = ConnectMysql(false)
		defer db.Close()
	}

	dbmodel := beedb.New(db)
	if q.Fields == "" {
		q.Fields = "*"
	}
	return dbmodel.SetTable(q.Table).Where(q.Where).OrderBy(q.OrderBy).Limit(q.Size, q.Offset).Select(q.Fields).Find(result)
}

func (q *MysqlQuery) FindAll(result interface{}, newOrm bool) error {
	db := SlaveDb
	if newOrm {
		db = ConnectMysql(false)
		defer db.Close()
	}

	dbmodel := beedb.New(db)
	if q.Fields == "" {
		q.Fields = "*"
	}
	return dbmodel.SetTable(q.Table).Where(q.Where).OrderBy(q.OrderBy).Limit(q.Size, q.Offset).Select(q.Fields).FindAll(result)
}

func (q *MysqlQuery) Upsert(data interface{}, newOrm bool) error {
	db := MasterDB
	if newOrm {
		db = ConnectMysql(true)
		defer db.Close()
	}

	dbmodel := beedb.New(db)
	return dbmodel.SetTable(q.Table).Save(data)
}

func (q *MysqlQuery) Delete(newOrm bool) (int64, error) {
	db := MasterDB
	if newOrm {
		db = ConnectMysql(true)
		defer db.Close()
	}
	dbmodel := beedb.New(db)
	return dbmodel.SetTable(q.Table).Where(q.Where).DeleteRow()
}

// @Title Tarn
// @Description exec sql by transaction
// @Param    sql1, sql2, sql3...
// @Success true, nil
// @Failure false, error
func (q *MysqlQuery) Tarn(sql ...string) (bool, error) {
	tran := orm.NewOrm()
	tran.Begin()
	var (
		Err  error
		Flag bool
	)
	for _, v := range sql {
		_, err := tran.Raw(v).Exec()
		if err != nil {
			Flag = false
			Err = err
			goto RESULT
		}
	}
	tran.Commit()
	Flag = true
	return Flag, Err
RESULT:
	tran.Rollback()
	return Flag, Err

}

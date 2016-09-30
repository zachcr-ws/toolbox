package mysql_v3

import (
	"database/sql"
	"github.com/astaxie/beedb"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var Config MysqlConfig

func InitMysqlConfig(master, slave, user, password, name string, idle, open int) {
	Config.MasterAddress = master
	Config.SlaveAddress = slave
	Config.User = user
	Config.Password = password
	Config.DbName = name

	ConnectMysql(true)
	ConnectMysql(false)
}

type MysqlConfig struct {
	MasterAddress string
	SlaveAddress  string
	User          string
	Password      string
	DbName        string
	IdleConns     int
	OpenConns     int
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

func ConnectMysql(master bool) {
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
	orm.SetMaxIdleConns(dbtype, Config.IdleConns)
	orm.SetMaxOpenConns(dbtype, Config.OpenConns)
}

func (q *MysqlQuery) Exec(newOrm bool, query string, args ...interface{}) (s sql.Result, err error) {
	db, err := orm.GetDB("master")
	if err != nil {
		return
	}
	defer db.Close()

	dbmodel := beedb.New(db)
	return dbmodel.Exec(query, args...)
}

func (q *MysqlQuery) FindOne(result interface{}) error {
	db, err := orm.GetDB("slave")
	if err != nil {
		return err
	}
	defer db.Close()

	dbmodel := beedb.New(db)
	if q.Fields == "" {
		q.Fields = "*"
	}
	return dbmodel.SetTable(q.Table).Where(q.Where).OrderBy(q.OrderBy).Limit(q.Size, q.Offset).Select(q.Fields).Find(result)
}

func (q *MysqlQuery) FindAll(result interface{}) error {
	db, err := orm.GetDB("slave")
	if err != nil {
		return err
	}
	defer db.Close()

	dbmodel := beedb.New(db)
	if q.Fields == "" {
		q.Fields = "*"
	}
	return dbmodel.SetTable(q.Table).Where(q.Where).OrderBy(q.OrderBy).Limit(q.Size, q.Offset).Select(q.Fields).FindAll(result)
}

func (q *MysqlQuery) Upsert(data interface{}) error {
	db, err := orm.GetDB("master")
	if err != nil {
		return err
	}
	defer db.Close()

	dbmodel := beedb.New(db)
	return dbmodel.SetTable(q.Table).Save(data)
}

func (q *MysqlQuery) Delete() (int64, error) {
	db, err := orm.GetDB("master")
	if err != nil {
		return 0, err
	}
	defer db.Close()

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

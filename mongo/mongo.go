package mongo

import (
	"gopkg.in/mgo.v2"
)

type MgoClient struct {
	Addr   string
	User   string
	Passwd string
	Dbname string
}

func (m *MgoClient) Session() (*mgo.Session, error) {
	return GetSession(
		m.Addr,
		m.User,
		m.Passwd,
		m.Dbname,
	)
}

func (m *MgoClient) Count(c string, query map[string]interface{}) (count int, err error) {
	session, err := m.Session()
	if err != nil {
		return
	}

	s := session.Copy()
	defer s.Close()
	count, err = s.DB(m.Dbname).C(c).Find(query).Count()
	return
}

func (m *MgoClient) FindOne(c string, query map[string]interface{}, response interface{}) error {
	session, err := m.Session()
	if err != nil {
		return err
	}

	s := session.Copy()
	defer s.Close()
	return s.DB(m.Dbname).C(c).Find(query).One(response)
}

func (m *MgoClient) FindAll(c string, query map[string]interface{}, response interface{}) error {
	session, err := m.Session()
	if err != nil {
		return err
	}

	s := session.Copy()
	defer s.Close()
	return s.DB(m.Dbname).C(c).Find(query).All(response)
}

func (m *MgoClient) Find(c string, query map[string]interface{}, skip, limit int, response interface{}, sorts ...string) error {
	session, err := m.Session()
	if err != nil {
		return err
	}

	s := session.Copy()
	defer s.Close()
	return s.DB(m.Dbname).C(c).Find(query).Skip(skip).Limit(limit).Sort(sorts...).All(response)
}

func (m *MgoClient) Upsert(c string, query map[string]interface{}, response interface{}) error {
	session, err := m.Session()
	if err != nil {
		return err
	}

	s := session.Copy()
	defer s.Close()
	_, err = s.DB(m.Dbname).C(c).Upsert(query, response)
	return err
}

func (m *MgoClient) Remove(c string, query map[string]interface{}) error {
	session, err := m.Session()
	if err != nil {
		return err
	}

	s := session.Copy()
	defer s.Close()
	return s.DB(m.Dbname).C(c).Remove(query)
}

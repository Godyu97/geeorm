package gee

import (
	"database/sql"
	"fmt"
	"github.com/Godyu97/geeorm/gee/dialect"
	"github.com/Godyu97/geeorm/gee/session"
	"github.com/Godyu97/geeorm/log"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}
	// make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{
		db:      db,
		dialect: dial,
	}
	log.Info("Connect database success")
	return e, nil
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

type TxFunc func(session2 *session.Session) (any, error)

func (e *Engine) Transaction(f TxFunc) (result any, err error) {
	s := e.NewSession()
	err = s.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback()
		} else {
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()
			err = s.Commit()
		}
	}()
	return f(s)
}

// difference returns s1 - s2
func SliceSub[T comparable](s1 []T, s2 []T) []T {
	diff := make([]T, 0)
	tmp := make(map[T]struct{}, 0)
	for _, v2 := range s2 {
		if _, ok := tmp[v2]; ok == false {
			tmp[v2] = struct{}{}
		}
	}
	for _, v1 := range s1 {
		if _, ok := tmp[v1]; ok == false {
			diff = append(diff, v1)
		}
	}
	return diff
}

// Migrate table
func (e *Engine) Migrate(value interface{}) error {
	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := SliceSub(table.FieldNames, columns)
		delCols := SliceSub(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}

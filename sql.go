package helpers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/jinzhu/gorm"
)

type SqlQuery struct {
	fileName, query string
	db              *gorm.DB
	params          []interface{}

	rows    *sql.Rows
	total   int
	columns []string
	length  int
}

func (q *SqlQuery) init() (err error) {
	if q.query == "" {
		if q.query, err = (&File{Name: q.fileName}).Content(); err != nil {
			return err
		}
	}

	if q.rows, err = q.db.Raw(q.query, q.params...).Rows(); err != nil {
		return err
	}

	if q.columns, err = q.rows.Columns(); err != nil {
		return err
	}

	q.length = len(q.columns)

	return nil
}

func (q *SqlQuery) scanRowToMap() (map[string]interface{}, error) {
	current, value := q.makeResultReceiver(), map[string]interface{}{}

	if err := q.rows.Scan(current...); err != nil {
		return value, err
	}

	for i := 0; i < q.length; i++ {
		k := q.columns[i]
		val := *(current[i]).(*interface{})
		value[k] = val
	}

	return value, nil
}

func (q *SqlQuery) makeResultReceiver() []interface{} {
	result := make([]interface{}, 0, q.length)
	for i := 0; i < q.length; i++ {
		var current interface{} = struct{}{}
		result = append(result, &current)
	}
	return result
}

func (q *SqlQuery) Fetch(obj interface{}) (next bool, err error) {
	if next = q.rows.Next(); next {
		switch v := obj.(type) {
		case *map[string]interface{}:
			if *v, err = q.scanRowToMap(); err != nil {
				q.rows.Close()
				return false, err
			}
		default:
			if err = q.db.ScanRows(q.rows, obj); err != nil {
				q.rows.Close()
				return false, err
			}
		}
		q.total++
		return next, nil
	} else {
		return next, q.rows.Close()
	}
}

func (q *SqlQuery) FetchAll(obj interface{}) ([]interface{}, error) {
	var (
		next bool
		err  error
		list = []interface{}{}
	)

	for next, err = q.Fetch(obj); err == nil && next; next, err = q.Fetch(obj) {
		list = append(list, reflect.ValueOf(obj).Elem().Interface())
	}

	return list, err
}

func (q *SqlQuery) ToList() ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}
	if rows, err := q.FetchAll(&map[string]interface{}{}); err != nil {
		return result, err
	} else {
		for _, item := range rows {
			result = append(result, item.(map[string]interface{}))
		}
		return result, nil
	}
}

func (q *SqlQuery) First(obj interface{}) error {
	if next, err := q.Fetch(obj); err != nil {
		return err
	} else if !next {
		return errors.New("record not found")
	}
	q.rows.Close()
	return nil
}

func (q *SqlQuery) Write(w http.ResponseWriter, start, end string, obj interface{}) (err error) {

	if start == "" {
		start = "["
	}
	if end == "" {
		end = "]"
	}

	w.Write([]byte(start))

	var next bool
	for next, err = q.Fetch(obj); err == nil && next; next, err = q.Fetch(obj) {
		if out, err := json.Marshal(obj); err != nil {
			return err
		} else {
			if q.total > 1 {
				w.Write([]byte(","))
			}
			w.Write(out)
		}
	}
	w.Write([]byte(end))
	return err
}

func (q *SqlQuery) Total() int {
	return q.total
}

func NewSqlFromFile(fileName string, db *gorm.DB, params ...interface{}) (*SqlQuery, error) {
	q := &SqlQuery{fileName: fileName, db: db, params: params}
	return q, q.init()
}

func NewSql(query string, db *gorm.DB, params ...interface{}) (*SqlQuery, error) {
	q := &SqlQuery{query: query, db: db, params: params}
	return q, q.init()
}

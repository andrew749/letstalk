// Code generated by ModelQ
// authentication_data.go contains model for the database table [letstalk.authentication_data]

package data

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/mijia/modelq/gmq"
	"strings"
)

type AuthenticationData struct {
	UserId       int    `json:"user_id"`
	PasswordHash string `json:"password_hash"`
}

// Start of the AuthenticationData APIs.

func (obj AuthenticationData) String() string {
	if data, err := json.Marshal(obj); err != nil {
		return fmt.Sprintf("<AuthenticationData UserId=%v>", obj.UserId)
	} else {
		return string(data)
	}
}

func (obj AuthenticationData) Get(dbtx gmq.DbTx) (AuthenticationData, error) {
	filter := AuthenticationDataObjs.FilterUserId("=", obj.UserId)
	if result, err := AuthenticationDataObjs.Select().Where(filter).One(dbtx); err != nil {
		return obj, err
	} else {
		return result, nil
	}
}

func (obj AuthenticationData) Insert(dbtx gmq.DbTx) (AuthenticationData, error) {
	_, err := AuthenticationDataObjs.Insert(obj).Run(dbtx)
	return obj, err
}

func (obj AuthenticationData) Update(dbtx gmq.DbTx) (int64, error) {
	fields := []string{"PasswordHash"}
	filter := AuthenticationDataObjs.FilterUserId("=", obj.UserId)
	if result, err := AuthenticationDataObjs.Update(obj, fields...).Where(filter).Run(dbtx); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (obj AuthenticationData) Delete(dbtx gmq.DbTx) (int64, error) {
	filter := AuthenticationDataObjs.FilterUserId("=", obj.UserId)
	if result, err := AuthenticationDataObjs.Delete().Where(filter).Run(dbtx); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

// Start of the inner Query Api

type _AuthenticationDataQuery struct {
	gmq.Query
}

func (q _AuthenticationDataQuery) Where(f gmq.Filter) _AuthenticationDataQuery {
	q.Query = q.Query.Where(f)
	return q
}

func (q _AuthenticationDataQuery) OrderBy(by ...string) _AuthenticationDataQuery {
	tBy := make([]string, 0, len(by))
	for _, b := range by {
		sortDir := ""
		if b[0] == '-' || b[0] == '+' {
			sortDir = string(b[0])
			b = b[1:]
		}
		if col, ok := AuthenticationDataObjs.fcMap[b]; ok {
			tBy = append(tBy, sortDir+col)
		}
	}
	q.Query = q.Query.OrderBy(tBy...)
	return q
}

func (q _AuthenticationDataQuery) GroupBy(by ...string) _AuthenticationDataQuery {
	tBy := make([]string, 0, len(by))
	for _, b := range by {
		if col, ok := AuthenticationDataObjs.fcMap[b]; ok {
			tBy = append(tBy, col)
		}
	}
	q.Query = q.Query.GroupBy(tBy...)
	return q
}

func (q _AuthenticationDataQuery) Limit(offsets ...int64) _AuthenticationDataQuery {
	q.Query = q.Query.Limit(offsets...)
	return q
}

func (q _AuthenticationDataQuery) Page(number, size int) _AuthenticationDataQuery {
	q.Query = q.Query.Page(number, size)
	return q
}

func (q _AuthenticationDataQuery) Run(dbtx gmq.DbTx) (sql.Result, error) {
	return q.Query.Exec(dbtx)
}

type AuthenticationDataRowVisitor func(obj AuthenticationData) bool

func (q _AuthenticationDataQuery) Iterate(dbtx gmq.DbTx, functor AuthenticationDataRowVisitor) error {
	return q.Query.SelectList(dbtx, func(columns []gmq.Column, rb []sql.RawBytes) bool {
		obj := AuthenticationDataObjs.toAuthenticationData(columns, rb)
		return functor(obj)
	})
}

func (q _AuthenticationDataQuery) One(dbtx gmq.DbTx) (AuthenticationData, error) {
	var obj AuthenticationData
	err := q.Query.SelectOne(dbtx, func(columns []gmq.Column, rb []sql.RawBytes) bool {
		obj = AuthenticationDataObjs.toAuthenticationData(columns, rb)
		return true
	})
	return obj, err
}

func (q _AuthenticationDataQuery) List(dbtx gmq.DbTx) ([]AuthenticationData, error) {
	result := make([]AuthenticationData, 0, 10)
	err := q.Query.SelectList(dbtx, func(columns []gmq.Column, rb []sql.RawBytes) bool {
		obj := AuthenticationDataObjs.toAuthenticationData(columns, rb)
		result = append(result, obj)
		return true
	})
	return result, err
}

func (q _AuthenticationDataQuery) Count(dbtx gmq.DbTx) (int, error) {
	result := 0

	err := q.Query.SelectCount(dbtx, func(columns []gmq.Column, rb []sql.RawBytes) bool {
		if len(columns) == len(rb) {
			for i := range columns {
				if "_count" == columns[i].Name {
					result = gmq.AsInt(rb[i])

					return true
				}
			}
		}

		return true
	})

	return result, err
}

// Start of the model facade Apis.

type _AuthenticationDataObjs struct {
	fcMap map[string]string
}

func (o _AuthenticationDataObjs) Names() (schema, tbl, alias string) {
	return "letstalk", "authentication_data", "AuthenticationData"
}

func (o _AuthenticationDataObjs) Select(fields ...string) _AuthenticationDataQuery {
	q := _AuthenticationDataQuery{}
	if len(fields) == 0 {
		fields = []string{"UserId", "PasswordHash"}
	}
	q.Query = gmq.Select(o, o.columns(fields...))
	return q
}

func (o _AuthenticationDataObjs) Insert(obj AuthenticationData) _AuthenticationDataQuery {
	q := _AuthenticationDataQuery{}
	q.Query = gmq.Insert(o, o.columnsWithData(obj, "UserId", "PasswordHash"))
	return q
}

func (o _AuthenticationDataObjs) Update(obj AuthenticationData, fields ...string) _AuthenticationDataQuery {
	q := _AuthenticationDataQuery{}
	q.Query = gmq.Update(o, o.columnsWithData(obj, fields...))
	return q
}

func (o _AuthenticationDataObjs) Delete() _AuthenticationDataQuery {
	q := _AuthenticationDataQuery{}
	q.Query = gmq.Delete(o)
	return q
}

///// Managed Objects Filters definition

func (o _AuthenticationDataObjs) FilterUserId(op string, p int, ps ...int) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("user_id", op, params...)
}

func (o _AuthenticationDataObjs) FilterPasswordHash(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("password_hash", op, params...)
}

///// Managed Objects Columns definition

func (o _AuthenticationDataObjs) ColumnUserId(p ...int) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"user_id", value}
}

func (o _AuthenticationDataObjs) ColumnPasswordHash(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"password_hash", value}
}

////// Internal helper funcs

func (o _AuthenticationDataObjs) newFilter(name, op string, params ...interface{}) gmq.Filter {
	if strings.ToUpper(op) == "IN" {
		return gmq.InFilter(name, params)
	}
	return gmq.UnitFilter(name, op, params[0])
}

func (o _AuthenticationDataObjs) toAuthenticationData(columns []gmq.Column, rb []sql.RawBytes) AuthenticationData {
	obj := AuthenticationData{}
	if len(columns) == len(rb) {
		for i := range columns {
			switch columns[i].Name {
			case "user_id":
				obj.UserId = gmq.AsInt(rb[i])
			case "password_hash":
				obj.PasswordHash = gmq.AsString(rb[i])
			}
		}
	}
	return obj
}

func (o _AuthenticationDataObjs) columns(fields ...string) []gmq.Column {
	data := make([]gmq.Column, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "UserId":
			data = append(data, o.ColumnUserId())
		case "PasswordHash":
			data = append(data, o.ColumnPasswordHash())
		}
	}
	return data
}

func (o _AuthenticationDataObjs) columnsWithData(obj AuthenticationData, fields ...string) []gmq.Column {
	data := make([]gmq.Column, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "UserId":
			data = append(data, o.ColumnUserId(obj.UserId))
		case "PasswordHash":
			data = append(data, o.ColumnPasswordHash(obj.PasswordHash))
		}
	}
	return data
}

var AuthenticationDataObjs _AuthenticationDataObjs

func init() {
	AuthenticationDataObjs.fcMap = map[string]string{
		"UserId":       "user_id",
		"PasswordHash": "password_hash",
	}
	gob.Register(AuthenticationData{})
}

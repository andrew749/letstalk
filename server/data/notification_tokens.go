// Code generated by ModelQ
// notification_tokens.go contains model for the database table [letstalk.notification_tokens]

package data

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/mijia/modelq/gmq"
	"strings"
)

type NotificationTokens struct {
	Id     string `json:"id"`
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

// Start of the NotificationTokens APIs.

func (obj NotificationTokens) String() string {
	if data, err := json.Marshal(obj); err != nil {
		return fmt.Sprintf("<NotificationTokens Id=%v>", obj.Id)
	} else {
		return string(data)
	}
}

func (obj NotificationTokens) Get(dbtx gmq.DbTx) (NotificationTokens, error) {
	filter := NotificationTokensObjs.FilterId("=", obj.Id)
	if result, err := NotificationTokensObjs.Select().Where(filter).One(dbtx); err != nil {
		return obj, err
	} else {
		return result, nil
	}
}

func (obj NotificationTokens) Insert(dbtx gmq.DbTx) (NotificationTokens, error) {
	_, err := NotificationTokensObjs.Insert(obj).Run(dbtx)
	return obj, err
}

func (obj NotificationTokens) Update(dbtx gmq.DbTx) (int64, error) {
	fields := []string{"UserId", "Token"}
	filter := NotificationTokensObjs.FilterId("=", obj.Id)
	if result, err := NotificationTokensObjs.Update(obj, fields...).Where(filter).Run(dbtx); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (obj NotificationTokens) Delete(dbtx gmq.DbTx) (int64, error) {
	filter := NotificationTokensObjs.FilterId("=", obj.Id)
	if result, err := NotificationTokensObjs.Delete().Where(filter).Run(dbtx); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

// Start of the inner Query Api

type _NotificationTokensQuery struct {
	gmq.Query
}

func (q _NotificationTokensQuery) Where(f gmq.Filter) _NotificationTokensQuery {
	q.Query = q.Query.Where(f)
	return q
}

func (q _NotificationTokensQuery) OrderBy(by ...string) _NotificationTokensQuery {
	tBy := make([]string, 0, len(by))
	for _, b := range by {
		sortDir := ""
		if b[0] == '-' || b[0] == '+' {
			sortDir = string(b[0])
			b = b[1:]
		}
		if col, ok := NotificationTokensObjs.fcMap[b]; ok {
			tBy = append(tBy, sortDir+col)
		}
	}
	q.Query = q.Query.OrderBy(tBy...)
	return q
}

func (q _NotificationTokensQuery) GroupBy(by ...string) _NotificationTokensQuery {
	tBy := make([]string, 0, len(by))
	for _, b := range by {
		if col, ok := NotificationTokensObjs.fcMap[b]; ok {
			tBy = append(tBy, col)
		}
	}
	q.Query = q.Query.GroupBy(tBy...)
	return q
}

func (q _NotificationTokensQuery) Limit(offsets ...int64) _NotificationTokensQuery {
	q.Query = q.Query.Limit(offsets...)
	return q
}

func (q _NotificationTokensQuery) Page(number, size int) _NotificationTokensQuery {
	q.Query = q.Query.Page(number, size)
	return q
}

func (q _NotificationTokensQuery) Run(dbtx gmq.DbTx) (sql.Result, error) {
	return q.Query.Exec(dbtx)
}

type NotificationTokensRowVisitor func(obj NotificationTokens) bool

func (q _NotificationTokensQuery) Iterate(dbtx gmq.DbTx, functor NotificationTokensRowVisitor) error {
	return q.Query.SelectList(dbtx, func(columns []gmq.Column, rb []sql.RawBytes) bool {
		obj := NotificationTokensObjs.toNotificationTokens(columns, rb)
		return functor(obj)
	})
}

func (q _NotificationTokensQuery) One(dbtx gmq.DbTx) (NotificationTokens, error) {
	var obj NotificationTokens
	err := q.Query.SelectOne(dbtx, func(columns []gmq.Column, rb []sql.RawBytes) bool {
		obj = NotificationTokensObjs.toNotificationTokens(columns, rb)
		return true
	})
	return obj, err
}

func (q _NotificationTokensQuery) List(dbtx gmq.DbTx) ([]NotificationTokens, error) {
	result := make([]NotificationTokens, 0, 10)
	err := q.Query.SelectList(dbtx, func(columns []gmq.Column, rb []sql.RawBytes) bool {
		obj := NotificationTokensObjs.toNotificationTokens(columns, rb)
		result = append(result, obj)
		return true
	})
	return result, err
}

func (q _NotificationTokensQuery) Count(dbtx gmq.DbTx) (int, error) {
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

type _NotificationTokensObjs struct {
	fcMap map[string]string
}

func (o _NotificationTokensObjs) Names() (schema, tbl, alias string) {
	return "letstalk", "notification_tokens", "NotificationTokens"
}

func (o _NotificationTokensObjs) Select(fields ...string) _NotificationTokensQuery {
	q := _NotificationTokensQuery{}
	if len(fields) == 0 {
		fields = []string{"Id", "UserId", "Token"}
	}
	q.Query = gmq.Select(o, o.columns(fields...))
	return q
}

func (o _NotificationTokensObjs) Insert(obj NotificationTokens) _NotificationTokensQuery {
	q := _NotificationTokensQuery{}
	q.Query = gmq.Insert(o, o.columnsWithData(obj, "Id", "UserId", "Token"))
	return q
}

func (o _NotificationTokensObjs) Update(obj NotificationTokens, fields ...string) _NotificationTokensQuery {
	q := _NotificationTokensQuery{}
	q.Query = gmq.Update(o, o.columnsWithData(obj, fields...))
	return q
}

func (o _NotificationTokensObjs) Delete() _NotificationTokensQuery {
	q := _NotificationTokensQuery{}
	q.Query = gmq.Delete(o)
	return q
}

///// Managed Objects Filters definition

func (o _NotificationTokensObjs) FilterId(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("id", op, params...)
}

func (o _NotificationTokensObjs) FilterUserId(op string, p int, ps ...int) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("user_id", op, params...)
}

func (o _NotificationTokensObjs) FilterToken(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("token", op, params...)
}

///// Managed Objects Columns definition

func (o _NotificationTokensObjs) ColumnId(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"id", value}
}

func (o _NotificationTokensObjs) ColumnUserId(p ...int) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"user_id", value}
}

func (o _NotificationTokensObjs) ColumnToken(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"token", value}
}

////// Internal helper funcs

func (o _NotificationTokensObjs) newFilter(name, op string, params ...interface{}) gmq.Filter {
	if strings.ToUpper(op) == "IN" {
		return gmq.InFilter(name, params)
	}
	return gmq.UnitFilter(name, op, params[0])
}

func (o _NotificationTokensObjs) toNotificationTokens(columns []gmq.Column, rb []sql.RawBytes) NotificationTokens {
	obj := NotificationTokens{}
	if len(columns) == len(rb) {
		for i := range columns {
			switch columns[i].Name {
			case "id":
				obj.Id = gmq.AsString(rb[i])
			case "user_id":
				obj.UserId = gmq.AsInt(rb[i])
			case "token":
				obj.Token = gmq.AsString(rb[i])
			}
		}
	}
	return obj
}

func (o _NotificationTokensObjs) columns(fields ...string) []gmq.Column {
	data := make([]gmq.Column, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "Id":
			data = append(data, o.ColumnId())
		case "UserId":
			data = append(data, o.ColumnUserId())
		case "Token":
			data = append(data, o.ColumnToken())
		}
	}
	return data
}

func (o _NotificationTokensObjs) columnsWithData(obj NotificationTokens, fields ...string) []gmq.Column {
	data := make([]gmq.Column, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "Id":
			data = append(data, o.ColumnId(obj.Id))
		case "UserId":
			data = append(data, o.ColumnUserId(obj.UserId))
		case "Token":
			data = append(data, o.ColumnToken(obj.Token))
		}
	}
	return data
}

var NotificationTokensObjs _NotificationTokensObjs

func init() {
	NotificationTokensObjs.fcMap = map[string]string{
		"Id":     "id",
		"UserId": "user_id",
		"Token":  "token",
	}
	gob.Register(NotificationTokens{})
}

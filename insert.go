package sqb

import "bytes"

func Insert(table string) *InsertQuery {
	q := InsertQuery{
		Query: Query{
			tableName: table,
		},
		queryChanges: &queryChanges{
			changes:         make(map[string]interface{}),
			changeFields:    make([]string, 0),
			rawChanges:      make(map[string]string),
			rawChangeFields: make([]string, 0),
		},
	}
	return &q
}

type InsertQuery struct {
	Query
	*queryChanges
}

func (q *InsertQuery) Set(key string, value interface{}) *InsertQuery {
	q.queryChanges.Set(key, value)
	return q
}

func (q *InsertQuery) SetRaw(key string, value string) *InsertQuery {
	q.queryChanges.SetRaw(key, value)
	return q
}

func (q *InsertQuery) SetData(data map[string]interface{}) *InsertQuery {
	q.queryChanges.SetData(data)
	return q
}

func (q *InsertQuery) SetRawData(data map[string]string) *InsertQuery {
	q.queryChanges.SetRawData(data)
	return q
}

func (q *InsertQuery) Build() (string, []interface{}) {
	bs := bytes.NewBufferString("INSERT INTO `")
	bs.WriteString(q.tableName)
	bs.WriteString("` (")
	bs2 := bytes.NewBufferString(" VALUES (")
	args := make([]interface{}, 0)
	n := 1
	fieldsNum := len(q.changes) + len(q.rawChanges)
	for k, v := range q.changes {
		args = append(args, v)
		bs.WriteString("`" + k + "`")
		if n < fieldsNum {
			bs.WriteRune(',')
		}
		bs2.WriteRune('?')
		if n < fieldsNum {
			bs2.WriteRune(',')
		}
		n++
	}
	for k, v := range q.rawChanges {
		bs.WriteString("`" + k + "`")
		if n < fieldsNum {
			bs.WriteRune(',')
		}
		bs2.WriteString(v)
		if n < fieldsNum {
			bs2.WriteRune(',')
		}
		n++
	}
	bs.WriteString(")")
	bs2.WriteString(")")
	bs.Write(bs2.Bytes())
	return bs.String(), args
}

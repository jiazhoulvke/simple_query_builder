package sqb

import "bytes"

func Update(table string) *UpdateQuery {
	q := UpdateQuery{
		Query: Query{
			tableName: table,
		},
		queryChanges: &queryChanges{
			changes:         make(map[string]interface{}),
			changeFields:    make([]string, 0),
			rawChanges:      make(map[string]string),
			rawChangeFields: make([]string, 0),
		},
		queryConditions: &queryConditions{
			conditions: make([]Conditioner, 0),
		},
	}
	return &q
}

type UpdateQuery struct {
	Query
	*queryChanges
	*queryConditions
}

func (q *UpdateQuery) Set(key string, value interface{}) *UpdateQuery {
	q.queryChanges.Set(key, value)
	return q
}

func (q *UpdateQuery) SetRaw(key string, value string) *UpdateQuery {
	q.queryChanges.SetRaw(key, value)
	return q
}

func (q *UpdateQuery) SetData(data map[string]interface{}) *UpdateQuery {
	q.queryChanges.SetData(data)
	return q
}

func (q *UpdateQuery) SetRawData(data map[string]string) *UpdateQuery {
	q.queryChanges.SetRawData(data)
	return q
}

func (q *UpdateQuery) Where(conditions ...Conditioner) *UpdateQuery {
	q.queryConditions.Where(conditions...)
	return q
}

func (q *UpdateQuery) Build() (string, []interface{}) {
	bs := bytes.NewBufferString("UPDATE ")
	bs.WriteString(q.tableName)
	bs.WriteString(" SET ")
	args := make([]interface{}, 0)
	n := 1
	fieldsNum := len(q.changes) + len(q.rawChanges)
	for k, v := range q.changes {
		args = append(args, v)
		bs.WriteString(k + "=?")
		if n != fieldsNum {
			bs.WriteByte(',')
		}
		n++
	}
	for k, v := range q.rawChanges {
		if n == fieldsNum {
			bs.WriteString(k + "=" + v)
		} else {
			bs.WriteString(k + "=" + v + ",")
		}
		n++
	}
	if len(q.conditions) == 0 {
		return bs.String(), []interface{}{}
	}
	bs.WriteString(" WHERE ")
	cStr, cArgs := joinConditions("AND", q.conditions)
	bs.WriteString(" " + cStr)
	args = append(args, cArgs...)
	return bs.String(), args
}

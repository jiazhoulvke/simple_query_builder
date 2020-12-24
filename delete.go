package sqb

import "bytes"

func Delete(table string) *DeleteQuery {
	q := DeleteQuery{
		Query: Query{
			tableName: table,
		},
		queryConditions: &queryConditions{
			conditions: make([]Conditioner, 0),
		},
	}
	return &q
}

type DeleteQuery struct {
	Query
	*queryConditions
}

func (q *DeleteQuery) Where(conditions ...Conditioner) *DeleteQuery {
	q.queryConditions.Where(conditions...)
	return q
}

func (q *DeleteQuery) Build() (string, []interface{}) {
	bs := bytes.NewBufferString("DELETE FROM ")
	bs.WriteString(q.tableName)
	if len(q.conditions) == 0 {
		return bs.String(), []interface{}{}
	}
	bs.WriteString(" WHERE ")
	cStr, args := joinConditions("AND", q.conditions)
	bs.WriteString(" " + cStr)
	return bs.String(), args
}

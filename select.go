package sqb

import (
	"bytes"
	"strconv"
)

func Select(table string) *SelectQuery {
	q := SelectQuery{
		Query: Query{
			tableName: table,
		},
		fields: "*",
		queryJoins: &queryJoins{
			joins: make([]string, 0),
		},
		queryConditions: &queryConditions{
			conditions: make([]Conditioner, 0),
		},
		groupBys: make([]string, 0),
		orderBys: make([]string, 0),
	}
	return &q
}

type SelectQuery struct {
	Query
	isDistinct bool
	fields     string
	*queryJoins
	*queryConditions
	groupBys []string
	having   string
	orderBys []string
	limit    int
	offset   int
}

func (q *SelectQuery) Distinct() *SelectQuery {
	q.isDistinct = true
	return q
}

func (q *SelectQuery) NoDistinct() *SelectQuery {
	q.isDistinct = false
	return q
}

func (q *SelectQuery) Fields(fields string) *SelectQuery {
	q.fields = fields
	return q
}

func (q *SelectQuery) Join(joinExpr string) *SelectQuery {
	q.queryJoins.Join(joinExpr)
	return q
}

func (q *SelectQuery) InnerJoin(tableName string, expr string) *SelectQuery {
	q.queryJoins.InnerJoin(tableName, expr)
	return q
}

func (q *SelectQuery) LeftJoin(tableName string, expr string) *SelectQuery {
	q.queryJoins.LeftJoin(tableName, expr)
	return q
}

func (q *SelectQuery) RightJoin(tableName string, expr string) *SelectQuery {
	q.queryJoins.RightJoin(tableName, expr)
	return q
}

func (q *SelectQuery) Where(conditions ...Conditioner) *SelectQuery {
	q.queryConditions.Where(conditions...)
	return q
}

func (q *SelectQuery) Build() (string, []interface{}) {
	args := make([]interface{}, 0)
	bs := bytes.NewBufferString("SELECT ")
	if q.isDistinct {
		bs.WriteString("DISTINCT ")
	}
	bs.WriteString(q.fields)
	bs.WriteString(" FROM `")
	bs.WriteString(q.tableName)
	bs.WriteByte('`')
	//Join
	if len(q.joins) > 0 {
		for _, j := range q.joins {
			bs.WriteString(j)
			bs.WriteRune(' ')
		}
	}
	//Where
	if len(q.conditions) > 0 {
		bs.WriteString(" WHERE ")
		cStr, cArgs := joinConditions("AND", q.conditions)
		bs.WriteString(" " + cStr)
		args = append(args, cArgs...)
	}
	//GroupBy
	if len(q.groupBys) > 0 {
		bs.WriteString(" GROUP BY ")
		for i := 0; i < len(q.groupBys); i++ {
			bs.WriteString(q.groupBys[i])
			if i < len(q.groupBys)-1 {
				bs.WriteRune(',')
			}
		}
	}
	//Having
	if q.having != "" {
		bs.WriteString(" HAVING ")
		bs.WriteString(q.having)
	}
	//OrderBy
	if len(q.orderBys) > 0 {
		bs.WriteString(" ORDER BY ")
		for i := 0; i < len(q.orderBys); i++ {
			bs.WriteString(q.orderBys[i])
			if i < len(q.orderBys)-1 {
				bs.WriteRune(',')
			}
		}
	}
	//Limit
	if q.limit > 0 {
		bs.WriteString(" LIMIT ")
		bs.WriteString(strconv.FormatInt(int64(q.limit), 10))
	}
	//Offset
	if q.offset > 0 {
		bs.WriteString(" OFFSET ")
		bs.WriteString(strconv.FormatInt(int64(q.offset), 10))
	}
	return bs.String(), args
}

func (q *SelectQuery) GroupBy(str string) *SelectQuery {
	q.groupBys = append(q.groupBys, str)
	return q
}

func (q *SelectQuery) Having(str string) *SelectQuery {
	q.having = str
	return q
}

func (q *SelectQuery) Asc(str string) *SelectQuery {
	q.orderBys = append(q.orderBys, str+" ASC")
	return q
}

func (q *SelectQuery) Desc(str string) *SelectQuery {
	q.orderBys = append(q.orderBys, str+" DESC")
	return q
}

func (q *SelectQuery) Limit(n int) *SelectQuery {
	q.limit = n
	return q
}

func (q *SelectQuery) Offset(n int) *SelectQuery {
	q.offset = n
	return q
}

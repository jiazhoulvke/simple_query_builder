package sqb

import (
	"bytes"
	"strconv"
)

type Model struct {
	method     string
	isDistinct bool
	fields     string
	tableName  string
	joins      []string
	changes    map[string]interface{}
	conditions []Conditioner
	groupBys   []string
	having     string
	orderBys   []string
	limit      int
	offset     int
}

func Select(table string, fields string) *Model {
	m := NewModel(table, "select")
	m.fields = fields
	return m
}

func Insert(table string) *Model {
	return NewModel(table, "insert")
}

func Delete(table string) *Model {
	return NewModel(table, "delete")
}

func Update(table string) *Model {
	return NewModel(table, "update")
}

func NewModel(table string, method string) *Model {
	m := Model{
		fields:     "*",
		tableName:  table,
		method:     method,
		changes:    make(map[string]interface{}),
		joins:      make([]string, 0),
		groupBys:   make([]string, 0),
		conditions: make([]Conditioner, 0),
		orderBys:   make([]string, 0),
		limit:      -1,
		offset:     -1,
	}
	return &m
}

func (m *Model) Distinct() *Model {
	m.isDistinct = true
	return m
}

func (m *Model) NoDistinct() *Model {
	m.isDistinct = false
	return m
}

func (m *Model) Join(t string) *Model {
	m.joins = append(m.joins, t)
	return m
}

func (m *Model) InnerJoin(t string) *Model {
	m.joins = append(m.joins, " INNER JOIN "+t)
	return m
}

func (m *Model) LeftJoin(t string) *Model {
	m.joins = append(m.joins, " LEFT JOIN "+t)
	return m
}

func (m *Model) RightJoin(t string) *Model {
	m.joins = append(m.joins, " RIGHT JOIN "+t)
	return m
}

func (m *Model) Set(key string, value interface{}) *Model {
	m.changes[key] = value
	return m
}

func (m *Model) SetData(data map[string]interface{}) *Model {
	for k, v := range data {
		m.changes[k] = v
	}
	return m
}

func (m *Model) Where(conditions ...Conditioner) *Model {
	m.conditions = append(m.conditions, conditions...)
	return m
}

func (m *Model) GroupBy(str string) *Model {
	m.groupBys = append(m.groupBys, str)
	return m
}

func (m *Model) Having(str string) *Model {
	m.having = str
	return m
}

func (m *Model) Asc(str string) *Model {
	m.orderBys = append(m.orderBys, str+" ASC")
	return m
}

func (m *Model) Desc(str string) *Model {
	m.orderBys = append(m.orderBys, str+" DESC")
	return m
}

func (m *Model) Limit(n int) *Model {
	m.limit = n
	return m
}

func (m *Model) Offset(n int) *Model {
	m.offset = n
	return m
}

func (m *Model) Build() (string, []interface{}) {
	switch m.method {
	case "insert":
		return m.buildInsert()
	case "delete":
		return m.buildDelete()
	case "update":
		return m.buildUpdate()
	case "select":
		return m.buildSelect()
	}
	panic("method error")
}

func (m *Model) buildInsert() (string, []interface{}) {
	bs := bytes.NewBufferString("INSERT INTO `")
	bs.WriteString(m.tableName)
	bs.WriteString("` (")
	bs2 := bytes.NewBufferString(" (")
	args := make([]interface{}, 0)
	n := 1
	fieldsNum := len(m.changes)
	for k, v := range m.changes {
		args = append(args, v)
		if n >= fieldsNum {
			bs.WriteString("`" + k + "`")
			bs2.WriteString("?")
		} else {
			bs.WriteString("`" + k + "`,")
			bs2.WriteString("?,")
		}
		n++
	}
	bs.WriteString(")")
	bs2.WriteString(")")
	bs.Write(bs2.Bytes())
	return bs.String(), args
}

func (m *Model) buildDelete() (string, []interface{}) {
	bs := bytes.NewBufferString("DELETE FROM `")
	bs.WriteString(m.tableName)
	bs.WriteString("`")
	if len(m.conditions) == 0 {
		return bs.String(), []interface{}{}
	}
	bs.WriteString(" WHERE ")
	cStr, args := joinConditions("AND", m.conditions)
	bs.WriteString(" " + cStr)
	return bs.String(), args
}

func (m *Model) buildUpdate() (string, []interface{}) {
	bs := bytes.NewBufferString("UPDATE `")
	bs.WriteString(m.tableName)
	bs.WriteString("` SET ")
	n := 1
	args := make([]interface{}, 0)
	for k, v := range m.changes {
		args = append(args, v)
		if n >= len(m.changes) {
			bs.WriteString(k + "=?")
		} else {
			bs.WriteString(k + "=?,")
		}
		n++
	}
	if len(m.conditions) == 0 {
		return bs.String(), []interface{}{}
	}
	bs.WriteString(" WHERE ")
	cStr, cArgs := joinConditions("AND", m.conditions)
	bs.WriteString(" " + cStr)
	args = append(args, cArgs...)
	return bs.String(), args
}

func (m *Model) buildSelect() (string, []interface{}) {
	args := make([]interface{}, 0)
	bs := bytes.NewBufferString("SELECT ")
	if m.isDistinct {
		bs.WriteString("DISTINCT ")
	}
	bs.WriteString(m.fields)
	bs.WriteString(" FROM ")
	bs.WriteString(m.tableName)
	//Join
	if len(m.joins) > 0 {
		for _, j := range m.joins {
			bs.WriteString(j)
			bs.WriteRune(' ')
		}
	}
	//Where
	if len(m.conditions) > 0 {
		bs.WriteString(" WHERE ")
		cStr, cArgs := joinConditions("AND", m.conditions)
		bs.WriteString(" " + cStr)
		args = append(args, cArgs...)
	}
	//GroupBy
	if len(m.groupBys) > 0 {
		bs.WriteString(" GROUP BY ")
		for i := 0; i < len(m.groupBys); i++ {
			bs.WriteString(m.groupBys[i])
			if i < len(m.groupBys)-1 {
				bs.WriteRune(',')
			}
		}
	}
	//Having
	if m.having != "" {
		bs.WriteRune(' ')
		bs.WriteString(m.having)
	}
	//OrderBy
	if len(m.orderBys) > 0 {
		bs.WriteString(" ORDER BY ")
		for i := 0; i < len(m.orderBys); i++ {
			bs.WriteString(m.orderBys[i])
			if i < len(m.groupBys)-1 {
				bs.WriteRune(',')
			}
		}
	}
	//Limit
	if m.limit > 0 {
		bs.WriteString(" LIMIT ")
		bs.WriteString(strconv.FormatInt(int64(m.limit), 10))
	}
	//Offset
	if m.offset > 0 {
		bs.WriteString(" OFFSET ")
		bs.WriteString(strconv.FormatInt(int64(m.offset), 10))
	}
	return bs.String(), args
}

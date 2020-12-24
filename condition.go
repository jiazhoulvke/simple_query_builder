package sqb

import (
	"bytes"
	"reflect"
)

type Conditioner interface {
	Build() (string, []interface{})
}

type BaseCondition struct {
	Str  string
	args []interface{}
}

func (c BaseCondition) Build() (string, []interface{}) {
	return c.Str, c.args
}

func Equal(expr string, arg interface{}) BaseCondition {
	return BaseCondition{
		Str:  expr + "=?",
		args: []interface{}{arg},
	}
}

func NotEqual(expr string, arg interface{}) BaseCondition {
	return BaseCondition{
		Str:  expr + "<>?",
		args: []interface{}{arg},
	}
}

func LessThan(expr string, arg interface{}) BaseCondition {
	return BaseCondition{
		Str:  expr + "<?",
		args: []interface{}{arg},
	}
}

func LessEqualThan(expr string, arg interface{}) BaseCondition {
	return BaseCondition{
		Str:  expr + "<=?",
		args: []interface{}{arg},
	}
}

func GreaterThan(expr string, arg interface{}) BaseCondition {
	return BaseCondition{
		Str:  expr + ">?",
		args: []interface{}{arg},
	}
}

func GreaterEqualThan(expr string, arg interface{}) BaseCondition {
	return BaseCondition{
		Str:  expr + ">=?",
		args: []interface{}{arg},
	}
}

func IsNull(expr string) BaseCondition {
	return BaseCondition{
		Str:  expr + " IS NULL",
		args: []interface{}{},
	}
}

func IsNotNull(expr string) BaseCondition {
	return BaseCondition{
		Str:  expr + " IS NOT NULL",
		args: []interface{}{},
	}
}

func Like(expr string, str string) BaseCondition {
	return BaseCondition{
		Str:  expr + " LIKE '" + str + "'",
		args: []interface{}{},
	}
}

func NotLike(expr string, str string) BaseCondition {
	return BaseCondition{
		Str:  expr + " NOT LIKE '" + str + "'",
		args: []interface{}{},
	}
}

func inOrNotIn(method string, expr string, args ...interface{}) BaseCondition {
	bs := bytes.NewBufferString(expr)
	bs.WriteRune(' ')
	bs.WriteString(method)
	bs.WriteString(" (")
	for i := 0; i < len(args); i++ {
		if i == len(args)-1 {
			bs.WriteRune('?')
		} else {
			bs.WriteString("?,")
		}
	}
	bs.WriteRune(')')
	return BaseCondition{
		Str:  bs.String(),
		args: args,
	}
}

//Interfaces 将任意类型的切片转为[]interface{}
func Interfaces(slice interface{}) []interface{} {
	s := make([]interface{}, 0, 0)
	v := reflect.ValueOf(slice)
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		panic("value is not slice")
	}

	for i, l := 0, v.Len(); i < l; i++ {
		s = append(s, v.Index(i).Interface())
	}
	return s
}

func In(expr string, args ...interface{}) BaseCondition {
	return inOrNotIn("IN", expr, args...)
}

func NotIn(expr string, args ...interface{}) BaseCondition {
	return inOrNotIn("NOT IN", expr, args...)
}

func And(conditions ...Conditioner) AndCondition {
	return AndCondition{
		conditions: conditions,
	}
}

func Or(conditions ...Conditioner) OrCondition {
	return OrCondition{
		conditions: conditions,
	}
}

type OrCondition struct {
	conditions []Conditioner
}

func (c OrCondition) Build() (string, []interface{}) {
	return joinConditions("OR", c.conditions)
}

type AndCondition struct {
	conditions []Conditioner
}

func (c AndCondition) Build() (string, []interface{}) {
	return joinConditions("AND", c.conditions)
}

func joinConditions(joinType string, conditions []Conditioner) (string, []interface{}) {
	bs := bytes.NewBufferString("")
	args := make([]interface{}, 0)
	for i, cond := range conditions {
		cStr, cArgs := cond.Build()
		if i == len(conditions)-1 {
			bs.WriteString("(" + cStr + ")")
		} else {
			bs.WriteString("(" + cStr + ") " + joinType + " ")
		}
		args = append(args, cArgs...)
	}
	return bs.String(), args
}

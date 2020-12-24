package sqb

type Builder interface {
	Build() (string, []interface{})
}

type Query struct {
	tableName string
}

type queryChanges struct {
	changes         map[string]interface{}
	changeFields    []string
	rawChanges      map[string]string
	rawChangeFields []string
}

func (q *queryChanges) Set(key string, value interface{}) *queryChanges {
	q.changes[key] = value
	return q
}

func (q *queryChanges) SetRaw(key string, value string) *queryChanges {
	q.rawChanges[key] = value
	return q
}

func (q *queryChanges) SetData(data map[string]interface{}) *queryChanges {
	for k, v := range data {
		q.changes[k] = v
	}
	return q
}

func (q *queryChanges) SetRawData(data map[string]string) *queryChanges {
	for k, v := range data {
		q.rawChanges[k] = v
	}
	return q
}

type queryConditions struct {
	conditions []Conditioner
}

func (q *queryConditions) Where(conditions ...Conditioner) *queryConditions {
	q.conditions = append(q.conditions, conditions...)
	return q
}

type queryJoins struct {
	joins []string
}

func (q *queryJoins) Join(joinExpr string) *queryJoins {
	q.joins = append(q.joins, joinExpr)
	return q
}

func (q *queryJoins) InnerJoin(tableName string, expr string) *queryJoins {
	q.joins = append(q.joins, " INNER JOIN "+tableName+" ON "+expr)
	return q
}

func (q *queryJoins) LeftJoin(tableName string, expr string) *queryJoins {
	q.joins = append(q.joins, " LEFT JOIN "+tableName+" ON "+expr)
	return q
}

func (q *queryJoins) RightJoin(tableName string, expr string) *queryJoins {
	q.joins = append(q.joins, " RIGHT JOIN "+tableName+" ON "+expr)
	return q
}

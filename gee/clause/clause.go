package clause

import (
	"strings"
)

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]any
}

func (c *Clause) Set(name Type, vars ...any) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]any)
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

func (c *Clause) Build(orders ...Type) (string, []any) {
	sqls := make([]string, 0)
	vars := make([]any, 0)
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}

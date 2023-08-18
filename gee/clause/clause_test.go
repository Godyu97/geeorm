package clause

import (
	"reflect"
	"testing"
)

func testSelect(t *testing.T) {
	var c Clause
	c.Set(LIMIT, 3)
	c.Set(SELECT, "User", []string{"*"})
	c.Set(WHERE, "Name = ?", "Tom")
	c.Set(ORDERBY, "Age ASC")
	sql, vars := c.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []any{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}

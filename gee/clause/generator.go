package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...any) (string, []any)

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert()
	generators[VALUES] = _values()
	generators[SELECT] = _select()
	generators[LIMIT] = _limit()
	generators[WHERE] = _where()
	generators[ORDERBY] = _order()
	generators[UPDATE] = _update()
	generators[DELETE] = _delete()
	generators[COUNT] = _count()
}

func genBindVars(num int) string {
	vars := make([]string, 0)
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _insert() generator {
	return func(values ...any) (string, []any) {
		//	INSERT INTO $table_name ($fields)
		table := values[0]
		fields := strings.Join(values[1].([]string), ",")
		return fmt.Sprintf("INSERT INTO %s (%s)", table, fields), []any{}
	}
}

func _values() generator {
	return func(values ...any) (string, []any) {
		//	VALUES ($v1),($v2)
		var bindStr string
		var sql strings.Builder
		var vars []any
		sql.WriteString("VALUES ")
		for i, value := range values {
			v := value.([]any)
			if bindStr == "" {
				bindStr = genBindVars(len(v))
			}
			sql.WriteString(fmt.Sprintf("(%s)", bindStr))
			if i+1 != len(values) {
				sql.WriteString(", ")
			}
			vars = append(vars, v...)
		}
		return sql.String(), vars
	}
}

func _select() generator {
	return func(values ...any) (string, []any) {
		//	SELECT $fields FROM $table_name
		table := values[0]
		fields := strings.Join(values[1].([]string), ",")
		return fmt.Sprintf("SELECT %s FROM %s", fields, table), []any{}
	}
}

func _limit() generator {
	return func(values ...any) (string, []any) {
		// LIMIT $num
		return "LIMIT ?", values
	}
}

func _where() generator {
	return func(values ...any) (string, []any) {
		//	WHERE $condition
		condition, vars := values[0], values[1:]
		return fmt.Sprintf("WHERE %s", condition), vars
	}
}

func _order() generator {
	return func(values ...any) (string, []any) {
		return fmt.Sprintf("ORDER BY %s", values[0]), []any{}
	}
}

func _update() generator {
	return func(values ...any) (string, []any) {
		table := values[0]
		m := values[1].(map[string]any)
		var ks []string
		var vs []any
		for k, v := range m {
			ks = append(ks, k+" = ?")
			vs = append(vs, v)
		}
		return fmt.Sprintf("UPDATE %s SET %s", table, strings.Join(ks, ", ")), vs
	}
}

func _delete() generator {
	return func(values ...any) (string, []any) {
		return fmt.Sprintf("DELETE FROM %s", values[0]), []any{}
	}
}

func _count() generator {
	return func(values ...any) (string, []any) {
		return _select()(values[0], []string{"COUNT(*)"})
	}
}

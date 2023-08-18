# geeorm
geektutu orm learn

* 什么是对象关系映射

| 数据库              | 面向对象的编程语言  |
| ------------------- | ------------------- |
| 表(table)           | 类(class/struct)    |
| 记录(record, row)   | 对象 (object)       |
| 字段(field, column) | 对象属性(attribute) |

* Go 语言的反射机制(reflect)，通过反射，可以获取到对象对应的结构体名称，成员变量、方法等信息。根据任意类型的指针，得到其对应的结构体的信息。

	- `reflect.ValueOf()` 获取指针对应的反射值。

	- `reflect.Indirect()` 获取指针指向的对象的反射值。

	- `(reflect.Type).Name()` 返回类名(字符串)。

	- `(reflect.Type).Field(i)` 获取第 i 个成员变量。

* 基本思想：通过链式调用，构造sql语句，给sql.Exec执行
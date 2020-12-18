简单的SQL生成器，是一个只用了一个晚上写的简单玩具，仅供自娱自乐。

## 安装

`go get -v github.com/jiazhoulvke/simple_query_builer`

## 使用

```go
package main

import (
	"fmt"

	sqb "github.com/jiazhoulvke/simple_query_builder"
)

func main() {
	sql, args := sqb.Select("table1", "*").Where(
		sqb.Or(
		    sqb.Equal("field1", 123),
		    sqb.Like("field2", "abc"),
		),
		sqb.Or(
		    sqb.NotEqual("field3", 1.23),
		    sqb.Like("field4", "%foo%"),
		),
	).
	Limit(10).
	Offset(20).
	Build()
	fmt.Println(sql, args)
	//SELECT * FROM table1 WHERE  ((field1=?) OR (field2<>?)) AND ((field3 IS NULL) OR (field4 LIKE '%foo%')) LIMIT 10 OFFSET 20 [123 abc]
}
```

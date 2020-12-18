package sqb

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestModel(t *testing.T) {
	Convey("TestModel", t, func() {
		Convey("TestInsert", func() {
			insertSQL, args := Insert("table1").
				Set("field1", 1).
				Set("field2", "two").
				SetData(map[string]interface{}{
					"field3": 'a',
					"field4": "foobar",
				}).
				Build()
			So(insertSQL, ShouldNotEqual, "")
			So(len(args), ShouldEqual, 4)
			t.Log("INSERT:", insertSQL, args)
		})
		Convey("TestDelete", func() {
			deleteSQL, args := Delete("table1").
				Where(
					Or(
						Or(
							Equal("`field1`", 1),
							Or(
								LessThan("`field2`", 2),
								GreaterThan("field3", "three"),
							),
							LessEqualThan("field4", 1.23),
							GreaterEqualThan("field5", "hello"),
						),
						IsNull("n1"),
						IsNotNull("n2"),
						Like("l1", "%area%"),
						NotLike("l2", "area%"),
						In("id", 1, 2, 3),
						NotIn("id", 4, 5),
					),
				).
				Build()
			So(deleteSQL, ShouldNotEqual, "")
			So(len(args), ShouldEqual, 10)
			t.Log("DELETE:", deleteSQL, args)
		})
		Convey("TestUpdate", func() {
			updateSQL, args := Update("table1").
				Set("field1", 1).
				Set("field2", 2.2).
				Set("field3", "three").
				Set("field4", "4").
				SetData(map[string]interface{}{
					"field5": "foobar5",
				}).
				Where(
					Equal("`status`", "normal"),
					GreaterEqualThan("created_at", 11111111),
				).
				Build()
			So(updateSQL, ShouldNotEqual, "")
			So(len(args), ShouldEqual, 7)
			t.Log("UPDATE:", updateSQL, args)
		})
		Convey("TestSelect", func() {
			selectSQL, args := Select("table1", "*").
				Distinct().
				LeftJoin("table2 ON table2.tid=table1.id").
				RightJoin("table3 ON table3.tid=table1.id").
				Join("LEFT JOIN table4 ON table4.cid=table1.cid").
				Where(
					Equal("field1", "abc"),
					NotEqual("field2", "def"),
				).
				GroupBy("table1.field1").
				GroupBy("table2.field3").
				Having("C>5").
				Asc("id").
				Desc("rank").
				Limit(10).
				Offset(20).
				Build()
			So(selectSQL, ShouldNotEqual, "")
			So(len(args), ShouldEqual, 2)
			t.Log("SELECT:", selectSQL, args)
		})
	})
}

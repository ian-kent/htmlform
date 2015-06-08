package htmlform

import (
	"testing"
	"time"

	. "github.com/smartystreets/GoConvey/convey"
)

func TestMap(t *testing.T) {
	Convey("Map should return a map", t, func() {
		m, err := Map("a", 1)
		So(err, ShouldBeNil)
		So(m, ShouldHaveSameTypeAs, map[string]interface{}{})
	})
	Convey("Map should return an error for invalid arguments", t, func() {
		m, err := Map("a", 1, "b")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "expecting even number of arguments, got 3")
		So(m, ShouldBeNil)

		m, err = Map(1, 1)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "expecting string for odd numbered arguments, got 1")
		So(m, ShouldNotBeNil)
		So(len(m), ShouldEqual, 0)
	})
	Convey("Map should support arbitrary nesting", t, func() {
		m0, err := Map("e", 3)
		So(err, ShouldBeNil)
		m1, err := Map("c", 2, "d", m0)
		So(err, ShouldBeNil)
		m2, err := Map("a", 1, "b", m1)
		So(err, ShouldBeNil)
		So(m0, ShouldHaveSameTypeAs, map[string]interface{}{})
		So(m1, ShouldHaveSameTypeAs, map[string]interface{}{})
		So(m2, ShouldHaveSameTypeAs, map[string]interface{}{})
		So(m2["a"], ShouldEqual, 1)
		So(m2["b"], ShouldEqual, m1)
		So(m2["b"], ShouldHaveSameTypeAs, map[string]interface{}{})
		So(m2["b"].(map[string]interface{})["c"], ShouldEqual, 2)
		So(m2["b"].(map[string]interface{})["d"], ShouldEqual, m0)
		So(m2["b"].(map[string]interface{})["d"], ShouldHaveSameTypeAs, map[string]interface{}{})
		So(m2["b"].(map[string]interface{})["d"].(map[string]interface{})["e"], ShouldEqual, 3)
	})
}

func TestExtend(t *testing.T) {
	Convey("Extend should extend a map", t, func() {
		m0 := map[string]interface{}{"a": 1}
		So(m0["b"], ShouldBeNil)
		m1, err := Extend(m0, "b", 2)
		So(err, ShouldBeNil)
		So(m0["b"], ShouldEqual, 2)
		So(m1, ShouldEqual, m0)
		So(m1["b"], ShouldEqual, 2)
	})
	Convey("Map should return an error for invalid arguments", t, func() {
		m0 := map[string]interface{}{"a": 1}
		m0, err := Extend(m0, "a", 1, "b")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "expecting even number of arguments, got 3")

		m, err := Extend(m0, 1, 1)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "expecting string for odd numbered arguments, got 1")
		So(m, ShouldBeNil)
	})
}

func TestFirstNotNil(t *testing.T) {
	Convey("FirstNotNil should return the first non-nil parameter", t, func() {
		a := FirstNotNil(1, 2, 3)
		So(a, ShouldEqual, 1)
		a = FirstNotNil(0, 1, 2)
		So(a, ShouldEqual, 0)
		a = FirstNotNil(nil, "", "foo")
		So(a, ShouldEqual, "")
		a = FirstNotNil(nil, nil, nil, time.Now())
		So(a, ShouldHaveSameTypeAs, time.Now())
		a = FirstNotNil(nil, nil)
		So(a, ShouldEqual, nil)
	})
}

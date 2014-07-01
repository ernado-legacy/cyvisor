package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestSite(t *testing.T) {
	Convey("TestUrl", t, func() {
		So(testUrl("http://ya.ru"), ShouldBeTrue)
		So(testUrl("http://asdasd.asdasd"), ShouldBeFalse)
		So(Test("http://ya.ru", time.Second), ShouldBeTrue)
		So(Test("http://ya.ru", time.Nanosecond*100), ShouldBeFalse)
	})
}

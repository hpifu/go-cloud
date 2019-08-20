package fs

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	Convey("test", t, func() {
		s := NewFileSystem(".")

		err1 := s.CreateFile("test", "1.txt", []byte("hello world"))
		So(err1, ShouldBeNil)
		err2 := s.CreateFile("test", "2.txt", []byte("hello world"))
		So(err2, ShouldBeNil)
		err3 := s.Remove("test", "1.txt")
		So(err3, ShouldBeNil)
		err4 := s.Remove("test", "3.txt")
		So(err4, ShouldBeNil)
		err5 := s.Remove("test", "")
		So(err5, ShouldBeNil)
	})
}

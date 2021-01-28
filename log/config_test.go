package log

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomOption(t *testing.T) {
	Convey("CustomOptionRoot", t, func() {
		testCustomOptionVarients(t)
	})
}

func testCustomOptionVarients(t *testing.T) {
	newW := func() *wrapper {
		return &wrapper{new(ulLogger)}
	}

	testDirectValue := func() {
		w := newW()
		const a = "A set"

		//value set directly
		opt := CustomOption("SetA", a)
		err := opt(w)
		So(err, ShouldBeNil)
		So(w.ul.A, ShouldEqual, a)
	}

	testDirectValueFunc := func() {
		w := newW()
		const a = "A set"

		//value set via return from func
		opt := CustomOption("SetA", func() string { return a })
		err := opt(w)
		So(err, ShouldBeNil)
		So(w.ul.A, ShouldEqual, a)
	}

	testChainValue := func() {
		w := newW()
		w.ul.A = "old a"
		orig := w.ul

		//reset underlying logger
		opt := CustomOption("WithA", "new a")
		err := opt(w)
		So(err, ShouldBeNil)
		So(w.ul, ShouldNotPointTo, orig)
		So(w.ul.A, ShouldNotEqual, "old a")
	}

	testMultiReturnSuccess := func() {
		w := newW()
		w.ul.C = "reflect me"
		orig := w.ul

		//call method that returns several values and the last is a nil error
		opt := CustomOption("ChainClearCNil", nil)
		err := opt(w)
		So(err, ShouldBeNil)
		So(w.ul, ShouldPointTo, orig)
		So(w.ul.C, ShouldBeEmpty)
	}

	testChainFailure := func() {
		w := newW()
		orig := w.ul
		b := "B"

		//call method that returns several values and the last is a non nil error
		opt := CustomOption("ChainBFailure", func() string { return b })
		err := opt(w)
		So(err, ShouldEqual, errTheSentinel)
		So(w.ul, ShouldPointTo, orig)
		So(w.ul.B, ShouldBeEmpty)
	}

	testBadFunc := func() {
		w := newW()
		b := "B"

		//pass a func that accepts a value, which we don't support
		opt := CustomOption("SetB", func(i int) string { return b })
		err := opt(w)
		So(err, ShouldNotBeNil)
		So(w.ul.B, ShouldNotEqual, b)
	}

	testPanicRecover := func() {
		w := newW()
		orig := w.ul

		//pass func that returns value that is not appropriate for reflected method
		opt := CustomOption("WithA", func() int { return 42 })
		err := opt(w)
		So(err, ShouldNotBeNil)
		So(w.ul, ShouldPointTo, orig)
	}

	testMultiInputErrorChain := func() {
		w := newW()
		orig := w.ul
		a, b := "A", "B"

		//all together now
		opt := CustomOption("WithAB", func() (string, string, error) { return a, b, nil })
		err := opt(w)
		So(err, ShouldBeNil)
		So(w.ul, ShouldNotPointTo, orig)
		So(w.ul.A, ShouldEqual, a)
		So(w.ul.B, ShouldEqual, b)
		So(orig.A, ShouldBeEmpty)
		So(orig.B, ShouldBeEmpty)
	}

	Convey("DirectValue", testDirectValue)
	Convey("DirectValueFunc", testDirectValueFunc)
	Convey("ChainValue", testChainValue)
	Convey("MultiReturnSuccess", testMultiReturnSuccess)
	Convey("ChainFailure", testChainFailure)
	Convey("BadFunc", testBadFunc)
	Convey("PanicRecover", testPanicRecover)
	Convey("MultiInputErrorChain", testMultiInputErrorChain)
}

type wrapper struct {
	ul *ulLogger
}

func (w *wrapper) GetLogger() interface{} {
	return w.ul
}

func (w *wrapper) SetLogger(iface interface{}) {
	if ul, ok := iface.(*ulLogger); ok {
		w.ul = ul
	}
}

type ulLogger struct {
	A, B, C string
}

func (ul *ulLogger) SetA(a string) {
	ul.A = a
}

func (ul *ulLogger) SetB(b string) string {
	ul.B = b
	return b
}

func (ul *ulLogger) WithAB(a, b string) (*ulLogger, error) {
	cpy := *ul
	cpy.SetA(a)
	cpy.SetB(b)

	return &cpy, nil
}

func (ul *ulLogger) WithA(a string) *ulLogger {
	cpy := *ul
	cpy.SetA(a)
	return &cpy
}

func (ul *ulLogger) ChainClearCNil() (*ulLogger, error) {
	ul.C = ""
	return ul, nil
}

var errTheSentinel = errors.New("Oh noooo")

func (ul *ulLogger) ChainBFailure(b string) (*ulLogger, error) {
	return nil, errTheSentinel
}

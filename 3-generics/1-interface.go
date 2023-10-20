package main

type IntOrFloat32 interface {
	int | float32
}

func Foo[T IntOrFloat32](val T) {}

type MyInt int

var (
	a = Foo(1) // OK
	b = Foo(float32(1.0)) // OK
	c = Foo(MyInt(1)) // Not OK (MyInt does not implement IntOrFloat32)
)

type IntOrFloat32Derived interface {
	~int | ~float32
}

func Bar[T IntOrFloat32Derived](val T) {}

var (
	d = Bar(1) // OK
	e = Bar(1.0) // OK
	f = Bar(MyInt(1)) // OK
)

type Addable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
	~float32 | ~float64 |
	~complex64 | ~complex128 |
	~string
}

func Add[T Addable](lhs, rhs T) T {
	return lhs + rhs
}

package main

import (
	"errors"
	"fmt"
	"time"
)

// structs should be ordered from largest value to smallest value so that memory is padded by the smallest value for alignment
type example struct {
	flag    bool
	counter int16
	pi      float32
}

type aStruct struct {
	flag    bool
	counter int16
	pi      float32
}
type bStruct struct {
	flag    bool
	counter int16
	pi      float32
}

func increment(inc *int) {
	*inc++
	fmt.Println("inc: Value of: ", *inc, " Addr of: ", &*inc)
}

type user struct {
	email string
	name  string
}

func createUserV1() user {
	u := user{
		"test@test.com",
		"bill",
	}
	fmt.Printf("V1 %p\n", &u)

	return u
}

func createUserV2() *user {
	u := user{
		"test@test.com",
		"alice",
	}
	fmt.Printf("V2 %p\n", &u)

	return &u
}

func f() {
	fmt.Println("F Function")
}

func strlen(s string, c chan int) {
	c <- len(s)
}

type MyError string

func (e MyError) Error() string {
	return string(e)
}

func foo() error {
	return errors.New("Some Error Occurred")
}

func main() {

	//variables allocated unitialized are initialized using their zero value
	var abcd int

	fmt.Printf("var abcd \t %T [%v]\n", abcd, abcd)

	var a int = 3
	var b string = "Abc"
	var c float64 = 10.1
	var d bool = true

	fmt.Printf("var a \t %T [%v]\n", a, a)
	fmt.Printf("var b \t %T [%v]\n", b, b)
	fmt.Printf("var c \t %T [%v]\n", c, c)
	fmt.Printf("var d \t %T [%v]\n", d, d)

	aa := 3
	bb := "test"
	cc := 123.1
	dd := false

	fmt.Printf("var aa \t %T [%v]\n", aa, aa)
	fmt.Printf("var bb \t %T [%v]\n", bb, bb)
	fmt.Printf("var cc \t %T [%v]\n", cc, cc)
	fmt.Printf("var dd \t %T [%v]\n", dd, dd)

	aaa := int32(10)
	ccc := float32(324.1)

	fmt.Printf("var aaa \t %T [%v]\n", aaa, aaa)
	fmt.Printf("var ccc \t %T [%v]\n", ccc, ccc)

	fmt.Println("Hello, World!")

	//zero value construction / decleration
	var e1 example
	fmt.Printf("var e1 \t %T [%+v]\n", e1, e1)

	e2 := example{
		flag:    true,
		counter: 10,
		pi:      3.141592,
	}

	fmt.Println("Flag", e2.flag)
	fmt.Println("Counter", e2.counter)
	fmt.Println("Pi", e2.pi)

	//you can do literal structure creation in go for structures only used in one place. You dont need to name the struct
	//type in this case because it won't be used in other places.
	var unnamed struct {
		flag    bool
		counter int16
		pi      float32
	}

	fmt.Printf("%T\t%+v\n", unnamed, unnamed)

	//creating the struct for local usage and initializing it in the creation call. This is unnamed structure type
	unamed2 := struct {
		flag    bool
		counter int16
		pi      float32
	}{
		flag:    true,
		counter: 10,
		pi:      3.141592,
	}
	fmt.Printf("%T\t%+v\n", unamed2, unamed2)

	var aSt aStruct
	var bSt bStruct

	//performing explicit conversion since the structure types match we can convert to assign the value to a new var
	bSt = bStruct(aSt)

	fmt.Println(bSt, aSt)

	//we do not need to perform explicit conversion here since the unamed2 structure is a literal structure not a named one
	bSt = unamed2

	fmt.Println(bSt, aSt)

	count := 10

	fmt.Println("Count value: ", count, " Count Address: ", &count)

	increment(&count)

	fmt.Println("Count Value: ", count, " Count Address: ", &count)

	var bill user
	var alice *user

	bill = createUserV1()
	fmt.Printf("%s %p\n", bill, &bill)

	//this is bad way to initialize a variable because we will create a second copy of the data in the heap.
	//this has extra cost compared to the first function variable initialization.
	alice = createUserV2()
	fmt.Printf("%s %p\n", *alice, &alice)

	go f()
	time.Sleep(1 * time.Second)
	fmt.Println("main function")

	channel := make(chan int)
	go strlen("Saluations", channel)
	go strlen("World", channel)
	go strlen("testing this shit", channel)
	go strlen("one more time for the people in the back", channel)

	res1, res2, res3, res4 := <-channel, <-channel, <-channel, <-channel

	fmt.Println(res1, res2, res3, res4)

	if err := foo(); err != nil {
		fmt.Println(err)
	}

}

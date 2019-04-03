package test

import (
	"fmt"
	"testing"
)

func TestIntf(t *testing.T) {
	tempA := ClassA{
		Name: "ClassA",
	}
	fmt.Println("before:", tempA.Name)
	ChangeName(tempA) // ClassB实现了Intf接口的是ClassB类型，所以需要传递的接口的值， 传递给ChangeName的是一个指针
	fmt.Println("after:", tempA.Name)

	tempB := ClassB{
		Name: "ClassB",
	}
	fmt.Println("before:", tempB.Name)
	ChangeName(&tempB) // ClassB实现了Intf接口的是*ClassB类型，所以需要传递的是指针， 传递给ChangeName的是一个指针
	fmt.Println("after:", tempB.Name)
	t.Log("done")
}

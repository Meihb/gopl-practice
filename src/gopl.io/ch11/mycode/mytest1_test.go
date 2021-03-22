package mytest1

import (
	"fmt"
	"testing"
)

func TestCheckSave(t *testing.T) {

	saved := note
	defer func() { note = saved }() //恢复现场

	note = func(name, msg string) { fmt.Println(name, "never says: ", msg) }
	user := "125388094@qq.com"

	CheckSave(user)

}

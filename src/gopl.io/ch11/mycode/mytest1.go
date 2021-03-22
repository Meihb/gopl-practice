package mytest1

import (
	"fmt"
	"time"
)

var note = func(name string, msg string) {
	fmt.Println(name, " says: ", msg)
}

func CheckSave(username string) {
	//do nothing
	msg := fmt.Sprintf("%v : urgent", time.Now())
	note(username, msg)
}

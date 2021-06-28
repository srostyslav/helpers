package helpers

import (
	"encoding/json"
	"fmt"
	"runtime/debug"

	uuid "github.com/satori/go.uuid"
)

func IsEmptyUUID(u uuid.UUID) bool {
	return fmt.Sprint(u) == "00000000-0000-0000-0000-000000000000"
}

func PrettyPrint(v interface{}) (err error) {
	if b, err := json.MarshalIndent(v, "", "  "); err == nil {
		fmt.Println(string(b))
	}
	return err
}

func CatchPanic(title string) {
	if err := recover(); err != nil {
		stack := string(debug.Stack())
		ErrorLogger.Println(err, "panic: "+title+"\n"+stack)
		AdminBot.SendError(stack, "panic: "+title)
	}
}

package server

import (
	"chat-server/src/user"
	"encoding/gob"
	"reflect"
)

func Register() {
	gob.Register(reflect.ValueOf(user.User{}).Interface())
	gob.Register(reflect.ValueOf([]user.User{}).Interface())
}

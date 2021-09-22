package main

import "github.com/LCBHSStudent/xfw-core/src/poet"

type simpleFunc 		func () string
type idGroupFunc 		func (int64, int64)
type idGroupMsgFunc 	func (int64, int64, *string)
type idMsgFunc 			func (int64, *string)
type groupMsgFunc 		func (int64, *string)

var simpleFuncRouter map[string] simpleFunc

func init() {
	simpleFuncRouter = make(map[string] simpleFunc)
	simpleFuncRouter["xfw"] = poet.GetPoetry
	simpleFuncRouter["XFW"] = poet.GetPoetry
	simpleFuncRouter["小飞舞"] = poet.GetPoetry
}

package logic

import (
	"fmt"
)

type Tag uint64

type Tags []Tag

//TODO: here we just compare the tag as it is, but acturely the hierarcky category need more complex compare
//		for example: t = 101 contains t[0] = 1000101
//benchmark:455 ns/op
func (t Tag) In2(userTag Tags) bool {
	if len(userTag) == 0 {
		return false
	}
	tStr := fmt.Sprintf("%d", t)
	for _, ut := range userTag {
		if ut == t {
			return true
		}
		if ut < t {
			continue
		}
		utStr := fmt.Sprintf("%d", t)
		if tStr == utStr[len(utStr)-len(tStr):] {
			return true
		}

	}
	return false
}

//benchmark 87.6 ns/op
func (t Tag) In(userTag Tags) bool {
	if len(userTag) == 0 {
		return false
	}
	index := t.categoryLen()
	for _, ut := range userTag {
		if ut == t {
			return true
		}
		if ut < t {
			continue
		}
		if uint64(ut)%index == uint64(t) {
			return true
		}

	}
	return false
}

func (t Tag) categoryLen() uint64 {
	var i uint64
	for i = 10; i < 100000000000000000; i = i * 10 {
		if uint64(t)/i == 0 {
			return i
		}
	}
	return i
}

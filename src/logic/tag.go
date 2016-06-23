package logic

type Tag uint64

type Tags []Tag

//TODO: here we just compare the tag as it is, but acturely the hierarcky category need more complex compare
//		for example: 101 contains 1000101
func (t Tag) In(userTag Tags) bool {
	if len(userTag) == 0 {
		return false
	}
	for _, ut := range userTag {
		if ut == t {
			return true
		}
	}
	return false
}

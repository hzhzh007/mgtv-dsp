package logic

type Tag uint64

type Tags []Tag

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

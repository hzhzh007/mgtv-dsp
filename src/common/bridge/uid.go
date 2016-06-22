//some common conv
package bridge

import (
	"logic"
	"tag"
)

func LogicUids2RpcUser(uids logic.Uids) (user tag.User) {
	for i := 0; i < len(uids); i++ {
		user.Uids = append(user.Uids, &tag.UserId{Type: uids[i].Type, Value: uids[i].Value})
	}
	return user
}

func RpcUserTag2LogicTag(userTag *tag.UserTag) (tags logic.Tags) {
	if userTag == nil {
		return logic.Tags{}
	}
	tags = make(logic.Tags, 0, len(userTag.Tags))
	for i := 0; i < len(userTag.Tags); i++ {
		tags = append(tags, logic.Tag(userTag.Tags[i].Id))
	}
	return tags
}

package tag

func NewUser(idType, value string) User {
	return User{
		Uids: []*UserId{&UserId{Type: idType, Value: value}},
	}
}

func NewUserFromMap(idMap map[string]string) (user User) {
	//user = make(User)
	for t, v := range idMap {
		user.Uids = append(user.Uids,
			&UserId{Type: t, Value: v})
	}
	return
}

func UserTag2Array(userTag UserTag) (list []uint64) {
	list = make([]uint64, 0, len(userTag.Tags))
	for _, tag := range userTag.Tags {
		list = append(list, tag.Id)
	}
	return
}

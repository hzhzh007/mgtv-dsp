package bridge

import (
	"config"
	. "logic"
)

var (
	activity *Activities = nil
)

//TODO: according reqeust return candidate ads
func GetActivitiesCopy(request interface{}) (activities *Activities, err error) {
	var copyActivities Activities
	activities = &copyActivities
	if activity != nil {
		copyActivities = make(Activities, len(*activity))
		copy(copyActivities, *activity)
		return
		//		return &copyActivities[0:len(copyActivities)], nil
	}
	activity = &Activities{}
	err = config.LoadConfig("config/activity.yaml", activity)
	if err != nil {
		activity = nil
		return nil, err
	}
	copyActivities = make(Activities, len(*activity))
	copy(copyActivities, *activity)
	return
}

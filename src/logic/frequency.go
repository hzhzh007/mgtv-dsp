package logic

const (
	FrequencyPerDayType   = 0
	FrequencyPerWeekType  = 1
	FrequencyPerMonthType = 2
	FrequencyAllType      = 3
	FrequencyCustomType   = 4
)

type FrequencyType []FrequencyStratage

type FrequencyControllType int

type FrequencyStratage struct {
	Type FrequencyControllType
	Num  int
}

func (f *FrequencyStratage) GetType() FrequencyControllType {
	return f.Type
}

func (f *FrequencyStratage) GetNum() int {
	return f.Num
}

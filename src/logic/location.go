package logic

//the location format 北京市	北京市	1156110000
type Location int

const (
	CityMax = 10000
)

func (l Location) Include(cmpLocation Location) bool {
	city := l % CityMax
	if city > 0 {
		return l == cmpLocation
	}
	return l == ((cmpLocation / CityMax) * CityMax)
}

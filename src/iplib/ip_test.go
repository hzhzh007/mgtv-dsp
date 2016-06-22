package iplib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
func init() {
	utils.Init()
	Load(`/Users/tiankui/work/git/PlanB/src/data/ip.csv`)
}
*/

func Test_ipTest(t *testing.T) {
	assert := assert.New(t)
	iplib, err := Load("../data/ip.csv")
	assert.Nil(err)
	assert.NotNil(iplib)
	assert.Equal(uint32(1156430100), iplib.Ip2CityCode("61.187.53.135"))
	assert.Equal(uint32(1000000000), iplib.Ip2CityCode("0.0.0.0"))
	assert.Equal(uint32(1036000000), iplib.Ip2CityCode("1.0.0.0"))
	assert.Equal(uint32(1036000000), iplib.Ip2CityCode("1.0.0.255"))
	assert.Equal(uint32(1036000000), iplib.Ip2CityCode("1.0.0.255"))
	//assert.Equal(uint32(1826000000), iplib.Ip2CityCode("195.81.195.32"))
	/*
		if BinarySearchIp(`61.187.53.135`) != 1156430100 ||
			BinarySearchIp(`0.0.0.0`) != 1000000000 ||
			BinarySearchIp(`1.0.0.0`) != 1036000000 ||
			BinarySearchIp(`1.0.0.255`) != 1036000000 ||
			BinarySearchIp(`1.9.255.255`) != 1458000000 ||
			BinarySearchIp(`195.81.195.32`) != 1826000000 ||
			BinarySearchIp(`195.81.202.63`) != 1826000000 ||
			BinarySearchIp(`195.112.167.84`) != 1156110000 ||
			BinarySearchIp(`195.112.167.87`) != 1156110000 ||
			BinarySearchIp(`195.112.167.86`) != 1156110000 ||
			BinarySearchIp(`202.38.127.0`) != 1156110000 ||
			BinarySearchIp(`202.38.129.255`) != 1156110000 ||
			BinarySearchIp(`202.97.18.0`) != 1156340100 ||
			BinarySearchIp(`202.97.18.2`) != 1156340100 ||
			BinarySearchIp(`1.9.0.0`) != 1458000000 {
			t.Error("ip error")
		}
	*/
}

func BenchmarkLookingFor(b *testing.B) {
	iplib, err := Load("../data/ip.csv")
	assert := assert.New(b)
	assert.Nil(err)

	city := iplib.Ip2CityCode(`202.97.18.2`)
	for i := 0; i < b.N; i++ {
		city = iplib.Ip2CityCode(`202.97.18.2`)
		//	fmt.Sprintf("hello")
	}
	b.Log(city)
}

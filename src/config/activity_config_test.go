package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"mgtv-dsp/logic"
	"testing"
)

func Test_Load(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)
	config := logic.Activities{}
	err := LoadConfig("test/activity_test.yaml", &config)
	assert.Nil(err)
	assert.NotNil(config)
	t.Log(fmt.Sprintf("%+v", config))

	assert.Equal(1, len(config))
	activity1 := config[0]
	assert.Equal(1, activity1.Id)
	assert.Equal(1, len(activity1.ActiveTime))
}

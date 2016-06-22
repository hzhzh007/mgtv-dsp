package main

import (
	"adxmgtv"
	"golang.org/x/net/context"
)

func initHandler(ctx context.Context) {
	adxmgtv.InitHandler(ctx)
}

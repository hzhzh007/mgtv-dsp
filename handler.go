package main

import (
	"golang.org/x/net/context"
	"mgtv-dsp/adx_mgtv"
)

func initHandler(ctx context.Context) {
	adx_mgtv.InitHandler(ctx)
}

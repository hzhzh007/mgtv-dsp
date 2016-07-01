package main

import (
	. "common/consts"
	"config"
	"feedback"
	"flag"
	"github.com/guregu/kami"
	clog "github.com/hzhzh007/context_log"
	"golang.org/x/net/context"
	"iplib"
	_ "net/http/pprof"
	"redis"
	"resource"
	"tag"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "c", "config/config.yaml", "config file")
	flag.Parse()

	//config
	config, err := config.LoadAndSet(configFile)
	if err != nil {
		clog.Log.Fatalf("config parse error", err)
	}

	//log
	clog.InitLog(config.GetLogPath, config.GetLogLevel)
	clog.InitSignal(func() {})
	clog.Log.Debug("config: %+v", config)

	ctx := context.Background()
	ctx = context.WithValue(ctx, ConfKey, config)

	//redis
	redis_pool := redis.NewPool(config.Redis.Addr, config.Redis.PoolNum)
	ctx = context.WithValue(ctx, RedisKey, redis_pool)

	//iplib
	ipLib, err := iplib.Load(config.IpLib)
	if err != nil {
		clog.Log.Fatalf("load ip failed:%s", err)
	}
	ctx = context.WithValue(ctx, IPLibKey, ipLib)

	//activities
	res, err := resource.NewResource(ctx)
	if err != nil {
		clog.Log.Fatalf("load resource error:%s", err)
	}
	ctx = context.WithValue(ctx, ResourceKey, res)

	//feedback
	feedback, err := feedback.NewFeedback(ctx, 10000, 10)
	if err != nil {
		clog.Log.Fatalf("feedback init err:%s", err)
	}
	ctx = context.WithValue(ctx, FeedbackKey, feedback)

	//tag init
	tag.Init(config.Tag.Addr, config.Tag.PoolNum, config.Tag.Timeout)

	kami.Context = ctx
	initHandler(ctx)
	kami.Serve()
}

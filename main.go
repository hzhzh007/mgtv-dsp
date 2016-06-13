package main

import (
	"flag"
	"github.com/guregu/kami"
	cl "github.com/hzhzh007/context_log"
	"golang.org/x/net/context"
	. "mgtv-dsp/common/consts"
	"mgtv-dsp/config"
	"mgtv-dsp/redis"
	_ "net/http/pprof"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "c", "config/config.yaml", "config file")
	flag.Parse()

	config, err := config.LoadAndSet(configFile)
	if err != nil {
		cl.Log.Panic("config parse error", err)
	}

	cl.InitLog(config.GetLogPath, config.GetLogLevel)
	cl.Log.Debug("config: %+v", config)

	ctx := context.Background()
	ctx = context.WithValue(ctx, ConfKey, config)

	redis_pool := redis.NewPool(config.Redis.Addr, config.Redis.PoolNum)
	ctx = context.WithValue(ctx, RedisKey, redis_pool)

	//	mysql_conn, _ := db_mysql.Open(config.Mysql.Addr)
	//TODO check error
	//ctx = context.WithValue(ctx, MysqlKey, mysql_conn)

	kami.Context = ctx
	initHandler(ctx)
	kami.Serve()
}

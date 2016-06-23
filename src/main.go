package main

import (
	. "common/consts"
	"config"
	"flag"
	"github.com/guregu/kami"
	clog "github.com/hzhzh007/context_log"
	"golang.org/x/net/context"
	"iplib"
	_ "net/http/pprof"
	"redis"
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
		clog.Log.Panic("config parse error", err)
	}

	//log
	clog.InitLog(config.GetLogPath, config.GetLogLevel)
	clog.Log.Debug("config: %+v", config)

	ctx := context.Background()
	ctx = context.WithValue(ctx, ConfKey, config)

	//redis
	redis_pool := redis.NewPool(config.Redis.Addr, config.Redis.PoolNum)
	ctx = context.WithValue(ctx, RedisKey, redis_pool)

	//iplib
	ipLib, err := iplib.Load(config.IpLib)
	if err != nil {
		clog.Log.Panic("load ip failed:%s", err)
	}
	ctx = context.WithValue(ctx, IPLibKey, ipLib)

	//tag init
	tag.Init(config.Tag.Addr, config.Tag.PoolNum, config.Tag.Timeout)

	//	mysql_conn, _ := db_mysql.Open(config.Mysql.Addr)
	//TODO check error
	//ctx = context.WithValue(ctx, MysqlKey, mysql_conn)

	kami.Context = ctx
	initHandler(ctx)
	kami.Serve()
}

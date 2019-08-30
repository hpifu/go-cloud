package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-cloud/internal/cloud"
	"github.com/hpifu/go-kit/hhttp"
	"github.com/hpifu/go-kit/logger"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// AppVersion name
var AppVersion = "unknown"

func main() {
	version := flag.Bool("v", false, "print current version")
	configfile := flag.String("c", "configs/cloud.json", "config file path")
	flag.Parse()
	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	// load config
	config := viper.New()
	config.SetEnvPrefix("account")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()
	config.SetConfigType("json")
	fp, err := os.Open(*configfile)
	if err != nil {
		panic(err)
	}
	err = config.ReadConfig(fp)
	if err != nil {
		panic(err)
	}

	// init logger
	infoLog, err := logger.NewTextLoggerWithViper(config.Sub("logger.infoLog"))
	if err != nil {
		panic(err)
	}
	warnLog, err := logger.NewTextLoggerWithViper(config.Sub("logger.warnLog"))
	if err != nil {
		panic(err)
	}
	accessLog, err := logger.NewJsonLoggerWithViper(config.Sub("logger.accessLog"))
	if err != nil {
		panic(err)
	}
	cloud.InfoLog = infoLog
	cloud.WarnLog = warnLog
	cloud.AccessLog = accessLog

	client := hhttp.NewHttpClient(
		config.GetInt("pool.maxConn"),
		config.GetDuration("pool.connTimeout"),
		config.GetDuration("pool.recvTimeout"),
	)

	secure := config.GetBool("service.secure")
	domain := config.GetString("service.domain")

	service := cloud.NewService(
		config.GetString("service.root"),
		config.GetString("api.account"),
		client,
	)

	// init gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		if secure {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "https://"+domain)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "http://"+domain)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// set handler
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
	r.POST("/upload", service.Upload)
	r.GET("/resource", service.Resource)

	infoLog.Infof("%v init success, port [%v]", os.Args[0], config.GetString("service.port"))

	// run server
	if err := r.Run(config.GetString("service.port")); err != nil {
		panic(err)
	}
}

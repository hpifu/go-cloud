package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/pkg/account"
	"github.com/hpifu/go-cloud/internal/service"
	"github.com/hpifu/go-kit/hhttp"
	"github.com/hpifu/go-kit/logger"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/sohlich/elogrus.v7"
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
	config.SetEnvPrefix("cloud")
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
	// init logger
	infoLog, warnLog, accessLog, err := logger.NewLoggerGroupWithViper(config.Sub("logger"))
	if err != nil {
		panic(err)
	}
	esclient, err := elastic.NewClient(
		elastic.SetURL(config.GetString("es.uri")),
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	hook, err := elogrus.NewAsyncElasticHook(esclient, "go-cloud", logrus.InfoLevel, "go-cloud-log")
	if err != nil {
		panic(err)
	}
	accessLog.Hooks.Add(hook)

	client := account.NewClient(
		config.GetString("account.address"),
		config.GetInt("account.maxConn"),
		config.GetDuration("account.connTimeout"),
		config.GetDuration("account.recvTimeout"),
	)

	secure := config.GetBool("service.cookieSecure")
	domain := config.GetString("service.cookieDomain")
	origins := config.GetStringSlice("service.allowOrigins")

	svc := service.NewService(
		config.GetString("service.root"),
		config.GetString("api.account"),
		client, secure, domain,
	)
	svc.SetLogger(infoLog, warnLog, accessLog)

	// init gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"PUT", "POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
	}))

	// set handler
	d := hhttp.NewGinHttpDecorator(infoLog, warnLog, accessLog)
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
	r.POST("/upload/:id", d.Decorate(svc.Upload))
	r.GET("/resource/:id", d.Decorate(svc.Resource))
	r.POST("/avatar/:id", d.Decorate(svc.POSTAccountAvatar))
	r.GET("/avatar/:id", d.Decorate(svc.GETAccountAvatar))
	r.POST("/techimg/:id", d.Decorate(svc.POSTTechImg))
	r.GET("/techimg/:id", d.Decorate(svc.GETTechImg))

	infoLog.Infof("%v init success, port [%v]", os.Args[0], config.GetString("service.port"))

	// run server
	server := &http.Server{
		Addr:    config.GetString("service.port"),
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// graceful quit
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	infoLog.Infof("%v shutdown ...", os.Args[0])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		warnLog.Errorf("%v shutdown fail or timeout", os.Args[0])
		return
	}
	warnLog.Out.(*rotatelogs.RotateLogs).Close()
	accessLog.Out.(*rotatelogs.RotateLogs).Close()
	infoLog.Errorf("%v shutdown success", os.Args[0])
	infoLog.Out.(*rotatelogs.RotateLogs).Close()
}

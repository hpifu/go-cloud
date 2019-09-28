package cloud

import (
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/hhttp"
	"github.com/sirupsen/logrus"
)

var InfoLog *logrus.Logger
var WarnLog *logrus.Logger
var AccessLog *logrus.Logger

func init() {
	InfoLog = logrus.New()
	WarnLog = logrus.New()
	AccessLog = logrus.New()
}

type FileRes struct {
	Filename string
}

type Service struct {
	Root       string
	client     *hhttp.HttpClient
	apiAccount string
	secure     bool
	domain     string
}

func NewService(root string, apiAccount string, client *hhttp.HttpClient, secure bool, domain string) *Service {
	return &Service{
		Root:       root,
		client:     client,
		apiAccount: apiAccount,
		secure:     secure,
		domain:     domain,
	}
}

func Decorator(inner func(*gin.Context) (interface{}, interface{}, int, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		rid := c.DefaultQuery("rid", NewToken())
		req, res, status, err := inner(c)
		if err != nil {
			c.String(status, err.Error())
			WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		} else if res == nil {
			c.Status(status)
		} else {
			switch res.(type) {
			case string:
				c.String(status, res.(string))
			case *FileRes:
				c.File(res.(*FileRes).Filename)
			default:
				c.JSON(status, res)
			}
		}

		AccessLog.WithFields(logrus.Fields{
			"host":   c.Request.Host,
			"url":    c.Request.URL.String(),
			"req":    req,
			"res":    res,
			"rid":    rid,
			"err":    err,
			"status": status,
		}).Info()
	}
}

func NewToken() string {
	buf := make([]byte, 32)
	token := make([]byte, 16)
	rand.New(rand.NewSource(time.Now().UnixNano())).Read(token)
	hex.Encode(buf, token)
	return string(buf)
}

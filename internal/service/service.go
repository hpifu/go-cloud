package service

import (
	"encoding/hex"
	"fmt"
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
			"client":    c.ClientIP(),
			"userAgent": c.GetHeader("User-Agent"),
			"host":      c.Request.Host,
			"url":       c.Request.URL.String(),
			"req":       req,
			"res":       res,
			"rid":       rid,
			"err":       fmt.Sprintf("%v", err),
			"status":    status,
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

type Account struct {
	ID        int    `form:"id" json:"id,omitempty"`
	Email     string `form:"email" json:"email,omitempty"`
	Phone     string `form:"phone" json:"phone,omitempty"`
	FirstName string `form:"firstName" json:"firstName,omitempty"`
	LastName  string `form:"lastName" json:"lastName,omitempty"`
	Birthday  string `form:"birthday" json:"birthday,omitempty"`
	Password  string `form:"password" json:"password,omitempty"`
	Gender    int    `form:"gender" json:"gender"`
	Avatar    string `form:"avatar" json:"avatar"`
}

func (s *Service) getAccount(token string) (*Account, error) {
	res := &Account{}
	if err := s.client.GET("http://"+s.apiAccount+"/account/"+token, nil, nil).Interface(res); err != nil {
		return nil, err
	}

	return res, nil
}

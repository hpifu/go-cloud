package service

import (
	"github.com/hpifu/go-account/pkg/account"
	"github.com/sirupsen/logrus"
)

var InfoLog *logrus.Logger = logrus.New()
var WarnLog *logrus.Logger = logrus.New()
var AccessLog *logrus.Logger = logrus.New()

type Service struct {
	Root       string
	client     *account.Client
	apiAccount string
	secure     bool
	domain     string
}

func NewService(root string, apiAccount string, client *account.Client, secure bool, domain string) *Service {
	return &Service{
		Root:       root,
		client:     client,
		apiAccount: apiAccount,
		secure:     secure,
		domain:     domain,
	}
}

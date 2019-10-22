package service

import (
	"github.com/hpifu/go-account/pkg/account"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Root       string
	client     *account.Client
	apiAccount string
	secure     bool
	domain     string
	infoLog    *logrus.Logger
	warnLog    *logrus.Logger
	accessLog  *logrus.Logger
}

func (s *Service) SetLogger(infoLog, warnLog, accessLog *logrus.Logger) {
	s.infoLog = infoLog
	s.warnLog = warnLog
	s.accessLog = accessLog
}

func NewService(root string, apiAccount string, client *account.Client, secure bool, domain string) *Service {
	return &Service{
		Root:       root,
		client:     client,
		apiAccount: apiAccount,
		secure:     secure,
		domain:     domain,
		infoLog:    logrus.New(),
		warnLog:    logrus.New(),
		accessLog:  logrus.New(),
	}
}

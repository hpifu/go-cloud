package cloud

import (
	"encoding/hex"
	"fmt"
	"github.com/hpifu/go-kit/cpool"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

var InfoLog *logrus.Logger
var WarnLog *logrus.Logger
var AccessLog *logrus.Logger

func init() {
	InfoLog = logrus.New()
	WarnLog = logrus.New()
	AccessLog = logrus.New()
}

type Service struct {
	Root       string
	pool       *cpool.HttpPool
	apiAccount string
}

func NewService(root string, apiAccount string, pool *cpool.HttpPool) *Service {
	return &Service{
		Root:       root,
		pool:       pool,
		apiAccount: apiAccount,
	}
}

func NewToken() string {
	buf := make([]byte, 32)
	token := make([]byte, 16)
	rand.New(rand.NewSource(time.Now().UnixNano())).Read(token)
	hex.Encode(buf, token)
	return string(buf)
}

func NewCode() string {
	return fmt.Sprintf("%06d", int(rand.NewSource(time.Now().UnixNano()).(rand.Source64).Uint64()%1000000))
}

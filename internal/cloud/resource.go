package cloud

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
	"github.com/sirupsen/logrus"
)

type ResourceReqBody struct {
	Token string `json:"token,omitempty"`
	Name  string `json:"name,omitempty"`
}

func (s *Service) Resource(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &ResourceReqBody{
		Token: c.DefaultQuery("token", ""),
		Name:  c.DefaultQuery("name", ""),
	}
	var err error
	var buf []byte
	status := http.StatusOK

	defer func() {
		AccessLog.WithFields(logrus.Fields{
			"host":   c.Request.Host,
			"body":   string(buf),
			"url":    c.Request.URL.String(),
			"req":    req,
			"rid":    rid,
			"err":    err,
			"status": status,
		}).Info()
	}()

	if err = s.checkResourceReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. body: [%v], err: [%v]", string(buf), err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	a, err := s.getAccount(req.Token)
	if err != nil {
		err = fmt.Errorf("get account failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	if a == nil {
		status = http.StatusBadRequest
		c.String(status, "bad token")
		return
	}

	status = http.StatusOK
	c.File(filepath.Join(s.Root, strconv.Itoa(a.ID), req.Name))
}

func (s *Service) checkResourceReqBody(req *ResourceReqBody) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Token: {rule.Required},
		req.Name:  {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

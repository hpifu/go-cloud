package cloud

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
)

type ResourceReqBody struct {
	Token string `json:"token,omitempty"`
	Name  string `json:"name,omitempty"`
}

type ResourceResBody struct {
	OK bool `json:"ok"`
}

func (s *Service) Resource(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &ResourceReqBody{
		Token: c.DefaultQuery("token", ""),
		Name:  c.DefaultQuery("name", ""),
	}
	var res *ResourceResBody
	var err error
	var buf []byte
	status := http.StatusOK

	defer func() {
		AccessLog.WithFields(logrus.Fields{
			"host":   c.Request.Host,
			"body":   string(buf),
			"url":    c.Request.URL.String(),
			"req":    req,
			"res":    res,
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

	res, err = s.resource(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("resource failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	status = http.StatusOK
	// c.JSON(status, res)
	c.File(filepath.Join(s.Root, req.Name))
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

func (s *Service) resource(req *ResourceReqBody) (*ResourceResBody, error) {

	return nil, nil
}

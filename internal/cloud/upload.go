package cloud

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type UploadReqBody struct {
	Token string `json:"token,omitempty"`
	c     *gin.Context
}

type UploadResBody struct {
	OK bool `json:"ok"`
}

func (s *Service) Upload(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &UploadReqBody{
		Token: c.DefaultQuery("token", ""),
		c:     c,
	}
	var res *UploadResBody
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

	if err = s.checkUploadReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. body: [%v], err: [%v]", string(buf), err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.upload(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("upload failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	status = http.StatusOK
	c.JSON(status, res)
}

func (s *Service) checkUploadReqBody(req *UploadReqBody) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Token: {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) upload(req *UploadReqBody) (*UploadResBody, error) {
	fh, err := req.c.FormFile("file")
	if err != nil {
		return nil, err
	}

	src, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	out, err := os.Create(filepath.Join(s.Root, filepath.Base(fh.Filename)))
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return nil, err
	}

	return &UploadResBody{OK: false}, nil
}

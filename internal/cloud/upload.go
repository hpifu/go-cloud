package cloud

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/pkg/account"
	"github.com/hpifu/go-kit/rule"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type UploadReqBody struct {
	Token string `json:"token,omitempty"`
}

type UploadResBody struct {
	OK bool `json:"ok"`
}

func (s *Service) Upload(ctx *gin.Context) {
	rid := ctx.DefaultQuery("rid", NewToken())
	req := &UploadReqBody{
		Token: ctx.DefaultQuery("token", ""),
	}
	var res *UploadResBody
	var err error
	var buf []byte
	status := http.StatusOK

	defer func() {
		AccessLog.WithFields(logrus.Fields{
			"host":   ctx.Request.Host,
			"body":   string(buf),
			"url":    ctx.Request.URL.String(),
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
		ctx.String(status, err.Error())
		return
	}

	a, err := s.getAccount(req.Token)
	if err != nil {
		err = fmt.Errorf("get account failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusInternalServerError
		ctx.String(status, err.Error())
		return
	}

	if a == nil {
		status = http.StatusBadRequest
		ctx.String(status, "bad token")
		return
	}

	if err := s.upload(ctx, a); err != nil {
		err = fmt.Errorf("upload failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusInternalServerError
		ctx.String(status, err.Error())
		return
	}

	status = http.StatusOK
	ctx.JSON(status, res)
}

func (s *Service) checkUploadReqBody(req *UploadReqBody) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Token: {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) getAccount(token string) (*account.Account, error) {
	hreq := &account.GetAccountReq{
		Token: token,
	}
	hres := &account.GetAccountRes{}
	if err := s.client.GET("http://"+s.apiAccount+"/getaccount", hreq, hres); err != nil {
		return nil, err
	}

	return hres.Account, nil
}

func (s *Service) upload(ctx *gin.Context, a *account.Account) error {
	fh, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err := os.MkdirAll(path.Join(s.Root, strconv.Itoa(a.ID)), 0755); err != nil {
		return err
	}
	out, err := os.Create(filepath.Join(s.Root, strconv.Itoa(a.ID), filepath.Base(fh.Filename)))
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return err
	}

	return nil
}

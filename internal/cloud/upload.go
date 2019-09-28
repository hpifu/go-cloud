package cloud

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
)

type UploadReq struct {
	Token string `json:"token,omitempty" uri:"token"`
}

func (s *Service) Upload(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &UploadReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := s.validUpdate(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	account, err := s.getAccount(req.Token)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get account failed. err: [%v]", err)
	}

	if account == nil {
		return req, "bad token", http.StatusBadRequest, nil
	}

	if err := s.upload(c, account); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("upload failed. err: [%v]", err)
	}

	return req, nil, http.StatusOK, nil
}

func (s *Service) validUpdate(req *UploadReq) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Token: {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) upload(ctx *gin.Context, a *Account) error {
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

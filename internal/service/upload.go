package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
)

type UploadReq struct {
	ID    int    `json:"id,omitempty" uri:"id"`
	Token string `json:"token,omitempty"`
	Name  string `json:"name,omitempty" form:"name"`
}

func (s *Service) Upload(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &UploadReq{
		Token: c.GetHeader("Authorization"),
	}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := c.BindQuery(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validUpload(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	account, err := s.getAccount(req.Token)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get account failed. err: [%v]", err)
	}

	if account == nil {
		return req, nil, http.StatusForbidden, fmt.Errorf("授权失败，请重新登陆")
	}

	if account.ID != req.ID {
		return req, nil, http.StatusForbidden, fmt.Errorf("您没有该资源的权限")
	}

	if err := s.upload(c, account.ID, req.Name); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("upload failed. err: [%v]", err)
	}

	return req, nil, http.StatusOK, nil
}

func (s *Service) validUpload(req *UploadReq) error {
	if err := rule.Check([][3]interface{}{
		{"token", req.Token, []rule.Rule{rule.Required}},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) upload(ctx *gin.Context, id int, name string) error {
	fh, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if name == "" {
		name = filepath.Base(fh.Filename)
	}
	path := filepath.Join(s.Root, strconv.Itoa(id), name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return err
	}

	return nil
}

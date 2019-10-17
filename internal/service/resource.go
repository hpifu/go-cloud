package service

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/hhttp"
	"github.com/hpifu/go-kit/rule"
)

type ResourceReq struct {
	ID    int    `json:"id,omitempty" uri:"id"`
	Token string `json:"token,omitempty" form:"token"`
	Name  string `json:"name,omitempty" form:"name"`
}

func (s *Service) Resource(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &ResourceReq{
		Token: c.GetHeader("Authorization"),
	}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validResource(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	account, err := s.client.GETAccountToken(rid, req.Token)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get account failed. err: [%v]", err)
	}

	if account == nil {
		return req, nil, http.StatusForbidden, fmt.Errorf("授权失败，请重新登陆")
	}

	if account.ID != req.ID {
		return req, nil, http.StatusForbidden, fmt.Errorf("您没有该资源的权限")
	}

	return req, &hhttp.FileRes{Filename: filepath.Join(s.Root, strconv.Itoa(account.ID), req.Name)}, http.StatusOK, nil
}

func (s *Service) validResource(req *ResourceReq) error {
	if err := rule.Check([][3]interface{}{
		{"token", req.Token, []rule.Rule{rule.Required}},
		{"name", req.Name, []rule.Rule{rule.Required}},
	}); err != nil {
		return err
	}

	return nil
}

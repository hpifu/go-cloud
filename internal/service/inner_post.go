package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
)

func (s *Service) POSTTechImg(c *gin.Context) (interface{}, interface{}, int, error) {
	return s.InnerPOST(c, "_pub/tech/img")
}

func (s *Service) POSTAccountAvatar(c *gin.Context) (interface{}, interface{}, int, error) {
	return s.InnerPOST(c, "_pub/account/avatar")
}

type InnerPOSTReq struct {
	ID    int    `json:"id,omitempty" uri:"id"`
	Token string `json:"token,omitempty"`
	Name  string `json:"name,omitempty" form:"name"`
}

func (s *Service) InnerPOST(c *gin.Context, directory string) (interface{}, interface{}, int, error) {
	req := &InnerPOSTReq{
		Token: c.GetHeader("Authorization"),
	}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := c.BindQuery(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validInnerPOST(req); err != nil {
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

	if err := s.upload(c, account.ID, directory+"/"+req.Name); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("InnerPOST failed. err: [%v]", err)
	}

	return req, nil, http.StatusOK, nil
}

func (s *Service) validInnerPOST(req *InnerPOSTReq) error {
	if err := rule.Check([][3]interface{}{
		{"token", req.Token, []rule.Rule{rule.Required}},
		{"name", req.Name, []rule.Rule{rule.Required}},
	}); err != nil {
		return err
	}

	return nil
}

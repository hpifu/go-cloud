package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
)

type POSTAvatarReq struct {
	ID    int    `json:"id,omitempty" uri:"id"`
	Token string `json:"token,omitempty"`
	Name  string `json:"name,omitempty" form:"name"`
}

func (s *Service) POSTAvatar(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &POSTAvatarReq{
		Token: c.GetHeader("Authorization"),
	}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := c.BindQuery(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validPOSTAvatar(req); err != nil {
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

	if err := s.upload(c, account.ID, "_pub/account/avatar/"+req.Name); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("postAvatar failed. err: [%v]", err)
	}

	return req, nil, http.StatusOK, nil
}

func (s *Service) validPOSTAvatar(req *POSTAvatarReq) error {
	if err := rule.Check([][3]interface{}{
		{"token", req.Token, []rule.Rule{rule.Required}},
		{"name", req.Name, []rule.Rule{rule.Required}},
	}); err != nil {
		return err
	}

	return nil
}

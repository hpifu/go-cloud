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

type GetAvatarReq struct {
	ID   int    `json:"id,omitempty" uri:"id"`
	Name string `json:"name,omitempty" form:"name"`
}

func (s *Service) GetAvatar(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &GetAvatarReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validGETAvatar(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	fmt.Println(req, "hello 123")
	return req, &hhttp.FileRes{Filename: filepath.Join(s.Root, strconv.Itoa(req.ID), "_pub/account/avatar", req.Name)}, http.StatusOK, nil
}

func (s *Service) validGETAvatar(req *GetAvatarReq) error {
	if err := rule.Check([][3]interface{}{
		{"name", req.Name, []rule.Rule{rule.Required}},
	}); err != nil {
		return err
	}

	return nil
}

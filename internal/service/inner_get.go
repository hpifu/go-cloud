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

func (s *Service) GETTechImg(c *gin.Context) (interface{}, interface{}, int, error) {
	return s.InnerGET(c, "_pub/tech/img")
}

func (s *Service) GETAccountAvatar(c *gin.Context) (interface{}, interface{}, int, error) {
	return s.InnerGET(c, "_pub/account/avatar")
}

type InnerGETReq struct {
	ID   int    `json:"id,omitempty" uri:"id"`
	Name string `json:"name,omitempty" form:"name"`
}

func (s *Service) InnerGET(c *gin.Context, directory string) (interface{}, interface{}, int, error) {
	req := &InnerGETReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validInnerGET(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	return req, &hhttp.FileRes{Filename: filepath.Join(s.Root, strconv.Itoa(req.ID), directory, req.Name)}, http.StatusOK, nil
}

func (s *Service) validInnerGET(req *InnerGETReq) error {
	if err := rule.Check([][3]interface{}{
		{"name", req.Name, []rule.Rule{rule.Required}},
	}); err != nil {
		return err
	}

	return nil
}

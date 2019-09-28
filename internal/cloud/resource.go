package cloud

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
)

type ResourceReq struct {
	Token string `json:"token,omitempty" uri:"token"`
	Name  string `json:"name,omitempty" form:"name"`
}

func (s *Service) Resource(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &ResourceReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validResource(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	account, err := s.getAccount(req.Token)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get account failed. err: [%v]", err)
	}

	if account == nil {
		return req, "bad token", http.StatusBadRequest, nil
	}

	return req, &FileRes{Filename: filepath.Join(s.Root, strconv.Itoa(account.ID), req.Name)}, http.StatusOK, nil
}

func (s *Service) validResource(req *ResourceReq) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Token: {rule.Required},
		req.Name:  {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

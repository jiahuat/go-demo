package models

import "errors"

type Cluster struct {
	Name string
}

type CreateClusterReq struct {
	Name string `json:"name"`
}
type CreateClusterRes struct {
	ThisIsRes string `json:"response"`
}

func (r *CreateClusterReq) Validate() error {
	if r == nil {
		return errors.New("nil req")
	}
	if r.Name == "" {
		return errors.New("empty name")
	}

	return nil
}

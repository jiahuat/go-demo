package service

import (
	"context"

	appv1 "k8s.io/api/apps/v1"

	k8s "westone.com/wscp-restful/pkg/adapters/k8s"
	models "westone.com/wscp-restful/pkg/cluster/models"
)

// cluster service
type Service struct {
	Conf      *Conf
	k8sClient k8s.InterfaceK8s
}

type Conf struct {
	A bool
	B int64
	C string
}

func NewService(kc k8s.InterfaceK8s, conf *Conf) *Service {

	return &Service{
		Conf:      conf,
		k8sClient: kc,
	}
}

func (s *Service) CreateCluster(ctx context.Context, req *models.CreateClusterReq) error {
	if err := req.Validate(); err != nil {
		return err
	}

	err := s.k8sClient.Create(ctx, &appv1.Deployment{})

	return err
}

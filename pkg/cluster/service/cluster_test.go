package service

import (
	"context"
	"testing"

	"github.com/jiahuat/go-demo/pkg/adapters/k8s/mocks"
	models "github.com/jiahuat/go-demo/pkg/cluster/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCluster(t *testing.T) {
	k8sI := &mocks.InterfaceK8s{}
	k8sI.On("Create", mock.Anything, mock.Anything).Return(nil)

	svc := NewService(k8sI, nil)
	// with err
	err := svc.CreateCluster(context.Background(), &models.CreateClusterReq{})

	assert.NotNil(t, err)
	// nil err
	err = svc.CreateCluster(context.Background(), &models.CreateClusterReq{Name: "name"})
	assert.Nil(t, err)
}

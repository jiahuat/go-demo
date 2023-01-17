package k8s

import (
	"context"

	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type InterfaceK8s interface {
	Create(ctx context.Context, d *appv1.Deployment) error
}

type WscpK8sClientConf struct {
	KubeConfPath string
}

type wscpK8sClient struct {
	kubeClient kubernetes.Interface
}

func NewInerfaceK8s(conf *WscpK8sClientConf) (InterfaceK8s, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", conf.KubeConfPath)
	if err != nil {
		return nil, err
	}
	c, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &wscpK8sClient{
		kubeClient: c,
	}, nil
}

func (c *wscpK8sClient) Create(ctx context.Context, d *appv1.Deployment) error {
	_, err := c.kubeClient.AppsV1().Deployments("").Create(ctx, d, v1.CreateOptions{})

	return err
}

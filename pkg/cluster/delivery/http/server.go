package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	k8s "westone.com/wscp-restful/pkg/adapters/k8s"
	"westone.com/wscp-restful/pkg/cluster/models"
	"westone.com/wscp-restful/pkg/cluster/service"

	"westone.com/wscp-restful/pkg/option"
)

// http server
type Server struct {
	Service service.Service
}

func newServer(svc service.Service) *Server {

	return &Server{Service: svc}
}

func RegisHandlers(r *gin.RouterGroup, conf *option.Config) error {
	// new adapter
	k8sInterfae, err := k8s.NewInerfaceK8s(&k8s.WscpK8sClientConf{
		KubeConfPath: conf.K8s.KubeConfPath,
	})
	if err != nil {
		return err
	}

	// new service
	svc := service.NewService(k8sInterfae, &service.Conf{})
	s := newServer(*svc)

	// register handlers
	r.POST("/", s.CreateCluster)

	return nil
}

// CreateCluster godoc
// @Summary      Create a  cluster
// @Description  This is the description
// @Tags         cluster
// @Accept       json
// @Produce      json
// @Param        req   body    models.CreateClusterReq  true  "create cluster req"
// @Success      200  {object}  models.CreateClusterRes
// @Router       /cluster [post]
func (s Server) CreateCluster(c *gin.Context) {
	req := models.CreateClusterReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	ctx := c.Request.Context()
	if err := s.Service.CreateCluster(ctx, &req); err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

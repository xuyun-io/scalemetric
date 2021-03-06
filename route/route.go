package route

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/xuyun-io/scalemetric/api/v1"
)

// InitRoute define init route.
func InitRoute() *gin.Engine {
	router := gin.Default()
	router.POST("/api/cluster/scalemetric", v1.ClusterPodRequestScheduling)
	return router
}

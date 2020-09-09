package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuyun-io/scalemetric/pkg/calculate"
	"github.com/xuyun-io/scalemetric/pkg/clientset"
	"github.com/xuyun-io/scalemetric/pkg/resources"
	v1 "k8s.io/api/core/v1"
)

// ErrorResp err resp.
func ErrorResp(errstr string) gin.H {
	return gin.H{
		"code":  "1",
		"error": errstr,
	}
}

// OK ok resp.
func OK() gin.H {
	return gin.H{
		"code":    "0",
		"message": "successful receive alert notification message!",
	}
}

// ClusterPodRequestScheduling cluster pod request scheduling.
func ClusterPodRequestScheduling(c *gin.Context) {
	pod := &v1.Pod{}
	if err := c.BindJSON(pod); err != nil {
		errStr := NewError(1, fmt.Sprintf("request parameter parsing failed, %v", err.Error()))
		c.JSON(http.StatusOK, errStr)
		return
	}
	client := clientset.InclusterClientset()
	nodeList, err := resources.GetNodes(client)
	if err != nil {
		errStr := fmt.Sprintf("cluster node get failed, %v", err)
		c.JSON(http.StatusOK, NewError(1, errStr))
		return
	}
	pods, err := resources.GetPods(client)
	if err != nil {
		errStr := fmt.Sprintf("cluster pods get failed, %v", err)
		c.JSON(http.StatusOK, NewError(1, errStr))
		return
	}
	status := calculate.ClusterPodRequestScheduling(pod, nodeList, pods)
	c.JSON(http.StatusOK, NewOKWithData(status))

}

//Response defines standard response structure
type Response struct {
	// Code is a manual reference, for troubleshooting
	Code int `json:"code"`
	// Message is human readable message for users go get know what's the matter
	Message string `json:"message"`
	// RawError is raw kubestar error, it's often a wrapper of errors
	RawError string `json:"raw"`
	// Data payload need to be handled by frontend
	Data interface{} `json:"data,omitempty"`
}

func NewOK() *Response {
	return &Response{
		Code:    0,
		Message: "execution succeed",
	}
}

func NewError(code int, errmsg string) *Response {
	return &Response{
		Code:     code,
		RawError: errmsg,
	}
}

func NewOKWithData(data interface{}) *Response {
	return &Response{
		Code:    0,
		Message: "execution succeed",
		Data:    data,
	}
}

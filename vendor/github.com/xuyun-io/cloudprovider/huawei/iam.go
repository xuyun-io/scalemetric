package huawei

import (
	"fmt"

	"github.com/imroc/req"
	"github.com/xuyun-io/cloudprovider/kubeconfig"
	"github.com/xuyun-io/cloudprovider/types"
)

// getToken return iam token.
// eg https://iam.cn-north-4.myhuaweicloud.com/v3/auth/tokens
func getToken(auth Auth, projectName string) (string, error) {
	url := generateIAMURL(auth.Base.Region)
	reqBody := generateIAMBody(auth.Base.DomainName, auth.Base.Account, auth.Base.Password, projectName)
	bj := req.BodyJSON(reqBody)
	resp, err := req.Post(url, bj)
	if err != nil {
		return "", err
	}
	if resp.Response().StatusCode >= 500 {
		return "", fmt.Errorf("huawei auth failed, http code != %d", resp.Response().StatusCode)
	}
	token := resp.Response().Header.Get("X-Subject-Token")
	if len(token) <= 0 {
		return "", fmt.Errorf("auth token not found")
	}
	return token, nil
}

func getSWRToken(auth Auth) (string, error) {
	return getToken(auth, auth.SWR.ProjectName)

}

func getCCEToken(auth Auth) (string, error) {
	return getToken(auth, auth.CCE.ProjectName)
}
func getKubeconfig(token, region, projectID, clusterID string) (types.Client, error) {
	url := generateKubeconfigURL(region, projectID, clusterID)
	header := generateTokenHeader(token)
	resp, err := req.Get(url, header)
	if err != nil {
		return nil, err
	}
	if resp.Response().StatusCode >= 500 {
		return nil, fmt.Errorf("hua wei auth failed, http code != %d", resp.Response().StatusCode)
	}

	kubecfg, err := resp.ToString()
	if err != nil {
		return nil, err
	}

	return kubeconfig.NewClientByRegion(region, kubecfg)

}

func generateKubeconfigURL(region, projectID, clusterID string) string {
	// https: //cce.cn-north-4.myhuaweicloud.com/api/v3/projects/项目ID/clusters/集群id/clustercert
	return fmt.Sprintf("https://cce.%s.myhuaweicloud.com/api/v3/projects/%s/clusters/%s/clustercert", region, projectID, clusterID)
}

func generateIAMURL(region string) string {
	return fmt.Sprintf("https://iam.%s.myhuaweicloud.com/v3/auth/tokens", region)
	// return "https://iam.cn-north-1.myhuaweicloud.com/v3/auth/tokens"
}

func generateClustersURL(region, projectID string) string {
	return fmt.Sprintf("https://cce.%s.myhuaweicloud.com/api/v3/projects/%s/clusters", region, projectID)
}

func generateClusterURL(region, projectID, clusterID string) string {
	return fmt.Sprintf("https://cce.%s.myhuaweicloud.com/api/v3/projects/%s/clusters/%s", region, projectID, clusterID)
}

func generateIAMBody(domainName, subAccount, password, projectName string) *IAMBody {
	return &IAMBody{
		ReqAuth: ReqAuth{
			Identity: Identity{
				Methods: []string{"password"},
				Password: Password{
					User: User{
						Name:     subAccount,
						Password: password,
						Domain: Domain{
							Name: domainName,
						},
					},
				},
			},
			Scope: Scope{
				Project: Project{
					Name: projectName,
				},
			},
		},
	}
}

// IAMBody define huawei iam body
type IAMBody struct {
	ReqAuth `json:"auth"`
}

// Project define project
type Project struct {
	Name string `json:"name"`
}

// Scope define scope
type Scope struct {
	Project Project `json:"project"`
}

// Domain define main account
type Domain struct {
	Name string `json:"name"`
}

// User define user
type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Domain   `json:"domain"`
}

// Password define password resource.
type Password struct {
	User `json:"user"`
}

// Identity define huawei Identity.
type Identity struct {
	Methods  []string `json:"methods"`
	Password `json:"password"`
}

// ReqAuth define huawei auth body
type ReqAuth struct {
	Identity `json:"identity"`
	Scope    Scope `json:"scope"`
}

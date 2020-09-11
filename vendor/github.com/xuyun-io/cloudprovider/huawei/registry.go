package huawei

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/xuyun-io/cloudprovider/types"
)

// RepoClientset return aws repo clientset
func (hw *Client) RepoClientset(repo ...types.RepoInterface) types.RepoInterface {
	if len(repo) > 0 {
		return repo[0]
	}
	// use default repo
	c, _ := newClient(hw.Auth)
	return c

}

// Repositories return huawei swr repos
func (hw *Client) Repositories(ctx context.Context) ([]*types.Repository, error) {
	return hw.getRepoListByNamespace(hw.Auth.Base.Region, hw.Auth.SWR.Namespace)
}

// TagList return repo tag list
func (hw *Client) TagList(ctx context.Context, repoName string) ([]*types.ImageDetail, error) {
	return hw.tags(ctx, hw.Auth.Base.Region, repoName)
}

// ImageDescribe return image describe
func (hw *Client) ImageDescribe(ctx context.Context, imageName string) (interface{}, error) {
	return nil, errors.New("swr image warehouse does not support external access")
}

// GenerateDockerLogin return generate docker login user password.
func (r SWRAuth) GenerateDockerLogin() (*DockerUser, error) {
	if len(r.ProjectName) <= 0 {
		return nil, fmt.Errorf("project name is not allowed to be empty, eg: cn-north-1")
	}
	if len(r.AccessKey) <= 0 || len(r.SecretAccessKey) <= 0 {
		return nil, errors.New("accesskey or secretaccesskey is not allowed to be empty")
	}

	pwd := generateDockerLoginPassword(r.AccessKey, r.SecretAccessKey)
	if len(pwd) <= 0 {
		return nil, errors.New("use openssl generate swr password failed")
	}
	user := fmt.Sprintf("%s@%s", r.ProjectName, r.AccessKey)
	host := repoHost(r.ProjectName)
	return &DockerUser{
		User:     user,
		Password: pwd,
		Host:     host,
	}, nil
}

func repoHost(projectName string) string {
	return fmt.Sprintf("swr.%s.myhuaweicloud.com", projectName)
}
func generateDockerLoginPassword(ak, sk string) string {
	// printf "$AK" | openssl dgst -binary -sha256 -hmac "$SK" | od -An -vtx1 | sed 's/[ \n]//g' | sed 'N;s/\n//'
	cmd1 := exec.Command("printf", fmt.Sprintf("%s", ak))
	cmd2 := exec.Command("openssl", []string{"dgst", "-binary", "-sha256", "-hmac", fmt.Sprintf("%s", sk)}...)
	cmd3 := exec.Command("od", []string{"-An", "-vtx1"}...)
	cmd4 := exec.Command("sed", []string{`s/[ \n]//g`}...)
	cmd5 := exec.Command("sed", []string{`N;s/\n//`}...)
	cmd5.Stdout = os.Stdout
	// cmd5.Stderr = os.Stderr
	pipe(cmd1, cmd2, cmd3, cmd4, cmd5)
	cmd1.Run()
	cmd2.Run()
	cmd3.Run()
	cmd4.Run()
	var out bytes.Buffer
	cmd5.Stdout = &out
	cmd5.Run()
	return out.String()
}

func pipe(cmds ...*exec.Cmd) {
	for i, cmd := range cmds {
		if i > 0 {
			out, _ := cmds[i-1].StdoutPipe()
			in, _ := cmd.StdinPipe()
			go func() {
				defer func() {
					out.Close()
					in.Close()
				}()
				io.Copy(in, out)
			}()
		}
	}
}

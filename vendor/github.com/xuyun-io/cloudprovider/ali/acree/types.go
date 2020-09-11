package acree

// RepoParams define repo params
type RepoEEParams struct {
	// required
	InstanceId string `json:"instance_id,omitempty"`
	// required, docker host
	Host string `json:"host,omitempty"`
	// required, Scheme = "https" // https | http
	Scheme string `json:"scheme,omitempty"`
	//
	PullImageUserName string `json:"pullImageUserName,omitempty"`
	PullImagePassword string `json:"pullImagePassword,omitempty"`
	// optional
	RepoNamespaceName string `json:"repoNamespaceName,omitempty"`
	// optional
	RepoName string `json:"repoName,omitempty"`
}

type Image struct {
	Status      string `json:"Status"`
	ImageCreate int64  `json:"ImageCreate"`
	ImageSize   int64  `json:"ImageSize"`
	Digest      string `json:"Digest"`
	ImageID     string `json:"ImageId"`
	ImageUpdate int64  `json:"ImageUpdate"`
	Tag         string `json:"Tag"`
}

type RepositoriesBody struct {
	Repositories []Repository `json:"Repositories"`
	IsSuccess    bool         `json:"IsSuccess"`
	TotalCount   int          `json:"TotalCount"`
	PageSize     int          `json:"PageSize"`
	RequestID    string       `json:"RequestId"`
	PageNo       int          `json:"PageNo"`
	Code         string       `json:"Code"`
	Message      string       `json:"Message"`
}

type Repository struct {
	RepoNamespaceName string `json:"RepoNamespaceName"`
	RepoBuildType     string `json:"RepoBuildType"`
	ModifiedTime      int64  `json:"ModifiedTime"`
	RepoType          string `json:"RepoType"`
	RepoStatus        string `json:"RepoStatus"`
	InstanceID        string `json:"InstanceId"`
	CreateTime        int64  `json:"CreateTime"`
	RepoName          string `json:"RepoName"`
	RepoID            string `json:"RepoId"`
}

type RepositoryDetail struct {
	IsSuccess         bool   `json:"IsSuccess"`
	RepoNamespaceName string `json:"RepoNamespaceName"`
	RepoBuildType     string `json:"RepoBuildType"`
	RepoType          string `json:"RepoType"`
	RepoStatus        string `json:"RepoStatus"`
	RequestID         string `json:"RequestId"`
	InstanceID        string `json:"InstanceId"`
	RepoName          string `json:"RepoName"`
	Summary           string `json:"Summary"`
	RepoID            string `json:"RepoId"`
	Code              string `json:"Code"`
	Detail            string `json:"Detail"`
	CreateTime        string `json:"CreateTime"`
	ModifiedTime      string `json:"ModifiedTime"`
	Message           string `json:"Message"`
}

type ImageTagList struct {
	IsSuccess  bool    `json:"IsSuccess"`
	TotalCount int     `json:"TotalCount"`
	PageSize   int     `json:"PageSize"`
	RequestID  string  `json:"RequestId"`
	Images     []Image `json:"Images"`
	PageNo     int     `json:"PageNo"`
	Message    string  `json:"Message"`
	Code       string  `json:"Code"`
}

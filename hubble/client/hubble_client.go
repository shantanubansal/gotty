package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/shantanubansal/gotty/hubble/util"
	"io/ioutil"
	nethttp "net/http"
	"time"
)

type HubbleClient struct {
	Endpoint string
}

type UserInfo struct {
	cli          HubbleClient
	Token        string
	ClusterUid   string
	UserUid      string
	ProjectUid   string
	HeartBeat    time.Time
	KubeConfig   string
	UserName     string
	Password     string
	creationTime string
	User         V1UserMe
}

var users = map[string]*UserInfo{}

func propertyCannotBeEmpty(property string) error {
	return fmt.Errorf("%s cannot be empty", property)
}
func User(params map[string][]string) (*UserInfo, error) {
	userinfo, err := getUserInfo(params)
	if err != nil {
		return nil, err
	}
	return userinfo.getInfoFromHubble()
}

func (u *UserInfo) getInfoFromHubble() (*UserInfo, error) {
	user, err := u.GetUserInfo()
	if err != nil {
		return nil, err
	}
	u.User = user
	userName := fmt.Sprintf("%s-%s", u.User.Spec.FirstName, u.creationTime)
	users[userName] = u
	kubeConfig, err := u.GetKubeConfig()
	if err != nil {
		return nil, err
	}
	u.UserName = userName
	if kubeConfig == "" {
		return nil, propertyCannotBeEmpty("KubeConfig")
	}
	u.KubeConfig = base64.StdEncoding.EncodeToString([]byte(kubeConfig))
	return u, nil
}

func getUserInfo(params map[string][]string) (*UserInfo, error) {
	token := getParam("Authorization", params)
	if token == "" {
		return nil, propertyCannotBeEmpty("Authorization Token")
	}
	userUid := getParam("userUid", params)
	if userUid == "" {
		return nil, propertyCannotBeEmpty("User Uid")
	}
	projectUid := getParam("projectUid", params)
	if projectUid == "" {
		return nil, propertyCannotBeEmpty("Project Uid")
	}
	spectroClusterUid := getParam("spectroClusterUid", params)
	if spectroClusterUid == "" {
		return nil, propertyCannotBeEmpty("SpectroCluster Uid")
	}
	endpoint := getParam("endpoint", params)
	if endpoint == "" {
		return nil, propertyCannotBeEmpty("Endpoint")
	}
	userinfo := &UserInfo{
		cli: HubbleClient{
			Endpoint: endpoint,
		},
		Token:        token,
		UserUid:      userUid,
		ProjectUid:   projectUid,
		ClusterUid:   spectroClusterUid,
		creationTime: fmt.Sprintf("%v", time.Now().Unix()),
	}
	return userinfo, nil
}

func getParam(paramConstant string, params map[string][]string) string {
	endpoints := params[paramConstant]
	endpoint := ""
	if len(endpoints) > 0 {
		endpoint = endpoints[0]
	}
	return endpoint
}

func GetHttpClientWithCert() *nethttp.Client {
	return util.GetHttpClientWithTls()
}

func (u *UserInfo) GetKubeConfig() (string, error) {
	cli := GetHttpClientWithCert()
	subPath := fmt.Sprintf("v1/spectroclusters/%s/assets/kubeconfig?Authorization=%v&ProjectUid=%s", u.ClusterUid, u.Token, u.ProjectUid)
	res, err := cli.Get(fmt.Sprintf("https://%v/%v", u.cli.Endpoint, subPath))
	if err != nil {
		return "", err
	}
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if err := u.unmarshallToError(respBody); err != nil {
		return "", err
	}
	return string(respBody), nil
}

func (u *UserInfo) GetUserInfo() (*V1UserMe, error) {
	var respMap *V1UserMe
	cli := GetHttpClientWithCert()
	subPath := fmt.Sprintf("v1/users/me?Authorization=%s", u.Token)
	res, err := cli.Get(fmt.Sprintf("https://%v/%v", u.cli.Endpoint, subPath))
	if err != nil {
		return respMap, err
	}
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return respMap, err
	}
	if err := u.unmarshallToError(respBody); err != nil {
		return nil, err
	}
	e := json.Unmarshal(respBody, respMap)
	if e != nil {
		return nil, e
	}
	return respMap, nil
}

func (u *UserInfo) unmarshallToError(respBody []byte) error {
	var respMap *V1Error
	e := json.Unmarshal(respBody, respMap)
	if e != nil {
		return nil
	}
	if respMap.Code != "" {
		return fmt.Errorf("%s: %s", respMap.Code, respMap.Message)
	}
	return nil
}

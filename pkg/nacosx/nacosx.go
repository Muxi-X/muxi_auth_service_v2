package nacosx

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func GetConfigFromNacos(confName string) (string, error) {
	conf, err := parseNacosDSN(confName)
	if err != nil {
		return "", err
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: conf.Server,
			Port:   conf.Port,
			Scheme: "http",
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         conf.Namespace,
		Username:            conf.User,
		Password:            conf.Password,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		CacheDir:            "./data/configCache",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return "", err
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: conf.DataId,
		Group:  conf.Group,
	})
	if err != nil {
		return "", err
	}
	return content, nil
}

// DSN 示例： localhost:8848?namespace=default&username=nacos&password=1234&group=QA&dataId=autossl-qiniuyun-qa
func parseNacosDSN(conf string) (*DSNConf, error) {
	var dsnConf DSNConf
	var _ = godotenv.Load()

	dsn := os.Getenv(conf)
	if dsn == "" {
		return nil, fmt.Errorf("env %s not set", conf)
	}

	parts := strings.SplitN(dsn, "?", 2)
	host := parts[0]
	params := url.Values{}

	if len(parts) == 2 {
		params, _ = url.ParseQuery(parts[1])
	}

	hostParts := strings.Split(host, ":")
	dsnConf.Server = hostParts[0]
	if len(hostParts) > 1 {
		p, _ := strconv.Atoi(hostParts[1])
		dsnConf.Port = uint64(p)
	} else {
		dsnConf.Port = 8848
	}

	dsnConf.Namespace = params.Get("namespace")
	if dsnConf.Namespace == "" {
		dsnConf.Namespace = "public"
	}

	dsnConf.User = params.Get("username")
	dsnConf.Password = params.Get("password")
	dsnConf.Group = params.Get("group")
	dsnConf.DataId = params.Get("dataId")
	return &dsnConf, nil
}

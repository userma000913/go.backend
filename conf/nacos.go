package conf

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
)

type Nacos struct {
	Sc []constant.ServerConfig
	Cc constant.ClientConfig
}

func InitNacos(sc []constant.ServerConfig, cc constant.ClientConfig) *Nacos {
	return &Nacos{
		Sc: sc,
		Cc: cc,
	}
}

// nacos做配置中心
func (n *Nacos) GetNacosConfigClient() config_client.IConfigClient {
	// 创建nacos config客户端
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &n.Cc,
			ServerConfigs: n.Sc,
		},
	)

	if err != nil {
		log.Println("nacos client创建失败 。。。", err.Error())
		return nil
	}

	return client

}

func (n *Nacos) GetNacosConfigContent() (string, error) {
	content, err := n.GetNacosConfigClient().
		GetConfig(vo.ConfigParam{DataId: "go.backend", Group: "DEFAULT_GROUP"})
	if err != nil {
		log.Println("配置获取失败 。。。", err.Error())
		return "", err
	}
	return content, nil
}

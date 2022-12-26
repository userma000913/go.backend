package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"time"
)

type AppConfig struct {
	Server *Server      `json:"*service" yaml:"server"`
	Mysql  *MySQLConfig `yaml:"mysql" json:"mysql"`
	Redis  *RedisConfig `yaml:"redis" json:"redis"`
}

type Server struct {
	Env  string `yaml:"env" json:"env"`
	Name string `json:"name"`
	Port int    `json:"port"`
}

// MySQLConfig 定义mysql配置文件的结构体
type MySQLConfig struct {
	Host         string `mapstructure:"host" yaml:"host"`
	User         string `mapstructure:"user" yaml:"user"`
	Password     string `mapstructure:"password" yaml:"password"`
	DB           string `mapstructure:"dbname" yaml:"dbname"`
	Port         int    `mapstructure:"port" yaml:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
}

// RedisConfig 定义redis配置文件的结构体
type RedisConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Password string `mapstructure:"password" yaml:"password"`
	Port     int    `mapstructure:"port" yaml:"port"`
	DB       int    `mapstructure:"db" yaml:"db"`
	PoolSize int    `mapstructure:"pool_size" yaml:"pool_size"`
}

type EsConfig struct {
	Host string `mapstructure:"host" yaml:"host"`
}

func InitConfigWithNacos() *AppConfig {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := constant.ClientConfig{
		NamespaceId:         "public",
		Username:            "nacos",
		Password:            "nacos",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	n := InitNacos(sc, cc)
	// 获取nacos配置
	config, err := n.GetNacosConfigContent()
	if err != nil {
		log.Println(err)
		return nil
	}

	var c *AppConfig
	if err := yaml.Unmarshal([]byte(config), &c); err != nil {
		log.Printf("Unmarshal is err;err=%s", err)
		return nil
	}
	fmt.Println("conf ok")
	return c
}

func InitConfigWithFile(path string) *AppConfig {

	// 读取文件
	viper.SetConfigFile(path)
	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return nil
	}
	var c *AppConfig
	if err := viper.Unmarshal(&c); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		return nil
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(&c); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})

	return c
}

// InitConfigWithRemote todo viper远程读取etcd/consul
func InitConfigWithRemote(path string) {
	// alternatively, you can create a new viper instance.
	var runtimeViper = viper.New()

	runtimeViper.AddRemoteProvider("etcd", "server://127.0.0.1:4001", path)
	runtimeViper.SetConfigType("yaml") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"

	// read from remote conf the first time.
	_ = runtimeViper.ReadRemoteConfig()

	// unmarshal conf
	var c *AppConfig
	runtimeViper.Unmarshal(&c)

	// open a goroutine to watch remote changes forever
	// 必须开一个协程后台监听，且必须是死循环监听，因为WatchRemoteConfig并没有死循环 而文件监听的WatchConfig有！！！
	go func() {
		for {
			time.Sleep(time.Second * 5) // delay after each request

			// currently, only tested with etcd support
			err := runtimeViper.WatchRemoteConfig()
			if err != nil {
				log.Printf("unable to read remote conf: %v", err)
				continue
			}

			// unmarshal new conf into our runtime conf struct. you can also use channel
			// to implement a signal to notify the system of the changes
			runtimeViper.Unmarshal(&c)
		}
	}()
}

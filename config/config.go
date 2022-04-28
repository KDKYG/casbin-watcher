package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	WatcherConfig *GlobalConfig
)

var (
	confPath string
	TemplateDir string
)

const (
	ConfigFile = "./config.yaml"
)

type GlobalConfig struct {
	ListenPort         string                   `yaml:"ListenPort"` //监听的端口
	RedisConf          RedisConf                `yaml:"RedisConf"` //缓存
}

type RedisConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Instance int    `yaml:"instance"`
}

func Init() {
	flag.StringVar(&confPath, "f", ConfigFile, "Usage:./out -f xx")
	flag.StringVar(&TemplateDir, "d", "./site/*", "Usage:./out -d xx")
	//解析上面定义的标签
	flag.Parse()

	WatcherConfig = &GlobalConfig{}
	fileBuffer, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatalln("config ReadFile error:" + confPath)
	}
	err = yaml.Unmarshal(fileBuffer, WatcherConfig)
	if err != nil {
		log.Fatalln("config Unmarshal failed")
	}
	log.Printf("parse config :%v", *WatcherConfig)
}

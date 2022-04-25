package main

import (
	"fmt"
	"os"

	"github.com/chenjiayao/sidergo"
	"github.com/chenjiayao/sidergo/cluster"
	"github.com/chenjiayao/sidergo/config"
	"github.com/chenjiayao/sidergo/interface/server"
	"github.com/sirupsen/logrus"

	"github.com/chenjiayao/sidergo/redis"
	_ "github.com/chenjiayao/sidergo/redis/datatype"
)

func main() {

	banner := `
           __        __                                         
          |  \      |  \                                        
  _______  \$$  ____| $$  ______    ______    ______    ______  
 /       \|  \ /      $$ /      \  /      \  /      \  /      \ 
|  $$$$$$$| $$|  $$$$$$$|  $$$$$$\|  $$$$$$\|  $$$$$$\|  $$$$$$\
 \$$    \ | $$| $$  | $$| $$    $$| $$   \$$| $$  | $$| $$  | $$
 _\$$$$$$\| $$| $$__| $$| $$$$$$$$| $$      | $$__| $$| $$__/ $$
|       $$| $$ \$$    $$ \$$     \| $$       \$$    $$ \$$    $$
 \$$$$$$$  \$$  \$$$$$$$  \$$$$$$$ \$$       _\$$$$$$$  \$$$$$$ 
                                            |  \__| $$          
                                             \$$    $$          
                                              \$$$$$$           

	`
	fmt.Println(banner)

	customFormatter := new(logrus.TextFormatter)
	customFormatter.FullTimestamp = true                    // 显示完整时间
	customFormatter.TimestampFormat = "2006-01-02 15:04:05" // 时间格式
	customFormatter.DisableTimestamp = false                // 禁止显示时间
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(customFormatter)

	configFile := os.Getenv("REDIS_CONFIG")
	if configFile == "" {
		config.LoadDefaultConfig()
	} else {
		config.LoadConfig(configFile)
	}

	s := makeServer()

	sidergo.ListenAndServe(s)
}

func makeServer() server.Server {
	if config.Config.EnableCluster {
		return cluster.MakeCluster()
	} else {
		return redis.MakeRedisServer()
	}
}

package config

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type ServerConfig struct {
	Bind        string
	Port        int
	Databases   int
	RequirePass string
	Appendonly  bool //是否开启 aof

	AppendFilename string //aof 文件名称
}

// golang 的 code style：如果一个变量是全局单例，直接设为全局变量
var Config *ServerConfig

func LoadConfig(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	Config = parseConfig(file)
}

func LoadDefaultConfig() {
	Config = &ServerConfig{
		Bind:           "127.0.01",
		Port:           3101,
		Databases:      16,
		RequirePass:    "123456",
		Appendonly:     false,
		AppendFilename: "",
	}
}

func parseConfig(reader io.Reader) *ServerConfig {
	c := &ServerConfig{}
	configMap := loadConfig(reader)

	//使用反射来解析 ServerConfig
	t := reflect.TypeOf(c)
	v := reflect.ValueOf(c)

	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		fieldVal := v.Elem().Field(i)

		configName, ok := field.Tag.Lookup("config")
		if !ok {
			configName = field.Name
		}
		configValue, ok := configMap[strings.ToLower(configName)]

		if ok {
			switch field.Type.Kind() {
			case reflect.String:
				fieldVal.SetString(configValue)
			case reflect.Int:
				intValue, err := strconv.ParseInt(configValue, 10, 64)
				if err == nil {
					fieldVal.SetInt(intValue)
				}
			case reflect.Bool:
				boolValue := "yes" == configValue
				fieldVal.SetBool(boolValue)
			case reflect.Slice:
				if field.Type.Elem().Kind() == reflect.String {
					s := strings.Split(configValue, ",")
					fieldVal.Set(reflect.ValueOf(s))
				}
			}
		}

	}
	return c
}

func loadConfig(reader io.Reader) map[string]string {
	m := make(map[string]string)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue //空行 或者注释直接跳过
		}

		configNameIndex := strings.IndexAny(line, " ")
		if configNameIndex < 0 || configNameIndex == len(line)-1 {
			continue
		}

		configName := line[0:configNameIndex]
		configValue := strings.Trim(line[configNameIndex+1:], " ")
		m[strings.ToLower(configName)] = configValue
	}
	return m
}

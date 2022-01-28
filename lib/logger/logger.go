package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/chenjiayao/goredistraning/lib/file"
)

type logLevel int

const (
	DEBUG logLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var (
	logPrefixs = []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}
)

var (
	logger *log.Logger
	mu     sync.Mutex
)

func Setting() {

	fileFullPath := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
	file, err := file.OpenFile("./logs", fileFullPath)
	if err != nil {
		log.Fatalf("打开文件失败: %s", err)
	}

	mw := io.MultiWriter(file, os.Stdout)
	logger = log.New(mw, "", log.LstdFlags)
}

func setPrefix(level logLevel) {
	prefix := fmt.Sprintf("[%s]", logPrefixs[level])
	logger.SetPrefix(prefix)
}

func Debug(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	setPrefix(DEBUG)
	logger.Println(v...)
}

func Info(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	setPrefix(INFO)
	logger.Println(v...)
}

func Error(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	setPrefix(ERROR)
	logger.Println(v...)
}

func Fatal(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	setPrefix(FATAL)
	logger.Println(v...)
}

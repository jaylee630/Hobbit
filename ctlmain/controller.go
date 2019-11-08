package ctlmain

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/jaylee630/Hobbit/api"
	"github.com/jaylee630/Hobbit/config"
	"github.com/jaylee630/Hobbit/utils/log"
)

var (
	buildDate      string
	buildVersion   string
	buildGoVersion string

	configPath string
)

func usage() {

	fmt.Printf(
		"Usage:\n  %s [--version] [--config=value]\n", os.Args[0])
	flag.PrintDefaults()
}

// parseFlags parse command line flags
// and notice caller to run server or not by return value
func parseFlags() bool {

	showVersion := false

	flag.BoolVar(&showVersion, "version", false, "show version info")
	flag.StringVar(&configPath, "config", "./etc/golang-admin-basic.yml", "set config file path")

	flag.Usage = usage
	flag.Parse()

	if showVersion {
		fmt.Printf("%s\n  Version: %s\n  Build Date: %s\n  Go Version: %s\n",
			path.Base(os.Args[0]), buildVersion, buildDate, buildGoVersion)
		return false
	}

	return true

}

type Controller struct {
	*log.Logger
	config  *config.BasicConfig
	server  *api.Server
	db      *gorm.DB
	Engine  *gin.Engine
	Context context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func (ctl *Controller) Run() {

	// start http server
	ctl.server.Start()
}

func (ctl *Controller) Stop() {

	ctl.Info("Shutdown Server ...")

	ctl.cancel()
	ctl.server.Stop()
	ctl.wg.Wait()

	ctl.Info("Server Controller stopped")

}

func (ctl *Controller) Init() {

	ctl.Context, ctl.cancel = context.WithCancel(context.TODO())

	// init config
	cfg, err := config.LoadConfig(configPath, "")
	if err != nil {
		panic(err.Error())
	}
	ctl.config = cfg

	if err := ctl.initLogger(); err != nil {
		panic(err.Error())
	}

	if err := ctl.initDB(); err != nil {
		panic(err.Error())
	}

	if err := ctl.initEngine(); err != nil {
		panic(err.Error())
	}

	if err := ctl.initModule(); err != nil {
		panic(err.Error())
	}

	if err := ctl.initAPIServer(); err != nil {
		panic(err.Error())
	}
}

func (ctl *Controller) initLogger() error {

	// init logger
	logger, err := log.NewLogger(&ctl.config.Logger)
	if err != nil {
		return err
	}
	ctl.Logger = logger
	return nil
}

func (ctl *Controller) initDB() error {
	// init db
	var err error
	dbInfo := ctl.config.GetDBSource()
	if _, ok := dbInfo["mysql"]; !ok {
		ctl.Fatalf("fail to get db config for mysql: no such config")
	}
	ctl.Info("connect to MySQL, dialect: mysql, data_source" + dbInfo["mysql"])
	ctl.db, err = NewGORM(ctl.Logger, ctl.config.Logger.Level, "mysql", dbInfo["mysql"])
	if err != nil {
		return err
	}

	return nil
}

func (ctl *Controller) initAPIServer() error {

	ctl.server = api.NewHTTPServer(ctl.Logger, ctl.config.Host, ctl.config.Port,
		ctl.config.AdminPort, ctl.config.Logger.Level == "debug", ctl.config.Logger.Level)

	ctl.server.Init(ctl.Context, ctl.Engine)

	return nil

}

func (ctl *Controller) initModule() error {
	// init module here
	return nil
}

func (ctl *Controller) initEngine() error {

	if ctl.config.Logger.Level == "debug" {
		ctl.Engine = NewGinEngine(os.Stdout, ctl.config.Logger.Level)
	} else {
		ctl.Engine = NewGinEngine(ctl.Logger, ctl.config.Logger.Level)
	}

	return nil
}

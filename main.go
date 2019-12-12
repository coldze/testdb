package main

import (
	"flag"
	"fmt"
	"github.com/coldze/testdb/logic/handlers"
	"github.com/coldze/testdb/logic/handlers/cabs"
	"github.com/coldze/testdb/logic/handlers/caches"
	"github.com/coldze/testdb/logic/sources"
	"github.com/coldze/testdb/logic/sources/stores"
	"github.com/coldze/testdb/logic/sources/wraps/mysql"
	"github.com/coldze/testdb/logic/sources/wraps/redis"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/coldze/testdb/logs"
	"github.com/coldze/testdb/utils"

	_ "github.com/go-sql-driver/mysql"
)

const (
	HEALTH_CHECK_PATH = "/ping"
	API_VERSION       = "v1"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(fmt.Sprintf("Health check at: %v\n", time.Now().UTC())))
	if r.Body == nil {
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()
	_, _ = ioutil.ReadAll(r.Body)
}

func buildRoutes(logger logs.Logger, source sources.Source, cache sources.Store) (http.Handler, error) {
	getHandler, err := cabs.NewGetHandler(handlers.NewDefaultLoggerFactory(logs.NewPrefixedLogger(logger, "[GET cab]")), source)
	if err != nil {
		return nil, err
	}
	deleteCacheHandler, err := caches.NewDeleteHandler(handlers.NewDefaultLoggerFactory(logs.NewPrefixedLogger(logger, "[DELETE cache]")), cache)
	if err != nil {
		return nil, err
	}
	wipeCacheHandler := caches.NewWipeHandler(handlers.NewDefaultLoggerFactory(logs.NewPrefixedLogger(logger, "[WIPE cache]")), cache)
	router := mux.NewRouter()
	router.Path(HEALTH_CHECK_PATH).HandlerFunc(healthCheck)

	sr := router.PathPrefix(fmt.Sprintf("/%s", API_VERSION)).Subrouter()
	sr.HandleFunc("/cab", getHandler).Methods(http.MethodGet)
	sr.HandleFunc("/caches", wipeCacheHandler).Methods(http.MethodDelete)
	sr.HandleFunc("/cache", deleteCacheHandler).Methods(http.MethodDelete)

	return router, nil
}

func newMainFunc(cfg *appCfg) utils.MainFunc {
	return func(logger logs.Logger, stop <-chan struct{}) int {

		dbWrap, err := mysql.NewMysqlDbWrap(cfg.GetMysqlURI(), mysql.NewScanner)
		if err != nil {
			panic(err)
		}
		defer dbWrap.Close()
		dSource := stores.NewDbDataSource(dbWrap, mysql.BuildQuery)

		redisWrap, err := redis.NewRedis(cfg.GetRedisOptions())
		if err != nil {
			panic(err)
		}
		defer redisWrap.Close()

		cache := stores.NewCache(redisWrap, cfg.GetCacheTtl())

		source := sources.NewSourceWithCache(dSource, cache)
		router, err := buildRoutes(logger, source, cache)
		if err != nil {
			panic(err)
		}

		bind := cfg.GetBind()
		srv, err := utils.NewService(bind, router)
		if err != nil {
			logger.Errorf("Failed to start service. Error: %v", err)
			return 1
		}
		defer func() {
			cErr := srv.Stop()
			if cErr != nil {
				logger.Errorf("Failed to stop service: %+v", cErr)
			}
		}()
		logger.Infof("Ready. Listening at '%s'", bind)
		<-stop
		return 0
	}
}

func main() {
	configPath := flag.String("config", "./config.json", "service's configuration in JSON format")
	redisPwd := flag.String("redispwd", "", "Redis password")
	mysqlPwd := flag.String("mysqlpwd", "", "MySQL password")

	flag.Parse()
	logger := logs.NewStdLogger()
	logger.Infof("Starting...")
	cfg, err := getConfig(*configPath, *redisPwd, *mysqlPwd)
	if err != nil {
		logger.Errorf("Failed to load config. Error: %v", err)
		return
	}
	utils.Run(cfg.GetAppTimeout(), newMainFunc(cfg), logger)
	logger.Infof("Done")
}

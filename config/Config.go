package config

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/url"
	"os"
	"path"
	"strconv"
)

var (
	APPVIPER *viper.Viper
	DB       *gorm.DB
	CLIENT   []*ethclient.Client
)

func init() {
	APPVIPER = initAppConfig()
	DB = initDB()
	CLIENT = initClient()
}

func initAppConfig() *viper.Viper {
	workDir, _ := os.Getwd()
	appViper := viper.New()
	appViper.SetConfigName("application")
	appViper.SetConfigType("yml")
	appViper.AddConfigPath(path.Join(workDir, "config"))
	err := appViper.ReadInConfig()
	if err != nil {

	}
	return appViper
}

func initDB() *gorm.DB {
	host := APPVIPER.GetString("database.host")
	port := APPVIPER.GetString("database.port")
	database := APPVIPER.GetString("database.databaseName")
	username := APPVIPER.GetString("database.username")
	password := APPVIPER.GetString("database.password")
	charset := APPVIPER.GetString("database.charset")
	loc := APPVIPER.GetString("database.loc")

	sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc),
	)

	db, err := gorm.Open(mysql.Open(sqlStr), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
	if err != nil {
		panic("connected error" + err.Error())
	} else {
		log.Infoln("connect db success")
	}
	return db
}

func initClient() []*ethclient.Client {
	var clients []*ethclient.Client
	for i := 1; i < 102; i++ {
		url := APPVIPER.GetString("infura.url" + strconv.Itoa(i))
		client, err := ethclient.Dial(url)
		if err != nil {
			log.Error("client faild:", i, err)
		} else {
			clients = append(clients, client)
		}

	}
	log.Infoln("connect client success")
	return clients
}

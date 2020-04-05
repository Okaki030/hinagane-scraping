package config

import (
	"database/sql"
	"log"

	// _ "github.vom/liv/pq"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DBConfig DBConfig
}

type DBConfig struct {
	Port     string
	Host     string
	User     string
	Password string
}

var DB *sql.DB

func Connect() (*sql.DB, error) {

	if DB == nil {

		var err error

		// TODO:configファイルを使う
		// 実行の前に全てコンテナ側にコピーしてあげればいける？
		// env := os.Getenv("CONFIG_ENV")

		// var configName string
		// if env == "" {
		// 	configName = "config.toml"
		// } else {
		// 	configName = "/go/src/config." + env + ".toml"
		// }

		// var config Config
		// _, err = toml.DecodeFile(configName, &config)
		// if err != nil {
		// 	log.Println(err)
		// 	return nil, err
		// }

		// datasourceName := config.DBConfig.User + ":" + config.DBConfig.Password + "@tcp(" + config.DBConfig.Host + ":" + config.DBConfig.Port + ")/hinagane_db?parseTime=true&loc=Asia%2FTokyo"

		datasourceName := "root:root@tcp(mysql:3306)/hinagane_db?parseTime=true&loc=Asia%2FTokyo"
		DB, err = sql.Open("mysql", datasourceName)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return DB, nil
}

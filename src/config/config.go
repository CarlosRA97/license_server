package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
)

type Config struct {
	Port string
	Dsn string
}

func GetEnvOrDefault(key string, d string) string {
	env, ok := os.LookupEnv(key)
	if (!ok && env == "") {
		env = d
	}
	return env
}

func getSystemTimeZone() string {
	buf := new(bytes.Buffer)
	tz, ok := os.LookupEnv("TZ")
	if !ok {
		p, err := os.Readlink("/etc/localtime")
		if err != nil {
			log.Fatalln(err)
		}
		dir, lname := path.Split(p)
		_, fname := path.Split(dir[:len(dir)-1])
		if fname == "zoneinfo" {
			fmt.Fprint(buf, lname) // prints the timezone string, e.g. Japan
		} else {
			fmt.Fprint(buf, fname + "/" + lname) // prints the timezone string, e.g. Asia/Tokyo
		}
	} else {
		fmt.Fprint(buf, tz)
	}
	return buf.String()
}

func GetEnv() Config {
	port := GetEnvOrDefault("PORT", ":3000")

	dsn := "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"
	postgresHost 		:= GetEnvOrDefault("POSTGRES_HOST", "127.0.0.1")
	postgresUser 		:= GetEnvOrDefault("POSTGRES_USER", "gorm")
	postgresPassword 	:= GetEnvOrDefault("POSTGRES_PASSWORD", "gorm")
	postgresDB 			:= GetEnvOrDefault("POSTGRES_DB", "gorm")
	postgresPort 		:= GetEnvOrDefault("POSTGRES_PORT", "5432")
	postgresSslMode 	:= GetEnvOrDefault("POSTGRES_SSLMODE", "disable")
	postgresTimeZone 	:= GetEnvOrDefault("POSTGRES_TZ", getSystemTimeZone())
	dsn = fmt.Sprintf(dsn, postgresHost, postgresUser, postgresPassword, postgresDB, postgresPort, postgresSslMode, postgresTimeZone)
	return Config{
		Port: port,
		Dsn: dsn,
	}
}
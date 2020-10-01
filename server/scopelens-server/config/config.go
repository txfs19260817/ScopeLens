package config

import (
	"fmt"
	"time"

	"github.com/go-ini/ini"
)

var (
	Mode     string
	App      *app
	Server   *server
	Database *database
	Jwt      *jwt
	Path     *path
	Aws      *aws
)

type app struct {
	PageSize int
}

type server struct {
	Url          string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	HttpsCrt     string
	HttpsKey     string
}

type database struct {
	Type     string
	AtlasURI string
	DBName   string
}

type jwt struct {
	SigningKey      string
	ReCaptchaSecret string
}

type path struct {
	ImageSavePath string
	SpritePath    string
	LogPath       string
}

type aws struct {
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
	TeamPath  string
}

func init() {
	// Load config from .ini file.
	cfgPath := "config/config.ini"
	cfg, err := ini.Load(cfgPath)
	if err != nil {
		panic(fmt.Sprintf("Fail to parse '%v': %v", cfgPath, err))
	}

	// Map from `cfg` to struct and save as global variable
	// Mode
	Mode = cfg.Section("").Key("Mode").MustString("debug")

	// [app]
	App = new(app)
	if err := cfg.Section("app").MapTo(App); err != nil {
		panic(fmt.Sprintf("Fail to map [app] to App: %v", err))
	}

	// [server]
	Server = new(server)
	if err := cfg.Section("server").MapTo(Server); err != nil {
		panic(fmt.Sprintf("Fail to map [server] to Server: %v", err))
	}

	// [database]
	Database = new(database)
	if err := cfg.Section("database").MapTo(Database); err != nil {
		panic(fmt.Sprintf("Fail to map [database] to Database: %v", err))
	}

	// [jwt]
	Jwt = new(jwt)
	if err := cfg.Section("jwt").MapTo(Jwt); err != nil {
		panic(fmt.Sprintf("Fail to map [jwt] to Jwt: %v", err))
	}

	// [jwt]
	Path = new(path)
	if err := cfg.Section("path").MapTo(Path); err != nil {
		panic(fmt.Sprintf("Fail to map [path] to Path: %v", err))
	}

	// [aws]
	Aws = new(aws)
	if err := cfg.Section("aws").MapTo(Aws); err != nil {
		panic(fmt.Sprintf("Fail to map [aws] to Aws: %v", err))
	}
}

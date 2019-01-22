package config

import (
	"github.com/pelletier/go-toml"
	"log"
	"os"
)

// Database ...
type Database struct {
	Type     string `toml:"type"`
	Addr     string `toml:"addr"`
	Port     string `toml:"port"`
	Password string `toml:"password"`
	Username string `toml:"username"`
}

// Callback ...
type Callback struct {
	Type     string `toml:"type"`
	BackType string `toml:"back_type"`
	BackAddr string `toml:"back_addr"`
}

// Media ...
type Media struct {
	Upload      string `toml:"upload"`        //上传路径
	Transfer    string `toml:"transfer"`      //转换路径
	M3U8        string `toml:"m3u8"`          //m3u8文件名
	KeyURL      string `toml:"key_url"`       //default url
	KeyDest     string `toml:"key_dest"`      //key 文件输出目录
	KeyFile     string `toml:"key_file"`      //key文件名
	KeyInfoFile string `toml:"key_info_file"` //keyFile文件名
}

// IPFS ...
type IPFS struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

// GRPC ...
type GRPC struct {
	Enable bool   `toml:"enable"`
	Type   string `toml:"type"`
	Path   string `toml:"path"`
	Port   string `toml:"port"`
}

// REST ...
type REST struct {
	Enable  bool   `toml:"enable"`
	Type    string `toml:"type"`
	Path    string `toml:"path"`
	BackURL string `toml:"back_url"`
	Port    string `toml:"port"`
}

// Queue ...
type Queue struct {
	Type     string `toml:"type"`
	HostPort string `toml:"host_port"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

// DownloadInfo ...
type DownloadInfo struct {
	DirPath  string `toml:"dir_path"`
	PartSize int64  `toml:"part_size"`
}

// OSS ...
type OSS struct {
	Endpoint        string `toml:"endpoint"`
	AccessKeyID     string `toml:"access_key_id"`
	AccessKeySecret string `toml:"access_key_secret"`
	BucketName      string `toml:"bucket_name"`
}

// Configure ...
type Configure struct {
	Database Database `toml:"database"`
	Media    Media    `toml:"media"`
	Queue    Queue    `toml:"queue"`
	GRPC     GRPC     `toml:"grpc"`
	REST     REST     `toml:"rest"`
	IPFS     IPFS     `toml:"ipfs"`
	OSS      []OSS    `toml:"oss"`

	Callback Callback `toml:"callback"`
}

var config *Configure

// NullString ...
func NullString(s string) bool {
	return s == ""
}

// Initialize ...
func Initialize(filePath ...string) error {
	if filePath == nil {
		filePath = []string{"config.toml"}
	}

	cfg := LoadConfig(filePath[0])

	cfg.Media.Upload = DefaultString(cfg.Media.Upload, "upload")
	if !IsExists(cfg.Media.Upload) {
		err := os.Mkdir(cfg.Media.Upload, os.ModePerm)
		if err != nil {
			return err
		}
	}

	cfg.Media.Transfer = DefaultString(cfg.Media.Transfer, "transfer")
	if !IsExists(cfg.Media.Transfer) {
		err := os.Mkdir(cfg.Media.Transfer, os.ModePerm)
		if err != nil {
			return err
		}
	}
	cfg.Media.KeyDest = DefaultString(cfg.Media.KeyDest, "keydest")
	if !IsExists(cfg.Media.KeyDest) {
		err := os.Mkdir(cfg.Media.KeyDest, os.ModePerm)
		if err != nil {
			return err
		}
	}

	config = cfg

	return nil
}

// IsExists ...
func IsExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Panicln(err)
	}
	return true
}

// LoadConfig ...
func LoadConfig(filePath string) *Configure {
	var cfg Configure
	openFile, err := os.OpenFile(filePath, os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	decoder := toml.NewDecoder(openFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("config: %+v", cfg)
	return &cfg
}

// Config ...
func Config() *Configure {
	if config == nil {
		panic("nil config")
	}
	return config
}

// DefaultString ...
func DefaultString(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

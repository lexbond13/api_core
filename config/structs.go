package config

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	Params        *Params        `json:"params"`
	DB            *DB            `json:"db"`
	Logger        *Logger        `json:"logger"`
	Server        *Server        `json:"server"`
	Notifications *Notifications `json:"notifications"`
	CDNStorage    *CDNStorage    `json:"cdn_storage"`
	Cache         *Cache         `json:"cache"`
}

type Params struct {
	AppName                string       `json:"app_name"`
	AppDevMode             bool         `json:"app_dev_mode"`
	OpenAuthMode           bool         `json:"open_auth_mode"`
	DomainApp              string       `json:"domain_app"`
	ApiURL                 string       `json:"api_url"`
	FrontendURL            string       `json:"frontend_url"`
	AuthSecret             string       `json:"auth_secret"`
	NodeName               string       `json:"node_name"`
	MaxCookieLifeTimeHours int64        `json:"max_cookie_life_time_hours"`
	UploadFiles            *UploadFiles `json:"upload_files"`
}

type DB struct {
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	PoolSize int    `json:"pool_size" validators:"min=0"`
}

type Logger struct {
	FileConfig   *FileConfig
	SentryConfig *SentryConfig
}

// FileConfig
type FileConfig struct {
	Path   string `json:"path"`
	Perm   string `json:"perm"`
	Level  string `json:"level"`
	Format string `json:"format,omitempty"`
}

// SentryConfig
type SentryConfig struct {
	DSN string `json:"dsn"`
}

type Notifications struct {
	Email *Email `json:"email"`
}

type Email struct {
	SMTPSender    *SMTPSender    `json:"smtp_sender"`
	MailGunSender *MailGunSender `json:"mail_gun_sender"`
}

type SMTPSender struct {
	Host        string `json:"host"`
	Port        int64  `json:"port"`
	AuthEmail   string `json:"auth_email"`
	SenderEmail string `json:"sender_email"`
	SenderName  string `json:"sender_name"`
	Password    string `json:"password"`
}

type MailGunSender struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UploadFiles struct {
	AllowSize       int64    `json:"allow_size"`
	AllowExtensions []string `json:"allow_extensions"`
}

type CDNStorage struct {
	SelCDN *SelCDN
}

type SelCDN struct {
	URL           string `json:"url"`
	ContainerName string `json:"container_name"`
	AuthURL       string `json:"auth_url"`
	User          string `json:"access_token"`
	Key           string `json:"key"`
	DevProdMode   bool   `json:"dev_prod_mode"`
}

type Cache struct {
	Redis *Redis
}

type Redis struct {
	Network  string `json:"network"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Database int64  `json:"database"`
	PoolSize int64  `json:"pool_size" validate:"min=0"`
	IsUsed   bool   `json:"is_used"`
}

type Server struct {
	Host string
	Port int64
}

// Init configs
func Init() (*AppConfig, error) {

	appConfig := &AppConfig{}

	// For default, the configuration is taken from the environment variables of the operating system.
	// If they are not, then they are loaded from the config .env.prod file in config path.
	configFile := "./config/.env"

	devMode := false
	_, ok := os.LookupEnv("APP_DEV_MODE")
	if ok {
		devMode = getEnvAsBool("APP_DEV_MODE", false)
		if devMode {
			configFile = "./config/.env.loc"
		}
	}

	if err := godotenv.Load(configFile); err != nil {
		return nil, errors.New("Load from env file... No .env file found")
	}

	appConfig.Params = &Params{
		AppName:                os.Getenv("APP_NAME"),
		DomainApp:              os.Getenv("DOMAIN_APP"),
		ApiURL:                 os.Getenv("API_URL"),
		FrontendURL:            os.Getenv("FRONTEND_URL"),
		NodeName:               os.Getenv("APP_NAME"),
		OpenAuthMode:           getEnvAsBool("OPEN_AUTH_MODE", false),
		AppDevMode:             getEnvAsBool("APP_DEV_MODE", false),
		AuthSecret:             os.Getenv("AUTH_SECRET"),
		MaxCookieLifeTimeHours: getEnvAsInt("MAX_COOKIE_LIFE_TIME_HOURS", 24),
		UploadFiles: &UploadFiles{
			AllowSize:       getEnvAsInt("UPLOAD_FILES_ALLOW_SIZE", 0),
			AllowExtensions: strings.Split(os.Getenv("UPLOAD_FILES_ALLOW_EXTENSIONS"), ","),
		},
	}

	// Check for set auth secret
	if appConfig.Params.AuthSecret == "" {
		return nil, errors.New("Auth secret not set! You must set it in environment param: AUTH_SECRET in OS")
	}

	appConfig.DB = &DB{
		Host:     os.Getenv("DB_HOST"),
		Port:     getEnvAsInt("DB_PORT", 0),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	}

	appConfig.Server = &Server{
		Host: os.Getenv("SERVER_HOST"),
		Port: getEnvAsInt("SERVER_PORT", 0),
	}

	appConfig.Logger = &Logger{
		FileConfig: &FileConfig{
			Path:  os.Getenv("LOGGER_FILE_PATH"),
			Perm:  os.Getenv("LOGGER_FILE_PERM"),
			Level: os.Getenv("LOGGER_FILE_LEVEL"),
		},
		SentryConfig: &SentryConfig{
			DSN: os.Getenv("SENTRY_DSN"),
		},
	}

	appConfig.Notifications = &Notifications{
		Email: &Email{
			SMTPSender: &SMTPSender{
				Host:        os.Getenv("SMTP_HOST"),
				Port:        getEnvAsInt("SMTP_PORT", 0),
				AuthEmail:   os.Getenv("SMTP_AUTH_EMAIL"),
				SenderEmail: os.Getenv("SMTP_SENDER_EMAIL"),
				SenderName:  os.Getenv("SMTP_SENDER_NAME"),
				Password:    os.Getenv("SMTP_PASSWORD"),
			},
			MailGunSender: &MailGunSender{
				Host:     os.Getenv("MAILGUN_HOST"),
				Username: os.Getenv("MAILGUN_USERNAME"),
				Password: os.Getenv("MAILGUN_PASSWORD"),
			},
		},
	}

	appConfig.CDNStorage = &CDNStorage{
		SelCDN: &SelCDN{
			URL:           os.Getenv("CDN_URL"),
			ContainerName: os.Getenv("CDN_CONTAINER"),
			AuthURL:       os.Getenv("CDN_AUTH_URL"),
			User:          os.Getenv("CDN_USER"),
			Key:           os.Getenv("CDN_KEY"),
			DevProdMode:   getEnvAsBool("CDN_DEV_PROD_MODE", false),
		},
	}

	appConfig.Cache = &Cache{
		Redis: &Redis{
			Network:  os.Getenv("REDIS_NETWORK"),
			Address:  os.Getenv("REDIS_ADDRESS"),
			Password: os.Getenv("REDIS_PASSWORD"),
			Database: getEnvAsInt("REDIS_DATABASE", 0),
			PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 0),
			IsUsed:   getEnvAsBool("REDIS_IS_USED", false),
		},
	}

	return appConfig, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int64) int64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

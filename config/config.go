package config

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Kafka     KafkaConfig     `mapstructure:"kafka"`
	Cassandra CassandraConfig `mapstructure:"cassandra"`
	Logger    LoggerConfig    `mapstructure:"logger"`
	Jaeger    JaegerConfig    `mapstructure:"jaeger"`
}

type ServerConfig struct {
	Host                string `mapstructure:"host"`
	Port                string `mapstructure:"port"`
	GrpcHost            string `mapstructure:"grpcHost"`
	GrpcPort            string `mapstructure:"grpcPort"`
	Debug               bool   `mapstructure:"debug"`
	ReadTimeout         int    `mapstructure:"readTimeout"`
	WriteTimeout        int    `mapstructure:"writeTimeout"`
	ContextTimeout      int    `mapstructure:"contextTimeout"`
	Timezone            string `mapstructure:"timezone"`
	AccessJwtSecret     string `mapstructure:"accessJwtSecret"`
	AccessJwtExpireTime int    `mapstructure:"AccessJwtExpireTime"`
	Location            *time.Location
	Hostname            string `mapstructure:"hostname"`
	InstanceID          int64
	InstanceIP          string `mapstructure:"instanceIP"`
}

type DatabaseConfig struct {
	Adapter         string `mapstructure:"adapter"`
	Host            string `mapstructure:"host"`
	Username        string `mapstructure:"username"`
	Db              string `mapstructure:"db"`
	Password        string `mapstructure:"password"`
	Port            int    `mapstructure:"port"`
	MaxConns        int    `mapstructure:"maxConns"`
	MaxLiftimeConns int    `mapstructure:"maxLiftimeConns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
}

type LoggerConfig struct {
	Development       bool   `mapstructure:"development"`
	DisableCaller     bool   `mapstructure:"disableCaller"`
	DisableStacktrace bool   `mapstructure:"disableStacktrace"`
	Encoding          string `mapstructure:"encoding"`
	Level             string `mapstructure:"level"`
	Filename          string `mapstructure:"filename"`
	FileMaxSize       int    `mapstructure:"fileMaxSize"`
	FileMaxAge        int    `mapstructure:"fileMaxAge"`
	FileMaxBackups    int    `mapstructure:"fileMaxBackups"`
	FileIsCompress    bool   `mapstructure:"fileIsCompress"`
}

type CassandraConfig struct {
	Host         []string `mapstructure:"host"`
	Keyspace     string   `mapstructure:"keyspace"`
	Consistency  string   `mapstructure:"consistency"`
	ProtoVersion int      `mapstructure:"protoVersion"`
	Username     string   `mapstructure:"username"`
	Password     string   `mapstructure:"password"`
}

type kafkaProducer struct {
	Brokers    string `mapstructure:"broker"`
	Acks       string `mapstructure:"acks"`
	Idepotence bool   `mapstructure:"idepotence"`
}
type kafkaConsumer struct {
	Brokers     string `mapstructure:"broker"`
	GroupID     string `mapstructure:"groupID"`
	Timeout     int    `mapstructure:"timeout"`
	OffsetReset string `mapstructure:"offsetReset"`
	AutoOffset  bool   `mapstructure:"autoOffset"`
}
type KafkaConfig struct {
	Producer kafkaProducer `mapstructure:"producer"`
	Consumer kafkaConsumer `mapstructure:"consumer"`
}

type JaegerConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"serviceName"`
	LogSpans    bool   `mapstructure:"logSpans"`
}

func GetConf() (*Config, error) {
	re := regexp.MustCompile(`\$\{([^{}]+)\}`)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if re.Match([]byte(value)) {
			env := string(re.ReplaceAll([]byte(value), []byte("$1")))
			viper.Set(k, os.Getenv(env))
		}

	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation(config.Server.Timezone)
	if err != nil {
		return nil, err
	}
	config.Server.Location = loc

	hostnameArr := strings.Split(config.Server.Hostname, "-")
	InstanceID, err := strconv.Atoi(hostnameArr[len(hostnameArr)-1])
	if err != nil {
		InstanceID = 1
	}
	config.Server.InstanceID = int64(InstanceID)

	return &config, nil
}

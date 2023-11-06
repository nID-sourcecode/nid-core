// Package database provides utility functionality on a PostgreSQL databse
package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	// Postgres driver
	_ "github.com/lib/pq"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
)

// database constants
const (
	dbConnection    = "dbname=%s user=%s password=%s port=%d host=%s application_name='%s' sslmode=%s"
	dbConnectionLog = "Start connection to dbname=%s user=%s port=%d host=%s"
	createExtension = "CREATE EXTENSION IF NOT EXISTS \"%s\";"
	maxTimeout      = 120

	defaultTimeout       = 60
	defaultPort          = 5432
	defaultRetryDuration = 500 * time.Millisecond

	// TestModeOff Different test modes for db
	TestModeOff         TestMode = 1
	TestModeDropTables  TestMode = 2
	TestModeNoDropTable TestMode = 3

	DefaultOpenConns            int           = 15
	DefaultIdleConns            int           = 15
	DefaultConnLifetimeDuration time.Duration = 30 * time.Minute
)

// TestMode Enum
type TestMode int64

// TestModeDesc description of different test modes
// nolint: gochecknoglobals
var TestModeDesc = map[TestMode]string{
	TestModeOff:         "Test mode off",
	TestModeDropTables:  "Test mode with dropping tables",
	TestModeNoDropTable: "Test mode without dropping tables",
}

// CreateDefaultExtensionConfig creates a default db extension config
func CreateDefaultExtensionConfig() []string {
	return []string{"uuid-ossp", "pg_trgm"}
}

// DBConfig config for database connection
type DBConfig struct {
	Host            string
	Port            int
	User            string
	Pass            string
	DBName          string
	RetryOnFailure  bool
	LogMode         bool
	TestMode        TestMode
	TimeOut         int
	AutoMigrate     bool
	Extensions      []string
	ApplicationName string
	SSLMode         string
}

// CreateTestDBConfig creates desired test DB configuration
func CreateTestDBConfig() *DBConfig {
	return &DBConfig{
		RetryOnFailure: true,
		TestMode:       TestModeNoDropTable,
		Port:           defaultPort,
		Pass:           "postgres",
		TimeOut:        defaultTimeout,
		AutoMigrate:    true,
		Extensions:     CreateDefaultExtensionConfig(),
		SSLMode:        "disable",
	}
}

// setDefaultOptions sets max connection lifetime below istio 2 hour TCP connection limit and defaults to 15 open and 15 idle connections
func setDefaultOptions(db *gorm.DB) *gorm.DB {
	db.DB().SetMaxOpenConns(DefaultOpenConns)
	db.DB().SetMaxIdleConns(DefaultIdleConns)
	db.DB().SetConnMaxLifetime(DefaultConnLifetimeDuration)
	return db
}

// MustConnectTest must connection for running tests
func MustConnectTest(dbName string, models []interface{}) *gorm.DB {
	db, err := ConnectTest(dbName, models)
	if err != nil {
		panic(err)
	}
	return setDefaultOptions(db)
}

// ConnectTest connection for running tests
func ConnectTest(dbName string, models []interface{}) (*gorm.DB, error) {
	config := CreateTestDBConfig()
	config.DBName = dbName
	client := NewClientFromConfig(config)
	return client.InitDB(models)
}

// MustConnectCustom must connect custom database connection with specified config
func MustConnectCustom(config *DBConfig, models []interface{}) *gorm.DB {
	db, err := ConnectCustom(config, models)
	if err != nil {
		panic(err)
	}
	return setDefaultOptions(db)
}

// MustConnectCustomWithCustomLogger must connect custom database connection with specified config and custom logger utility
func MustConnectCustomWithCustomLogger(config *DBConfig, models []interface{}, loggerUtility log.LoggerUtility) *gorm.DB {
	log.SetLoggerUtility(loggerUtility)
	db, err := ConnectCustom(config, models)
	if err != nil {
		panic(err)
	}
	return db
}

// ConnectCustom create custom database connection with specfied config
func ConnectCustom(config *DBConfig, models []interface{}) (*gorm.DB, error) {
	client := NewClientFromConfig(config)
	return client.InitDB(models)
}

// NewClientFromConfig create postgres client from config
func NewClientFromConfig(config *DBConfig) PostgresClient {
	client := PostgresClient{
		DBConfig:          config,
		connectionRetries: 0,
	}
	// Config prio: config > env > defaults
	client.evalEnvironment()
	client.ensureDefaults()
	return client
}

// CreateTestDatabase Creates a database in postgres.
func CreateTestDatabase(name string) error {
	conf := CreateTestDBConfig()
	conf.Host = "localhost"
	conf.User = "postgres"
	conf.TestMode = TestModeOff

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Pass)

	var db *sql.DB
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s_test;", name))
	if err != nil && !strings.Contains(err.Error(), fmt.Sprintf("database \"%s_test\" already exists", name)) {
		return err
	}

	return nil
}

// PostgresClient client for connecting to postgres
type PostgresClient struct {
	*DBConfig
	DB                *gorm.DB
	connectionRetries int
}

// InitDB connect, migrate & create extensions
func (c *PostgresClient) InitDB(models []interface{}) (*gorm.DB, error) {
	c.LogConfig()
	err := c.Connect()
	if err != nil {
		return nil, err
	}
	c.DB.LogMode(c.LogMode)
	c.CreateDBExtensions()
	c.Migrate(models)
	c.DB = setDefaultOptions(c.DB)
	return c.DB, nil
}

// Connect try to create a connection
func (c *PostgresClient) Connect() error {
	applicationName := "default-service"
	if c.ApplicationName != "" {
		applicationName = c.ApplicationName
	}
	var err error
	url := fmt.Sprintf(dbConnection, c.DBName, c.User, c.Pass, c.Port, c.Host, applicationName, c.SSLMode)
	logURL := fmt.Sprintf(dbConnectionLog, c.DBName, c.User, c.Port, c.Host)
	log.Infof(logURL)
	c.DB, err = gorm.Open("postgres", url)
	if c.RetryOnFailure {
		for err != nil && c.connectionRetries < c.TimeOut {
			c.connectionRetries++
			log.Warnf("unable to connect to %s, error: %s", logURL, err.Error())
			time.Sleep(defaultRetryDuration)
			c.DB, err = gorm.Open("postgres", url)
		}
	}
	return err
}

// Migrate migrate database according to client config
func (c *PostgresClient) Migrate(models []interface{}) {
	if c.TestMode == TestModeDropTables {
		// If running tests always drop db up front.
		// Because in case tests fail, they might mess with data for other tests
		c.DB.DropTableIfExists(models...)
	}
	if c.AutoMigrate {
		c.DB.AutoMigrate(models...)
	}
}

// CreateDBExtensions create database extensions
func (c *PostgresClient) CreateDBExtensions() {
	for _, extension := range c.DBConfig.Extensions {
		c.DB.Exec(fmt.Sprintf(createExtension, extension))
	}
}

// ensureDefaults ensure clients defaults are set
func (c *PostgresClient) ensureDefaults() {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.User == "" {
		c.User = "postgres"
	}
	if c.TestMode != TestModeOff {
		c.DBName += "_test"
	}
	if c.Port == 0 {
		c.Port = defaultPort
	}
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}
	if c.RetryOnFailure {
		if c.TimeOut <= 0 {
			c.TimeOut = maxTimeout
		}
		if c.TimeOut > maxTimeout {
			c.TimeOut = maxTimeout
		}
	}
}

// LogConfig log current client config
func (c *PostgresClient) LogConfig() {
	log.Infof("Configuration:")
	log.Infof("DB HOST: %s", c.Host)
	log.Infof("DB Name: %s", c.DBName)
	log.Infof("DB User: %s", c.User)
	log.Infof("DB Port: %d", c.Port)
	log.Infof("SSLMode: %s", c.SSLMode)
	log.Infof("TestMode: %s", TestModeDesc[c.TestMode])
	log.Infof("Retry on failure: %v", c.RetryOnFailure)
	if c.RetryOnFailure {
		log.Infof("Timeout: %d", c.TimeOut)
	}
}

func (c *PostgresClient) evalEnvironment() {
	host := os.Getenv("PG_HOST")
	user := os.Getenv("PG_USERNAME")
	pass := os.Getenv("PG_PASSWORD")
	port := os.Getenv("PG_PORT")
	sslmode := os.Getenv("PG_SSLMODE")
	dbName := os.Getenv("PG_DBNAME")
	log.Infof("HOST: %s", host)
	if c.DBConfig.Host == "" && host != "" {
		c.Host = host
	}
	if c.DBConfig.User == "" && user != "" {
		c.User = user
	}
	if c.DBConfig.Pass == "" && pass != "" {
		c.Pass = pass
	}
	if c.DBConfig.DBName == "" && dbName != "" {
		c.DBName = dbName
	}
	if c.DBConfig.SSLMode == "" && sslmode != "" {
		c.SSLMode = sslmode
	}
	if c.DBConfig.Port == 0 && port != "" {
		pgPort, err := strconv.Atoi(port)
		if err != nil {
			c.Port = defaultPort
		}
		c.Port = pgPort
	}
}

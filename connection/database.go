package connection

import (
	"chat/globals"
	"chat/utils"
	"crypto/tls"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

var DB *sql.DB

func InitMySQLSafe() *sql.DB {
	ConnectDatabase()

	// using DB as a global variable to point to the latest db connection
	MysqlWorker(DB)
	return DB
}

func getConnPostgres() *sql.DB {
	sslMode := "disable"
	if viper.GetBool("postgres.tls") {
		sslMode = "require"
	}

	postgresUrl := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		utils.GetStringConfs("postgres.database", "postgres.db"),
		sslMode,
	)

	db, err := sql.Open("postgres", postgresUrl)

	if pingErr := db.Ping(); err != nil || pingErr != nil {
		errMsg := utils.Multi(err != nil, utils.GetError(err), utils.GetError(pingErr))
		globals.Warn(
			fmt.Sprintf("[connection] failed to connect to postgres server: %s (message: %s), will retry in 5 seconds",
				viper.GetString("postgres.host"), errMsg,
			),
		)

		utils.Sleep(5000)
		db.Close()

		return getConnPostgres()
	}

	globals.Debug(fmt.Sprintf("[connection] connected to postgres server (host: %s)", viper.GetString("postgres.host")))
	return db
}

func getConnMysql() *sql.DB {
	if viper.GetString("mysql.host") == "" {
		globals.SqliteEngine = true
		globals.Warn("[connection] mysql host is not set, using sqlite (~/db/chatnio.db)")
		db, err := sql.Open("sqlite3", utils.FileSafe("./db/chatnio.db"))
		if err != nil {
			panic(err)
		}

		return db
	}

	mysqlUrl := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		utils.GetStringConfs("mysql.database", "mysql.db"),
	)
	if viper.GetBool("mysql.tls") {
		mysql.RegisterTLSConfig("tls", &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: viper.GetString("mysql.host"),
		})

		mysqlUrl += "?tls=tls"
	}

	// connect to MySQL
	db, err := sql.Open("mysql", mysqlUrl)

	if pingErr := db.Ping(); err != nil || pingErr != nil {
		errMsg := utils.Multi[string](err != nil, utils.GetError(err), utils.GetError(pingErr)) // err.Error() may contain nil pointer
		globals.Warn(
			fmt.Sprintf("[connection] failed to connect to mysql server: %s (message: %s), will retry in 5 seconds",
				viper.GetString("mysql.host"), errMsg,
			),
		)

		utils.Sleep(5000)
		db.Close()

		return getConnMysql()
	}

	globals.Debug(fmt.Sprintf("[connection] connected to mysql server (host: %s)", viper.GetString("mysql.host")))
	return db
}

func ConnectDatabase() *sql.DB {
	var db *sql.DB

	if viper.GetString("database.driver") == "postgres" {
		globals.PostgresEngine = true
		db = getConnPostgres()
	} else {
		// default to MySQL
		db = getConnMysql()
	}

	db.SetMaxOpenConns(512)
	db.SetMaxIdleConns(64)

	CreateUserTable(db)
	CreateConversationTable(db)
	CreateMaskTable(db)
	CreateSharingTable(db)
	CreatePackageTable(db)
	CreateQuotaTable(db)
	CreateSubscriptionTable(db)
	CreateApiKeyTable(db)
	CreateInvitationTable(db)
	CreateRedeemTable(db)
	CreateBroadcastTable(db)

	// nil when postgres
	if err := doMigration(db); err != nil {
		fmt.Println(fmt.Sprintf("migration error: %s", err))
	}

	DB = db

	return db
}

func InitRootUser(db *sql.DB) {
	// create root user if totally empty
	var count int
	err := globals.QueryRowDb(db, "SELECT COUNT(*) FROM auth").Scan(&count)
	if err != nil {
		globals.Warn(fmt.Sprintf("[service] failed to query user count: %s", err.Error()))
		return
	}

	if count == 0 {
		globals.Debug("[service] no user found, creating root user (username: root, password: chatnio123456, email: root@example.com)")
		_, err := globals.ExecDb(db, `
			INSERT INTO auth (username, password, email, is_admin, bind_id, token)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "root", utils.Sha2Encrypt("chatnio123456"), "root@example.com", true, 0, "root")
		if err != nil {
			globals.Warn(fmt.Sprintf("[service] failed to create root user: %s", err.Error()))
		}
	} else {
		globals.Debug(fmt.Sprintf("[service] %d user(s) found, skip creating root user", count))
	}
}

func CreateUserTable(db *sql.DB) {
	_, err := globals.ExecDb(db, `
		CREATE TABLE IF NOT EXISTS auth (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  bind_id INT UNIQUE,
		  username VARCHAR(24) UNIQUE,
		  token VARCHAR(255) NOT NULL,
		  email VARCHAR(255) UNIQUE,
		  password VARCHAR(64) NOT NULL,
		  is_admin BOOLEAN DEFAULT FALSE,
		  is_banned BOOLEAN DEFAULT FALSE
		);
	`)
	if err != nil {
		fmt.Println(err)
	}

	InitRootUser(db)
}

func CreatePackageTable(db *sql.DB) {
	_, err := globals.ExecDb(db, `
		CREATE TABLE IF NOT EXISTS package (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  user_id INT,
		  type VARCHAR(255),
		  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  FOREIGN KEY (user_id) REFERENCES auth(id),
		  UNIQUE KEY (user_id, type)
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateQuotaTable(db *sql.DB) {
	_, err := globals.ExecDb(db, `
		CREATE TABLE IF NOT EXISTS quota (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  user_id INT UNIQUE,
		  quota DECIMAL(24, 6),
		  used DECIMAL(24, 6),
		  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  FOREIGN KEY (user_id) REFERENCES auth(id)
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateConversationTable(db *sql.DB) {
	_, err := globals.ExecDb(db, `
		CREATE TABLE IF NOT EXISTS conversation (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  user_id INT,
		  conversation_id INT,
		  conversation_name VARCHAR(255),
		  data MEDIUMTEXT,
		  model VARCHAR(255) NOT NULL DEFAULT 'gpt-3.5-turbo-0613',
		  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  UNIQUE KEY (user_id, conversation_id)
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateMaskTable(db *sql.DB) {
	_, err := globals.ExecDb(db, `
		CREATE TABLE IF NOT EXISTS mask (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  user_id INT,
		  avatar VARCHAR(255),
		  name VARCHAR(255),
		  description TEXT,
		  context TEXT,
		  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  FOREIGN KEY (user_id) REFERENCES auth(id)
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateSharingTable(db *sql.DB) {
	// refs is an array of message id, separated by comma (-1 means all messages)
	_, err := globals.ExecDb(db, `
		CREATE TABLE IF NOT EXISTS sharing (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  hash CHAR(32) UNIQUE,
		  user_id INT,
		  conversation_id INT,
		  refs TEXT,
		  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  FOREIGN KEY (user_id) REFERENCES auth(id)
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateSubscriptionTable(db *sql.DB) {
	_, err := globals.ExecDb(db, `
		CREATE TABLE IF NOT EXISTS subscription (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  level INT DEFAULT 1,
		  user_id INT UNIQUE,
		  expired_at DATETIME,
		  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  total_month INT DEFAULT 0,
		  enterprise BOOLEAN DEFAULT FALSE,
		  FOREIGN KEY (user_id) REFERENCES auth(id)
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateApiKeyTable(db *sql.DB) {
	_, err := globals.ExecDb(db, `
		CREATE TABLE IF NOT EXISTS apikey (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  user_id INT UNIQUE,
		  api_key VARCHAR(255) UNIQUE,
		  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  FOREIGN KEY (user_id) REFERENCES auth(id)
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateInvitationTable(db *sql.DB) {
	_, err := globals.ExecDb(db, `
		CREATE TABLE IF NOT EXISTS invitation (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  code VARCHAR(255) UNIQUE,
		  quota DECIMAL(16, 4),
		  type VARCHAR(255),
		  used BOOLEAN DEFAULT FALSE,
		  used_id INT,
		  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  UNIQUE KEY (used_id, type),
		  FOREIGN KEY (used_id) REFERENCES auth(id)
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateRedeemTable(db *sql.DB) {
	_, err := globals.ExecDb(db, `
		CREATE TABLE IF NOT EXISTS redeem (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  code VARCHAR(255) UNIQUE,
		  quota DECIMAL(16, 4),
		  used BOOLEAN DEFAULT FALSE,
		  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateBroadcastTable(db *sql.DB) {
	_, err := globals.ExecDb(db, `
		CREATE TABLE IF NOT EXISTS broadcast (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  poster_id INT,
		  content TEXT,
		  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  FOREIGN KEY (poster_id) REFERENCES auth(id)
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
}

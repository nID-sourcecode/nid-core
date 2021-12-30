# Installation

Install using: `go get -u lab.weave.nl/weave/utilities/database/v2@v2.0.0`

# DB usage

This readme explains the usage of the database package for database connections with Postgres.

## Predefined connect

For example in your `models` package add the following function:

```golang
// InitDatabase inits the database models
func InitDatabase(logMode bool, testMode bool) *gorm.DB {
	m := []interface{}{
		Model1{},
		Model2{},
	}
	var db *gorm.DB
	if testMode {
		db = database.MustConnectTest("idtrust", m)
	} else {
		db = database.MustConnectIstio("idtrust", m)
	}
	db.DB().SetMaxOpenConns(50)

	return db
}
```

## Custom connect

You can also specify a custom connection config. The example below specifies:
1. No table migrations
2. Postgis extension installation if it not exists
3. Enable db logmode
```golang
// InitDatabase inits the database models
func InitDatabase(logMode bool, testMode bool) *gorm.DB {
	m := []interface{}{
		Model1{},
		Model2{},
    }
	var db *gorm.DB
	if testMode {
		db = database.MustConnectTest("idtrust", m)
	} else {
		if testMode {
		    db = database.MustConnectTest("idtrust", m)
        } else {
            db = database.MustConnectCustom(database.DBConfig{
                DBName:      "idtrust",
                Extensions:  []string{"postgis"},
                LogMode:     true,
                AutoMigrate: false,
            }, m)
        }
	}
	db.DB().SetMaxOpenConns(50)

	return db
}
```
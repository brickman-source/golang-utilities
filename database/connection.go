package database

import (
	"database/sql"
	"fmt"
	"github.com/brickman-source/golang-utilities/errors"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/brickman-source/golang-utilities/retry"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/luna-duclos/instrumentedsql/opentracing"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SQLConnection struct {
	db            *sqlx.DB
	DSN           string
	driverName    string
	driverPackage string
	options
}

func cleanURLQuery(in url.Values) (out url.Values) {
	out, _ = url.ParseQuery(in.Encode())
	out.Del("max_conns")
	out.Del("max_idle_conns")
	out.Del("max_conn_lifetime")
	out.Del("parseTime")
	return out
}

func (c *SQLConnection) GetDatabase() *sqlx.DB {
	if err := retry.Run(time.Second*15, time.Minute, func() (err error) {
		c.db, err = c.getDatabase()
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil
	}

	return c.db
}

func classifyDSN(dsn string) string {
	scheme := strings.Split(dsn, "://")[0]
	parts := strings.Split(dsn, "@")
	host := parts[len(parts)-1]
	return fmt.Sprintf("%s://*:*@%s", scheme, host)
}
func (c *SQLConnection) getDatabase() (*sqlx.DB, error) {
	if c.db != nil {
		return c.db, nil
	}

	driverName, driverPackage, err := c.registerDriver()
	if err != nil {
		return nil, errors.WithCode(err, "internal_server_error", "could not register driver")
	}

	dsn, err := connectionString(c.DSN)
	if err != nil {
		return nil, err
	}

	// fulfil classifyDSN()'s format, which must starts with schema://
	dsnToClassify := dsn
	if driverName == "mysql" {
		dsnToClassify = "mysql://" + strings.TrimPrefix(dsnToClassify, "mysql://")
	}
	classifiedDSN := classifyDSN(dsnToClassify)
	log.With(log.Fields{"dsn": classifiedDSN}).Infof("Establishing connection with SQL database backend")

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		log.With(log.Fields{"dsn": classifiedDSN}).Errorf("Unable to open SQL connection")
		return nil, errors.WithCode(err, "internal_server_error", "could not open SQL connection")
	}

	c.db = sqlx.NewDb(db, driverPackage) // This must be clean.Scheme otherwise things like `Rebind()` won't work
	if err := c.db.Ping(); err != nil {
		log.With(log.Fields{"dsn": classifiedDSN}).Errorf("Unable to ping SQL database backend")
		return nil, errors.WithCode(err, "internal_server_error", "could not ping SQL connection")
	}

	log.With(log.Fields{"dsn": classifiedDSN}).Infof("Successfully connected to SQL database backend")

	_, query, err := parseQuery(c.DSN)
	if err != nil {
		return nil, err
	}

	maxConns := maxParallelism() * 2
	if v := query.Get("max_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.With(log.Fields{"err": err.Error()}).Warnf(`Query parameter "max_conns" value %v could not be parsed to int, falling back to default value %d`, v, maxConns)
		} else {
			maxConns = int(s)
		}
	}

	maxIdleConns := maxParallelism()
	if v := query.Get("max_idle_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.With(log.Fields{"err": err.Error()}).Warnf(`Query parameter "max_idle_conns" value %v could not be parsed to int, falling back to default value %d`, v, maxIdleConns)
		} else {
			maxIdleConns = int(s)
		}
	}

	maxConnLifetime := time.Duration(0)
	if v := query.Get("max_conn_lifetime"); v != "" {
		s, err := time.ParseDuration(v)
		if err != nil {
			log.With(log.Fields{"err": err.Error()}).Warnf(`Query parameter "max_conn_lifetime" value %v could not be parsed to int, falling back to default value %d`, v, maxConnLifetime)
		} else {
			maxConnLifetime = s
		}
	}

	c.db.SetMaxOpenConns(maxConns)
	c.db.SetMaxIdleConns(maxIdleConns)
	c.db.SetConnMaxLifetime(maxConnLifetime)

	return c.db, nil
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
func parseQuery(dsn string) (clean string, query url.Values, err error) {
	query = url.Values{}
	parts := strings.Split(dsn, "?")
	clean = parts[0]
	if len(parts) == 2 {
		if query, err = url.ParseQuery(parts[1]); err != nil {
			return "", query, errors.WithStack(err)
		}
	}
	return
}
func connectionString(dsn string) (string, error) {
	dsn, query, err := parseQuery(dsn)
	if err != nil {
		return "", err
	}

	query = cleanURLQuery(query)
	if strings.HasPrefix(dsn, "mysql://") {
		query.Set("parseTime", "true")
		dsn = strings.TrimPrefix(dsn, "mysql://")
	}
	return dsn + "?" + query.Encode(), nil
}
func (c *SQLConnection) registerDriver() (string, string, error) {
	if c.driverName != "" {
		return c.driverName, c.driverPackage, nil
	}

	scheme := strings.Split(c.DSN, "://")[0]
	driverName := scheme
	driverPackage := scheme

	if c.UseTracedDriver {
		driverName = "instrumented-sql-driver"
		if len(c.options.forcedDriverName) > 0 {
			driverName = c.options.forcedDriverName
		}

		tracingOpts := []instrumentedsql.Opt{instrumentedsql.WithTracer(opentracing.NewTracer(c.AllowRoot))}
		if c.OmitArgs {
			tracingOpts = append(tracingOpts, instrumentedsql.WithOmitArgs())
		}
		drivers := sql.Drivers()
		foundIndex := sort.SearchStrings(drivers, driverName)

		if !(foundIndex < len(drivers) && drivers[foundIndex] == driverName) {
			switch scheme {
			case "mysql":
				sql.Register(driverName, instrumentedsql.WrapDriver(mysql.MySQLDriver{}, tracingOpts...))
			case "postgres":
				sql.Register(driverName, instrumentedsql.WrapDriver(&pq.Driver{}, tracingOpts...))
			default:
				return "", "", fmt.Errorf("unsupported scheme (%s) in DSN", scheme)
			}
		}
	}

	c.driverName = driverName
	c.driverPackage = driverPackage

	return driverName, driverPackage, nil
}

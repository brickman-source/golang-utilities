package database

import (
	"database/sql"
	"github.com/pkg/errors"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/brickman-source/golang-utilities/retry"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SQLConnection struct {
	db         *sqlx.DB
	URL        *url.URL
	driverName string
}

func cleanURLQuery(c *url.URL) *url.URL {
	cleanurl := new(url.URL)
	*cleanurl = *c

	q := cleanurl.Query()
	q.Del("max_conns")
	q.Del("max_idle_conns")
	q.Del("max_conn_lifetime")
	q.Del("parseTime")

	cleanurl.RawQuery = q.Encode()
	return cleanurl
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

func (c *SQLConnection) getDatabase() (*sqlx.DB, error) {
	if c.db != nil {
		return c.db, nil
	}

	var err error
	var registeredDriver string

	clean := cleanURLQuery(c.URL)
	if registeredDriver, err = c.registerDriver(); err != nil {
		return nil, errors.Wrap(err, "could not register driver")
	}

	log.Infof(nil, "Connecting with %s", c.URL.Scheme+"://*:*@"+c.URL.Host+c.URL.Path+"?"+clean.RawQuery)
	u := connectionString(clean)

	db, err := sql.Open(registeredDriver, u)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open SQL connection")
	}

	c.db = sqlx.NewDb(db, clean.Scheme)
	if err := c.db.Ping(); err != nil {
		return nil, errors.Wrapf(err, "could not ping SQL connection")
	}

	log.Infof(nil, "Connected to SQL!")

	maxConns := maxParallelism() * 2
	if v := c.URL.Query().Get("max_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Warnf(nil, "max_conns value %s could not be parsed to int: %s", v, err)
		} else {
			maxConns = int(s)
		}
	}

	maxIdleConns := maxParallelism()
	if v := c.URL.Query().Get("max_idle_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Warnf(nil, "max_idle_conns value %s could not be parsed to int: %s", v, err)
		} else {
			maxIdleConns = int(s)
		}
	}

	maxConnLifetime := time.Duration(0)
	if v := c.URL.Query().Get("max_conn_lifetime"); v != "" {
		s, err := time.ParseDuration(v)
		if err != nil {
			log.Warnf(nil, "max_conn_lifetime value %s could not be parsed to int: %s", v, err)
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

func connectionString(clean *url.URL) string {
	if clean.Scheme == "mysql" {
		q := clean.Query()
		q.Set("parseTime", "true")
		clean.RawQuery = q.Encode()
	}

	username := clean.User.Username()
	userinfo := username
	password, hasPassword := clean.User.Password()
	if hasPassword {
		userinfo = url.QueryEscape(userinfo) + ":" + url.QueryEscape(password)
	}
	clean.User = nil
	u := clean.String()
	clean.User = url.UserPassword(username, password)

	if strings.HasPrefix(u, clean.Scheme+"://") {
		u = strings.Replace(u, clean.Scheme+"://", clean.Scheme+"://"+userinfo+"@", 1)
	}
	if clean.Scheme == "mysql" {
		u = strings.Replace(u, "mysql://", "", -1)
	}
	return u
}

func (c *SQLConnection) registerDriver() (string, error) {
	if c.driverName != "" {
		return c.driverName, nil
	}

	driverName := c.URL.Scheme

	c.driverName = driverName
	return driverName, nil
}

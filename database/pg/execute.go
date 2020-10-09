package pg

import (
	"database/sql"
	"github.com/brickman-source/golang-utilities/errors"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/jmoiron/sqlx"
)

func GetItems(items []interface{}, countSqlClause, dataSqlClause string, params []interface{}, pageSize, pageIndex int, tx *sqlx.Tx) (totalCount int, err error) {

	if pageSize < 1 {
		pageSize = 10
	}
	if pageIndex < 0 {
		pageIndex = 0
	}

	if tx != nil {
		err = tx.Get(
			&totalCount,
			tx.Rebind(countSqlClause),
			params...)

		if err == nil {
			err = tx.Select(
				&items,
				tx.Rebind(dataSqlClause),
				params...)
			if err != nil {
				log.Errorf( "get items err:%s", err.Error())
			}
		} else {
			log.Errorf( "get count err:%s", err.Error())
		}
		if err == sql.ErrNoRows {
			return 0, err
		} else if err != nil {
			return 0, err
		}
	} else {
		return 0, errors.New("db error")
	}

	return totalCount, nil
}

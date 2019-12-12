package mysql

import (
	"fmt"
	"strings"

	"github.com/coldze/testdb/logic/sources/wraps"
	"github.com/coldze/testdb/logic/structs"
)

const (
	query_format      = "SELECT medallion, count(*) as trips FROM cabs_db.cab_trip_data WHERE pickup_datetime BETWEEN ? AND ? AND medallion IN ('%s')  GROUP BY medallion"
	mysql_date_format = "2006-01-02"
	end_day_time      = " 23:59:59"
	join_ids_with     = "', '"
)

func BuildQuery(key structs.Request) (res wraps.Query, err error) {
	date := key.Date.Format(mysql_date_format)
	res.Request = fmt.Sprintf(query_format, strings.Join(key.IDs, join_ids_with))
	res.Args = []interface{}{date, date + end_day_time}
	return
}

package main

import (
	"database/sql"
	errors2 "errors"
	"fmt"
	"github.com/pkg/errors"
)

//我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。
//为什么，应该怎么做请写出代码？

func dbQuery(sqlStr string) (*sql.Row, error) {
	// simulation db do some query and return sql.ErrNoRows
	return nil, sql.ErrNoRows
}

func queryFunc(params interface{}) (*sql.Row, error) {
	//var conn = sql.Conn{}
	rows, err := dbQuery("")
	if err != nil {
		// wrap 这个error
		err = errors.Wrapf(err, "query db error, params %v", params)
		return nil, err
	}
	return rows, nil
}

var responseErr = errors.New("no-data")

func getData(params interface{}) (interface{}, error) {
	dbRes, err := queryFunc(params)
	if err != nil {
		// 判断error是否需要在这层处理
		if errors2.Is(err, sql.ErrNoRows) {
			// 转换成业务统一error
			return nil, responseErr
			// 或者直接返回得到空数据
			// log.info("get query no data", err)
			// return nil, nil
		}
		// 往上抛
		return nil, err
	}
	return dbRes, nil
}

func main() {
	var params = map[string]interface{}{"id": 1}
	data, err := getData(params)
	fmt.Println("get error: ", err)
	fmt.Println("get data: ", data)
}
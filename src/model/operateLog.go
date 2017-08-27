/**
 * 操作日志model实现
 *
 * @author zhangqichao
 * Created on 2017-08-18
 */
package model

import (
	"fmt"
	"time"

	"../pdb"
)

type OperateLog struct {
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	OperateType string `json:"operateType"`
	Operator    string `json:"operator"`
	IP          string `json:"ip"`
	Content     string `json:"content"`

	CreatedAt time.Time `json:"createdAt"`
}

func OperateLogTableName() string {
	return "t_operate_log"
}

func (m *OperateLog) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(type,operate_type,operator,ip,content,created_at) "+
		"VALUES($1,$2,$3,$4,$5,$6)", OperateLogTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	m.CreatedAt = time.Now()

	_, err = stmt.Exec(m.Type, m.OperateType, m.Operator, m.IP, m.Content, m.CreatedAt)
	return
}

func FindOperateLogs(condition, limit, order string) ([]OperateLog, error) {
	result := []OperateLog{}
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,type,operate_type,operator,ip,content,created_at FROM %s %s %s %s", OperateLogTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}

	for rows.Next() {
		tmp := OperateLog{}
		err = rows.Scan(&tmp.ID, &tmp.Type, &tmp.OperateType, &tmp.Operator, &tmp.IP, &tmp.Content, &tmp.CreatedAt)
		result = append(result, tmp)
	}
	return result, err
}

func UpdateOperateLogs(update, condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("UPDATE %s SET %s %s", OperateLogTableName(), update, condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}

func CountOperateLogs(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", OperateLogTableName(), condition)).Scan(&count)
	return
}

func DeleteOperateLogs(condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("DELETE FROM %s %s", OperateLogTableName(), condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}

package model

import (
	"time"
	"fmt"
	"../pdb"
	"github.com/bitly/go-simplejson"
)

/*
DROP TABLE IF EXISTS t_app;
CREATE TABLE t_app (
"id" serial NOT NULL,
"icon" text,
"app_id" varchar(16) COLLATE "default",
"name" varchar(128) COLLATE "default",
"version" varchar(16) COLLATE "default",
"describe" varchar(255) COLLATE "default",
"developer" int4,
"valid" bool DEFAULT TRUE,
"file" varchar(256) COLLATE "default",
"src" varchar(256) COLLATE "default",
"expend" jsonb NOT NULL,
"download_count" int4 DEFAULT -1,
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */


type T_app struct {
	ID            int64     `json:"id"`
	Icon          string    `json:"icon"`
	AppId         string    `json:"app_id"`
	Name          string    `json:"name"`
	Version       string    `json:"version"`
	Describe      string    `json:"describe"`
	Developer     int       `json:"developer"`
	Valid         bool      `json:"valid"`
	File          string    `json:"file"`
	Src           string    `json:"src"`
	Expend        *simplejson.Json    `json:"expend"`
	DownloadCount int       `json:"download_count"`
	CreatedAt     time.Time `json:"created_at"`
}
func AppTableName() string {
	return "t_app"
}

func (c *T_app) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(icon,app_id,version,name,describe,developer,valid,file,src,expend,download_count,created_at) "+
			  "VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)", AppTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.CreatedAt = time.Now()
	expend :="{}"
	if c.Expend!=nil{
		bs,err:=c.Expend.MarshalJSON()
		if err!=nil{
			return err
		}
		expend =string(bs)
	}
	_, err = stmt.Exec([]byte(c.Icon),c.AppId,c.Version,c.Name,c.Describe, c.Developer, c.Valid,c.File,c.Src,expend,c.DownloadCount,c.CreatedAt)
	return
}

func FindApps(condition, limit, order string) (result []T_app,err error) {
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,icon,app_id,version,name,describe,developer,valid,file,src,expend,download_count,created_at" +
			  " FROM %s %s %s %s", AppTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}
	for rows.Next() {
		tmp := T_app{}
		bs:=new([]byte)
		err = rows.Scan(&tmp.ID, &tmp.Icon,&tmp.AppId, &tmp.Version, &tmp.Name, &tmp.Describe, &tmp.Developer, &tmp.Valid, &tmp.File, &tmp.Src, bs,&tmp.DownloadCount,&tmp.CreatedAt)
		tmp.Expend,_=simplejson.NewJson(*bs)
		if err==nil {
			result = append(result, tmp)
		}
	}
	return result, err
}

func CountApps(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", AppTableName(), condition)).Scan(&count)
	return
}

func UpdateApps(update, condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("UPDATE %s %s %s", AppTableName(), update, condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = stmt.Exec()
	return
}

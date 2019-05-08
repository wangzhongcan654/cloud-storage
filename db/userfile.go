package db

import (
	"cloud-storage/db/mysql"
	"time"
	"fmt"
)

type UserFile struct{
	UserName string
	FileHash string
	FileName string
	FileSize int64
	UploadAt string
	LastUpdated string
}

func OnUserFileUploadFinished(username,filehash,filename string,filesize int64)bool{
	stmt,err:=mysql.DBConn().Prepare(
		"insert ignore into tbl_user_file(`user_name`,`file_shah1`,`file_name`,"+
			"`file_size`,`upload_at`) values(?,?,?,?,?)")
	if err!=nil{
		return false
	}
	defer stmt.Close()

	_,err=stmt.Exec(username,filehash,filename,filesize,time.Now())
	if err!=nil{
		return false
	}
	return true
}
//QueryUserFileMeta: 批量获取文件元信息
func QueryUserFileMeta(username string,limit int)([]UserFile,error){
	stmt,err:=mysql.DBConn().Prepare(
		"select file_sha1,file_name,file_size,upload_at,last_update from"+
			"tbl_user_file where user_name=?limit?")
	if err!=nil{
		return nil,err
	}
	defer stmt.Close()

	rows,err:=stmt.Query(username,limit)
	if err!=nil{
		return nil,err
	}
	var userFiles []UserFile
	for rows.Next(){
		ufile:=UserFile{}
		err=rows.Scan(&ufile.FileHash,&ufile.FileName,&ufile.FileSize,&ufile.UploadAt,&ufile.LastUpdated)
		if err!=nil{
			fmt.Println(err.Error())
			break
		}
		userFiles=append(userFiles,ufile)
	}
	return userFiles,nil
}
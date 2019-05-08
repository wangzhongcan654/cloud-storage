package db

import (
	"cloud-storage/db/mysql"
	"fmt"
)

func UserSignup(username string,passwd string)bool{
	//防止SQL注入
	stmt,err:=mysql.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`,`user_pwd`)values(?,?)")
	if err!=nil{
		fmt.Println("Failed to")
		return false
	}
	defer stmt.Close()
	ret,err:=stmt.Exec(username,passwd)
	if err!=nil{
		fmt.Println("Failed to insert,err:"+err.Error())
		return false
	}
	if rowsAffected,err:=ret.RowsAffected();nil==err&&rowsAffected>0{
		return true
	}
	return false
}


func UserSignin(username string,password string)bool{
	stmt,err:=mysql.DBConn().Prepare("select *from tbl_user where user_name=? limit 1")
	if err!=nil{
		fmt.Println(err.Error())
		return false
	}
	row,err:=stmt.Query(username)
	if err!=nil{
		fmt.Println(err.Error())
		return false
	}else if row==nil{
		fmt.Println("username not found")
		return false
	}
	pRow:=mysql.ParseRows(row)
	if len(pRow)>0&&string(pRow[0]["user_pwd"].([]byte))==password{
		return true
	}
	return false
}

//UpdateToken:
func UpdateToken(username string,token string)bool{
	stmt,err:=mysql.DBConn().Prepare(
		"replace into tbl_user_token(`user_name`,`user_token`) values(?,?)")
	if err!=nil{
		fmt.Println(err.Error())
		return false
	}
	_,err=stmt.Exec(username,token)
	if err!=nil{
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct{
	Username string
	Email string
	Phone string
	SignupAt string
	LastActiveAt string
	Status int
}
func GetUserInfo(username string)(User,error){
	user:=User{}
	stmt,err:=mysql.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=?limit 1")
	if err!=nil{
		fmt.Println(err.Error())
		return user,err
	}
	defer stmt.Close()

	//执行查询的操作
	err=stmt.QueryRow(username).Scan(&user.Username,&user.SignupAt)
	if err!=nil{
		return user,err
	}
	return user,nil
}



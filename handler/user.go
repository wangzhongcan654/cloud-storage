package handler

import (
	"net/http"
	"io/ioutil"
	"cloud-storage/util"
	"cloud-storage/db"
	"fmt"
	"time"
)

const(
	pwd_salt="*#890"
)

//SignupHander：处理用户注册请求
func SignupHander(w http.ResponseWriter,r *http.Request){
	if r.Method==http.MethodGet{
		data,err:=ioutil.ReadFile("./static/view/signup.html")
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username:=r.Form.Get("username")
	passwd:=r.Form.Get("password")
	if len(username)<3||len(passwd)<5{
		w.Write([]byte("Invalid parameter"))
		return
	}
	enc_passwd:=util.Sha1([]byte(passwd+pwd_salt))
	suc:=db.UserSignup(username,enc_passwd)
	if suc{
		w.Write([]byte("sunccess"))
	}else{
		w.Write([]byte("failed"))
	}
}

//SignInHandler：登陆接口
func SignInHandler(w http.ResponseWriter,r *http.Request){
	//1.校验用户名或密码
	r.ParseForm()
	username:=r.Form.Get("username")
	password:=r.Form.Get("password")
	encPasswd:=util.Sha1([]byte(password+pwd_salt))
	pwdChecked:=db.UserSignin(username,encPasswd)

	if pwdChecked{
		w.Write([]byte("failed"))
		return
	}
	//2.生成访问凭证（token）
	token:=GetToken(username)
	upRes:=db.UpdateToken(username,token)
	if !upRes{
		w.Write([]byte("FAILED"))
		return
	}
	//3.登陆成功后重定向到首页
	//w.Write([]byte("http://"+r.Host+"/static/view/home.html"))
	resp:=util.RespMsg{
		Code:0,
		Msg:"OK",
		Data:struct{
			Location string
			Username string
			Token string
		}{
			Location:"http://"+r.Host+"/static/view/home.html",
			Username:username,
			Token:token,
		},
	}
	w.Write(resp.JSONBytes())
}

//UserInfoHandler: 查询用户信息
func UserInfoHandler(w http.ResponseWriter,r *http.Request){
	//1.解析请求参数
	r.ParseForm()
	username:=r.Form.Get("username")
	//token:=r.Form.Get("token")
	//2.验证Token是否有效
	//isValid:=IsTokenValid(token)
	//if !isValid{
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	//3.查询用户信息
	user,err:=db.GetUserInfo(username)
	if err!=nil{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//4.组装并且响应用户数据
	resp:=util.RespMsg{
		Code:0,
		Msg:"OK",
		Data:user,
	}
	w.Write(resp.JSONBytes())
}

//IsTokenValid: 判断token的有效性
func IsTokenValid(token string)bool{
	//TODO:判断token的时效性，是否过期
	//TODO:从数据库表查询对应的username对应的token
	//对比两个
	if len(token)!=40{

	}
	return false
}



func GetToken(username string)string{
	//md5(username+timestamp+token_salt)+timestamp[:8]
	ts:=fmt.Sprintf("%x",time.Now().Unix())
	tokenPrefix:=util.MD5([]byte(username+ts+"_tokensalt"))
	return tokenPrefix+ts[:8]
}
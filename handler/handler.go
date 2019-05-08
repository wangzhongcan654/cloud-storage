package handler

import (
	"net/http"
	"io/ioutil"
	"io"
	"fmt"
	"os"
	"cloud-storage/meta"
	"time"
	"cloud-storage/util"
	"encoding/json"
	"cloud-storage/db"
	"strconv"
)

//上传文件
func UploadHandler(w http.ResponseWriter,r *http.Request){
	fmt.Printf("received request\n")
	if r.Method=="GET"{
		//返回上传html页面
		data,err:=ioutil.ReadFile("./static/view/index.html")
		if err!=nil{
			io.WriteString(w,"internel server error")
			return
		}
		io.WriteString(w,string(data))
	}else if r.Method=="POST"{
		//接收文件流存储到本地目录
		file,head,err:=r.FormFile("file")
		if err!=nil{
			fmt.Printf("Failed to get data,err%s\n",err.Error())
			return
		}
		defer file.Close()

		fileMeta:=meta.FileMeta{
			FileName:head.Filename,
			Location:"./storefile/"+head.Filename,
			UploadAt:time.Now().Format("2006-01-02 15:04:32"),
		}
		newFile,err:=os.Create(fileMeta.Location)
		if err!=nil{
			fmt.Println(err)
			return
		}
		defer newFile.Close()

		fileMeta.FileSize,err=io.Copy(newFile,file)
		if err!=nil{
			fmt.Printf("Failed to save data into file,err:%s\n",err.Error())
			return
		}

		newFile.Seek(0,0)
		fileMeta.FileSha1=util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)
		_=meta.UpdateFileMetaDB(fileMeta)
		//TODO:更新用户文件表记录
		username:=r.Form.Get("username")
		suc:=db.OnUserFileUploadFinished(username,fileMeta.FileSha1,fileMeta.FileName,fileMeta.FileSize)
		if suc{
			http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
		}else{
			w.Write([]byte("Upload failed"))
		}
	}
}


func UploadSucHandler(w http.ResponseWriter,r *http.Request){
	io.WriteString(w,"Upload finished!")
}

func FileQueryHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	limitCnt,_:=strconv.Atoi(r.Form.Get("limit"))
	username:=r.Form.Get("username")

	userFiles,err:=db.QueryUserFileMeta(username,limitCnt)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data,err:=json.Marshal(userFiles)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

//GetFileMetaHandler：获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	filehash:=r.Form["filehash"][0]
	//fMeta:=meta.GetFileMeta(filehash)
	fMeta,err:=meta.GetFileMetaDB(filehash)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data,err:=json.Marshal(fMeta)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

//DownloadHander: 文件下载接口
func DownloadHander(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	fsha1:=r.Form.Get("filehash")
	fileMeta:=meta.GetFileMeta(fsha1)

	f,err:=os.Open(fileMeta.Location)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data,err:=ioutil.ReadAll(f)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/octect-stream")
	w.Header().Set("Content-Descrption","attachment;filename=\""+fileMeta.FileName+"\"")
	w.Write(data)
}

//FileMetaUpdateHandler: 更新元信息接口（重命名）
func FileMetaUpdateHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()

	opType:=r.Form.Get("op")
	fileSha1:=r.Form.Get("filehash")
	newFileName:=r.Form.Get("filename")

	if opType!="0"{
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method!="POST"{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	curFileMeta:=meta.GetFileMeta(fileSha1)

	curFileMeta.FileName=newFileName
	meta.UpdateFileMeta(curFileMeta)

	w.WriteHeader(http.StatusOK)

	data,err:=json.Marshal(curFileMeta)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

//FileDeleteHandler:删除文件及文件元信息
func FileDeleteHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()

	fileSha1:=r.Form.Get("filehash")

	fMeta:=meta.GetFileMeta(fileSha1)

	os.Remove(fMeta.Location)

	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}

//
func TryFastUploadHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	//1.解析请求参数
	username:=r.Form.Get("username")
	filehash:=r.Form.Get("filehash")
	filename:=r.Form.Get("filename")
	filesize,_:=strconv.Atoi(r.Form.Get("filesize"))

	//2.从文件表中查询相同的hash的文件记录
	_,err:=meta.GetFileMetaDB(filehash)
	if err!=nil{
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//4.上传过则将文件信息写入用户文件表，返回成功
	suc:=db.OnUserFileUploadFinished(username,filehash,filename,int64(filesize))
	if suc{
		resp:=util.RespMsg{
			Code:0,
			Msg:"秒传成功",
		}
		w.Write(resp.JSONBytes())
		return
	}else{
		resp:=util.RespMsg{
			Code:-2,
			Msg:"秒传失败，请稍后重试",
		}
		w.Write(resp.JSONBytes())
		return
	}
}
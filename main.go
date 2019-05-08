package main

import (
	"net/http"
	"cloud-storage/handler"
	"fmt"
)

func main(){

	http.HandleFunc("/file/upload",handler.UploadHandler)
	http.HandleFunc("/file/upload/suc",handler.UploadSucHandler)
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	http.HandleFunc("/file/fastupload",handler.HTTPInterceptor(handler.TryFastUploadHandler))
	http.HandleFunc("/file/download",handler.DownloadHander)
	http.HandleFunc("/file/update",handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete",handler.FileDeleteHandler)
	http.HandleFunc("/user/signup",handler.SignupHander)
	http.HandleFunc("/user/signin",handler.SignInHandler)
	http.HandleFunc("/user/info",handler.HTTPInterceptor(handler.UserInfoHandler))

	err:=http.ListenAndServe("localhost:8080",nil)
	if err!=nil{
		fmt.Printf("Failed to start server, err:%s",err.Error())
	}
}

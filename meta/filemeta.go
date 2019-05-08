package meta

import "cloud-storage/db"

//FileMeta: 文件元信息结构
type FileMeta struct{
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init(){
	fileMetas=make(map[string]FileMeta)
}

//UpdateFileMeta:新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta){
	fileMetas[fmeta.FileSha1]=fmeta
}

//UpdateFileMetaDB:新增文件元信息到DB
func UpdateFileMetaDB(fmeta FileMeta)bool {
	return db.OnFileUploadFinished(fmeta.FileSha1,fmeta.FileName,fmeta.FileSize,fmeta.Location)
}


func GetFileMetaDB(fileSha1 string)(FileMeta,error){
	tfile,err:=db.GetFileMeta(fileSha1)
	if err!=nil{
		return FileMeta{},err
	}
	fmeta:=FileMeta{
		FileSha1:tfile.FileHash,
		FileName:tfile.FileName.String,
		FileSize:tfile.FileSize.Int64,
		Location:tfile.FileAddr.String,
	}
	return fmeta,nil
}


func GetFileMeta(fileSha1 string)FileMeta{
	return fileMetas[fileSha1]
}

func RemoveFileMeta(fileSha1 string){
	delete(fileMetas,fileSha1)
}
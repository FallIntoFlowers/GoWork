package utils

import (
	"os"
	"fmt"
	"io"
	"crypto/rand"
	"encoding/base64"
		"crypto/md5"
	"github.com/astaxie/beego"
	"errors"
	"encoding/json"
)

type  TempInfo struct{
	ChunkSize int64 `json:"ChunkSize"`
	FileHash string `json:"FileHash"`
	FileId string `json:"FileId"`
	FileName string `json:"FileName"`
	FileSize int64 `json:"FileSize"`
	Label string `json:"Label"`
}

func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

//生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}

//根据传入路径判断文件或目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	//file, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
	}
	//beego.Debug(file.IsDir())
	return true
}

//根据传入路径判断文件是否是目录
func IsDir(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) && file.IsDir() {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
	}
	return file.IsDir()
}

//根据传入路径判断是否是文件
func IsFile(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) && !file.IsDir() {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
	}
	return !file.IsDir()
}

//根据传入路径判断文件是否可写
func IsWritable(path string) bool {
	return true;
}

//获取路径对应的文件
func GetFile(path string) os.FileInfo {
	file, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) && !file.IsDir() {
			return file
		}
		if os.IsNotExist(err) {
			return nil
		}
	}
	if !file.IsDir() {
		return file
	} else {
		return nil
	}
}

func GetFileSize(infoPath string) int64 {
	if fileInfo, err := os.Stat(infoPath); err == nil {
		return fileInfo.Size()
	} else {
		beego.Error("GetFileSize:" + err.Error())
		return 0
	}
}

func DeleteFile(fileName string) error {
	//删除文件
	if !IsFile(fileName) {
		return errors.New("file isn't exist")
	}
	if err := os.Remove(fileName); err != nil {
		beego.Error("DeleteFile: " + err.Error())
		return err
	}
	return nil
}

//获取文件中的json数据
func GetJsonFileInfo(infoPath string) (temp TempInfo, n int64)  {
	//data := ; //获取文件中数据
	tempInfoFile, err := os.Open(infoPath)
	defer tempInfoFile.Close()
	if err != nil {
		//this.ReturnFailedJson(err, "Failed to find file!")
	}
	tempInfoSize := GetFileSize(infoPath)
	if tempInfoSize == 0 {
		beego.Error("缓存文件大小为0")
		return TempInfo{}, 0
	}
	data := TempInfo{}
	var info = make([]byte, tempInfoSize)
	if _, err := tempInfoFile.Read(info); err == nil {
		if err := json.Unmarshal(info, &data); err != nil {
			beego.Error(err.Error())
			return TempInfo{}, 0
		}
	} else {
		beego.Error(err.Error())
		return TempInfo{}, 0
	}
	return data, tempInfoSize
}
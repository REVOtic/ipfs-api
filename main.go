package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"os"
	"io/ioutil"
	
	"github.com/gin-gonic/gin"
	"ipfs_api/pkg/our_infura"
	//"github.com/wabarc/ipfs-pinner/pkg/infura"
)

func main() {
	router := gin.Default()

	root := router.Group("/")
	{
		root.POST("upload", upload)
		root.POST("retrieve", retrieve)
	}
	
	router.Run(":9090")
}

func upload(c *gin.Context) {
	file, err := c.FormFile("uploadFile")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
	}
	
	if file != nil{
		filename := filepath.Base(file.Filename)
	
		if err := c.SaveUploadedFile(file, "uploadedFiles/"+filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s \n", err.Error()))
		}
		
		c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully \n", file.Filename))
		
		//------------------------------------------------------------------------------------------
		filePath := filepath.Join("/home/oem/go/src/ipfs_api/uploadedFiles/", filename)
		
		cid, err := ipfs_protocol.PinFile(filePath)
		//cid, err := infura.PinFile(filePath)

		if err != nil {
			c.String(http.StatusOK, fmt.Sprintf("ipfs-pinner: %s \n", err.Error()))
		} else {
			c.String(http.StatusOK, fmt.Sprintf("Pinned file hash: \n %s \n", cid))
		}
	}
}

func retrieve(c *gin.Context) {
	hash := c.PostForm("retrieveFile")		
	if hash != ""{
		data, err := ipfs_protocol.RetrieveFile(hash)
		if err != nil {
			c.String(http.StatusOK, fmt.Sprintf("File retrieval error: %s \n", err.Error()))
		} else {
			filePath := filepath.Join("/home/oem/go/src/ipfs_api/retrievedFile/", hash)
			var file, errf = os.Create(filePath)
			if errf != nil { return }
			defer file.Close()
			
			err = ioutil.WriteFile(filePath, []byte(data),0644)
			if err!=nil{
				c.String(http.StatusOK, fmt.Sprintf("Writing file error: %s \n", err.Error()))
			}
			
			c.String(http.StatusOK, fmt.Sprintf("File retried at:\n%s \n", filePath))
		}
	}
}


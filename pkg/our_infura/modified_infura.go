package ipfs_protocol

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	fp "path/filepath"
	"time"
)

const (
	PROTOCAL = "http"
	HOST = "127.0.0.1"
	PORT = 5001
)

func PinFile(filepath string) (string, error) {
	uri := fmt.Sprintf("%s://%s:%d/api/v0/add", PROTOCAL, HOST, PORT)

	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer m.Close()

		part, err := m.CreateFormFile("file", fp.Base(file.Name()))
		if err != nil {
			return
		}

		if _, err = io.Copy(part, file); err != nil {
			return
		}
	}()
		
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	req, err := http.NewRequest(http.MethodPost, uri, r)
	req.Header.Add("Content-Type", m.FormDataContentType())
	
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	fmt.Printf("%s", data)
	
	var dat map[string]interface{}
	if err := json.Unmarshal(data, &dat); err != nil {
		return "", err
	}

	if hash, ok := dat["Hash"].(string); ok {
		return hash, nil
	}
	
	return "", fmt.Errorf("Pin file to local failed.")	
}

func RetrieveFile(hash string) (string, error) {
	//url :=  fmt.Sprintf("%s://%s:%d/api/v0/get?arg=%s", PROTOCAL, HOST, PORT, hash)
	url :=  fmt.Sprintf("%s://%s:%d/api/v0/object/data?arg=/ipfs/%s", PROTOCAL, HOST, PORT, hash)

	client := &http.Client {
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	
	return string(body), nil
}

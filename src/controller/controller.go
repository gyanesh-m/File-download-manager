package controller

import (
	"../model"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func HomePage(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(200)
	fmt.Println("Health:OK")
}
func getUUID()(string){
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	//fmt.Println(uuid)
	return uuid
}

var storage []model.Response

func serialDownload(w http.ResponseWriter, urls []string)(model.Response){
	fmt.Println("session")
	sessionId := getUUID()
	var data model.Response
	data.ID = sessionId
	data.Status = "QUEUED"
	data.DownloadType = "SERIAL"
	wd := os.TempDir()  + sessionId +"/"
	os.MkdirAll(wd,0777)
	mapFiles := make(map[string]string)
	fmt.Println("starting")
	for _,url := range urls{
		data.StartTime = time.Now()
		//fmt.Println(url)
		wd += getUUID()

		//fmt.Println(mapFiles)
		if err := DownloadFile(wd, url); err != nil {
			fmt.Println(err)
			fmt.Println("ERROR ")
			data.Status = "FAILED"
		} else{
			mapFiles[url] = wd
		}
		//fmt.Println("map:",mapFiles)
		//fmt.Println(data)
	}
	if(data.Status=="QUEUED"){
		data.Status = "SUCCESSFUL"
	}
	//fmt.Println("error")
	data.Files = mapFiles
	data.EndTime = time.Now()
	fmt.Println("data::")
	fmt.Println(data)
	//jsonData ,_ := json.Marshal(data)
	json.NewEncoder(w).Encode(data)
	//fmt.Fprint(w,jsonData)
	return data

}


func concurrentDownload(urls []string){

}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func Download(w http.ResponseWriter, r* http.Request){
	reqBody, _ := ioutil.ReadAll(r.Body)
	var model model.Data
	json.Unmarshal(reqBody,&model)
	//fmt.Println(model)
	//fmt.Fprint(w,model)
	//var result src.Response
	switch model.Type{
	case "serial":
		serialDownload(w,model.Urls)
		//case "concurrent":
		//	result = concurrentDownload(model.Urls)
	}


}
func Status(w http.ResponseWriter, r* http.Request){

}
func Browse(w http.ResponseWriter, r* http.Request){

}


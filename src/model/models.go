package model

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)
type Data struct {
	Type string   `json:"type"`
	Urls []string `json:"urls"`
}


type Packet struct{
	Url string
	Path string
}
type Error struct{
	internalCode int `json:"internal_code"`
	message string `json:"message"`
}

type Id struct{
	Id string `json:id`
}
type Response struct {
	Id           string `json:id`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       string `json:"status"`
	DownloadType string `json:"download_type"`
	Files        map[string]string `json:files`
}

type Requests interface{
	Download(w http.ResponseWriter)
}
type Serial struct{
	Data Data
	Response Response

}
type Concurrent struct{
	Threads int
	Data Data
	Response *Response
}
func GetUUID()(string){
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

func(s *Serial) Download(w http.ResponseWriter){
	sessionId := GetUUID()
	s.Response.Id= sessionId
	s.Response.Status = "QUEUED"
	s.Response.DownloadType = "SERIAL"
	wd := os.TempDir()  + sessionId +"/"
	os.MkdirAll(wd,0777)
	mapFiles := make(map[string]string)
	fmt.Println("starting")
	for _,url := range s.Data.Urls{
		s.Response.StartTime = time.Now()
		wd += GetUUID()
		if err := DownloadFile(wd, url); err != nil {
			fmt.Println(err)
			fmt.Println("ERROR ")
			s.Response.Status = "FAILED"
		} else{
			mapFiles[url] = wd
		}
	}
	if(s.Response.Status=="QUEUED"){
		s.Response.Status = "SUCCESSFUL"
	}
	s.Response.Files = mapFiles
	s.Response.EndTime = time.Now()
	fmt.Println("data::")
	fmt.Println(s.Response)
	//json.NewEncoder(w).Encode(s.Response.Id)
}

func(c *Concurrent)Download(){
	//status := "QUEUED"

	//c.Response = resp

	c.Response.Status = "QUEUED"
	c.Response.DownloadType = "CONCURRENT"
	c.Response.StartTime = time.Now()
	dataChan := make(chan Packet)
	resChan := make(chan int)
	sessionId := GetUUID()
	c.Response.Id = sessionId
	wd := os.TempDir()  + sessionId +"/"
	os.MkdirAll(wd,0777)
	//end_time := c.Response.StartTime
	fmt.Println("starting concurrent")
	mapFiles := make(map[string]string)

	c.Response.Files = mapFiles
	c.Response.EndTime = c.Response.StartTime
	//spawning threads
	for i:=0;i<c.Threads;i++{
		fmt.Println("Spawning Thread no: "+string(i))
		go Fetch(dataChan,resChan,mapFiles,c)
	}
	fmt.Println("Sending packets")
	go populateChannel(c,wd,dataChan)
	fmt.Println("Evaluation started")
	go EvaluateEnd(dataChan,len(c.Data.Urls),resChan,c,mapFiles)
	fmt.Println("Evaluation done")

}
func populateChannel(c *Concurrent,wd string, dataChan chan Packet){

	for id,url := range c.Data.Urls{
		fmt.Println(id,url)
		temp := wd
		wd += GetUUID()
		var pkt Packet
		pkt.Url = url
		pkt.Path = wd
		dataChan<-pkt
		wd = temp
	}
	fmt.Println("Packets sent")
}
func EvaluateEnd(dataChan chan Packet,total int,resChan chan int,c *Concurrent, mapFiles map[string]string){
	counter := 0
	for range resChan{
		counter++
		//fmt.Println("counter:",counter)
		if counter==total{
			close(dataChan)
			//close(resChan)
			fmt.Printf("COMPLETED")
			c.Response.Files = mapFiles
			c.Response.EndTime = time.Now()
			//resp.EndTime = time.Now()
			if(c.Response.Status == "QUEUED"){
				c.Response.Status = "SUCCESSFUL"
			}

			return
			}
	}
}
func Fetch(dataChan chan Packet,resChan chan int,mapFiles map[string]string, c * Concurrent){
	for data := range dataChan{
			if(data.Url!="") {
				if err := DownloadFile(data.Path, data.Url); err != nil {
					fmt.Println(err)
					fmt.Println("ERROR ")
					c.Response.Status = "FAILED"
				} else {
					mapFiles[data.Url] = data.Path
				}
				resChan<-1
			}

		}
	fmt.Println("THREAD CLOSED.")
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

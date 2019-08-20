package model

import "time"

type Data struct {
	Type string   `json:"type"`
	Urls []string `json:"urls"`
}

type Error struct{
	internalCode int `json:"internal_code"`
	message string `json:"message"`
}

type Response struct {
	ID           string `json:"id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       string `json:"status"`
	DownloadType string `json:"download_type"`
	Files        map[string]string `json:files`
}
type serial struct{
	data Data
	response Response
}
type concurrent struct{
	threads int
	data Data
	response Response
}
type Download interface{
	download()
}
func(s serial) download(){
	
}
func(c concurrent) download(){

}

package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gyanesh-m/File-download-manager/src/model"
	"io/ioutil"
	"net/http"
)

var requests = make(map[string] *model.Response)

func HealthCheck(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(200)
	fmt.Println("Health:OK")
}

func Download(w http.ResponseWriter, r* http.Request){
	reqBody, _ := ioutil.ReadAll(r.Body)
	var data model.Data
	json.Unmarshal(reqBody,&data)
	fmt.Print(data)
	switch data.Type {
	case "serial":
		{
			s := new(model.Serial)
			s.Data = data
			s.Download(w)
			requests[s.Response.Id] = &s.Response
			idObj := model.Id{Id:s.Response.Id}
			json.NewEncoder(w).Encode(idObj)
			//requests[s.Response.ID.ID] = s
			//s.Response.Status = "something"
		}
	case "concurrent":
		{
			c := new(model.Concurrent)
			c.Data = data
			c.Threads = 5
			c.Response = &model.Response{}
			c.Download()
			idObj := model.Id{Id:c.Response.Id}
			json.NewEncoder(w).Encode(idObj)
			requests[c.Response.Id] = c.Response
		}

	}


}


func Status(w http.ResponseWriter, r* http.Request){
	id := mux.Vars(r)["id"]
	if val, ok := requests[id]; ok{

		json.NewEncoder(w).Encode(val)
	}else{
		w.WriteHeader(200)
		type  returnObj struct  {
			internal_code int
			message string
		}
		var ret returnObj
		ret.internal_code = 4002
		ret.message = "unknown download ID"
		json.NewEncoder(w).Encode(ret)

	}

}
func Browse(w http.ResponseWriter, r* http.Request){

}


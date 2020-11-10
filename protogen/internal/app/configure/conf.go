package configure

import (
	"encoding/json"
	"log"
	"os"
)

// Desp ...
type Desp struct {
	In string `json:"in"`
	Getway bool `json:"getway"`
	Out []string `json:"out"`
}

// Conf ...
type Conf struct {
	Dir  string `json:"dir"`
	Desp []Desp `json:"desp"`
}

// NewConf ...
func NewConf() *Conf {
	file, err := os.Open("./protogen/protoconf.json")
    if err != nil{
        log.Fatal(err) 
        os.Exit(1) 
    }
	defer file.Close() 
	
	// decode json to struct
	conf := Conf{}
	if err = json.NewDecoder(file).Decode(&conf); err != nil {
		log.Fatal(err) 
	}
	
	return &conf
}
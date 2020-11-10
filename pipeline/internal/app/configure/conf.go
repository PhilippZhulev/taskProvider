package configure

import (
	"encoding/json"
	"log"
	"os"
)

// Conf ...
type Conf struct {
	Name string `json:"name"`
	Version string `json:"version"`
	Services [][]string `json:"services"`
	Log bool `json:"log"`
	LogURL string `json:"logUrl"`
	ErrorLog bool `json:"errorLog"`
	ErrorLogURL string `json:"errorLogUrl"`
	Bot bool `json:"bot"`
	BotAuth []string `json:"botAuth"`
	BotToken string `json:"botToken"`
	BotDebug bool `json:"botDebug"`
}

// NewConf ...
func NewConf() *Conf {
	file, err := os.Open("./pipeline/pipeconf.json")
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
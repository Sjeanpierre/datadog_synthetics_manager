package main

import (
	"fmt"
	"github.com/sjeanpierre/datadog_synthetics_manager/lib"
	"io/ioutil"
	"log"
)

func main() {
	lib.ListSyntheticsChecks()
	fmt.Println("done")
	d := readFile("./data/api-app.example.com.yml")
	dd,err := lib.YAMLtoStruct(d)
	if err != nil {
		log.Printf("it's broken %s",err)
	}
	//fmt.Printf("%+v",dd)
	fmt.Printf("%s | %s\n",*dd.PublicId,*dd.Name)
	ddd, err := lib.UpdateSyntheticsTest(dd)
	if err != nil {
		log.Printf("error in Test update %s",err)
	}
	fmt.Printf("%+v",ddd)
}


func readFile(f string) (m []byte) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal("could not read build_trigger.yml in current directory")
	}
	return data
}
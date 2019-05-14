package lib

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"gopkg.in/zorkian/go-datadog-api.v2"
	"io/ioutil"
	"log"
	"os"
)


func ListSyntheticsChecks() {
	c := datadog.NewClient(os.Getenv("SYNTH_MAN_DD_API_KEY"),os.Getenv("SYNTH_MAN_DD_APP_KEY"))
	checks,err := c.GetSyntheticsTests()
	if err != nil {
		log.Fatalf("Could not list synth checks %s",err)
	}
	if len(checks) < 1 {
		fmt.Println("No checks exist")
	}
	for _, check := range checks {
		fmt.Printf("%s | %s | %s\n",*check.PublicId,*check.Name,*check.Status)
	}
}

func GetSyntheticsCheck(publicID string) {
	c := datadog.NewClient(os.Getenv("SYNTH_MAN_DD_API_KEY"),os.Getenv("SYNTH_MAN_DD_APP_KEY"))
	check,err := c.GetSyntheticsTest(publicID)
	if err != nil {
		log.Fatalf("Could not get synth checks %s",err)
	}
	if check == nil {
		fmt.Println("could not locate check with public ID of ",publicID)
	}
		fmt.Printf("%s | %s | %s\n",*check.PublicId,*check.Name,*check.Status)
}

// Encountered issue with provided struct from DD lib where only JSON tags were present as struct tags
// We want to provide users the ability to use YAML for their config due to the readability
// Without YAML struct tags, there would be no easy way to get our YAML data into the needed struct
// According to http://ghodss.com/2014/the-right-way-to-handle-yaml-in-golang/ we can convert the YAML bytes to JSON
// then unmarshal the json as needed into the struct. This seems to work.
func YAMLtoSynth(data []byte) (datadog.SyntheticsTest,error) {
	//fmt.Printf(string(data))
	var dataStruct datadog.SyntheticsTest
	//err := yaml.Unmarshal(data, &datadog.SyntheticsTest{})
	jsonBytes, err := yaml.YAMLToJSON(data)
	if err != nil {
		return dataStruct,fmt.Errorf("could not unmarshall YAML to Struct %s",err)
	}
	err = json.Unmarshal(jsonBytes,&dataStruct)
	if err != nil {
		return dataStruct,fmt.Errorf("could not unmarshall to JSON to Struct %s",err)
	}
	return dataStruct,nil
}

func UpdateSyntheticsTest(test datadog.SyntheticsTest) (datadog.SyntheticsTest, error) {
	c := datadog.NewClient(os.Getenv("SYNTH_MAN_DD_API_KEY"),os.Getenv("SYNTH_MAN_DD_APP_KEY"))
	id := *test.PublicId
	test.PublicId = nil
	test.ModifiedBy = nil
	t, err := c.UpdateSyntheticsTest(id, &test)
	if err != nil {
		return test, fmt.Errorf("encountered error updating synthetics %s",err)
	}
	return *t,nil
}

func CreateSyntheticsTest(test datadog.SyntheticsTest) (datadog.SyntheticsTest, error) {
	c := datadog.NewClient(os.Getenv("SYNTH_MAN_DD_API_KEY"),os.Getenv("SYNTH_MAN_DD_APP_KEY"))
	t, err := c.CreateSyntheticsTest(&test)
	if err != nil {
		return test, fmt.Errorf("encountered error creating synthetics %s",err)
	}
	return *t,nil
}

func ReadFile(f string) (m []byte) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf("error reading file %s, %s",f,err)
	}
	return data
}
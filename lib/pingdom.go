package lib

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin/json"
	"gopkg.in/zorkian/go-datadog-api.v2"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/russellcardullo/go-pingdom/pingdom"
)

var (
	PingdomUser         = os.Getenv("PINGDOM_USER")
	PingdomPassword     = os.Getenv("PINGDOM_PASSWORD")
	PingdomAPIKey       = os.Getenv("PINGDOM_API_KEY")
	PingdomAccountEmail = os.Getenv("PINGDOM_ACCOUNT_EMAIL")
	RegionMap           = map[string]string{"region: EU": "aws:eu-west-2",
		"region: APAC": "aws:ap-southeast-2", "region: LATAM": "aws:us-west-2", "region: NA": "aws:us-east-2"}
	defaultAssertion = datadog.SyntheticsAssertion{Operator: Stringp("is"), Type: Stringp("statusCode"), Target: 200}
)

// ListPingdomChecks lists checks from pingdom - this is only a lightweight list, and does not contain full check details
// if full check details are desired, use the GetPingdomChecks function.
func ListPingdomChecks() (checks []pingdom.CheckResponse, e error) {

	client, _ := pingdom.NewClientWithConfig(pingdom.ClientConfig{User: PingdomUser, Password: PingdomPassword, APIKey: PingdomAPIKey, AccountEmail: PingdomAccountEmail})
	c, err := client.Checks.List()
	if err != nil {
		return checks, fmt.Errorf("Could not get checks: %s\n", err)
	}
	return c, nil
}

func GetPingdomChecks(ids []int) []pingdom.CheckResponse {
	var checks []pingdom.CheckResponse
	var wg sync.WaitGroup
	checkDetailChan := make(chan pingdom.CheckResponse, len(ids))
	for _, id := range ids {
		wg.Add(1)
		go func(checkID int, x *sync.WaitGroup) {
			defer x.Done()
			check, err := GetPingdomCheckDetail(checkID)
			if err != nil {
				log.Printf("failed to retrieve check %d, error: %s", checkID, err)
				return
			}
			checkDetailChan <- check
		}(id, &wg)
	}
	wg.Wait()
	close(checkDetailChan)
	for c := range checkDetailChan {
		checks = append(checks, c)
	}
	return checks
}

func GetPingdomCheck(checkID int) (checks []pingdom.CheckResponse, e error) {
	check, err := GetPingdomCheckDetail(checkID)
	if err != nil {
		return checks, fmt.Errorf("Could not get check %d: %s\n", checkID, err)
	}
	checks = append(checks, check)
	return checks, nil
}

func GetPingdomCheckDetail(checkID int) (pingdom.CheckResponse, error) {
	client, _ := pingdom.NewClientWithConfig(pingdom.ClientConfig{User: PingdomUser, Password: PingdomPassword, APIKey: PingdomAPIKey, AccountEmail: PingdomAccountEmail})
	check, err := client.Checks.Read(checkID)
	if err != nil {
		return *check, fmt.Errorf("Could not get check %d: %s\n", checkID, err)
	}
	return *check, nil
}

func DownloadCheck(checks []pingdom.CheckResponse) error {
	fmt.Printf("Found %d checks to download\n", len(checks))
	for _, check := range checks {
		fmt.Printf("Downloading check %d\n", check.ID)
		synthCheck, err := pingdomToSynthetics(check)
		if err != nil {
			log.Printf("could not download Pingdom check %d, encoutered error %s", check.ID, err)
		}
		filePath, err := writeCheckFile(synthCheck)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Printf("file(s) %s written to disk\n", filePath)
	}
	return nil
}

func CheckJson(checks []pingdom.CheckResponse) {
	j, err := json.Marshal(checks)
	if err != nil {
		log.Println("could not covert checks to json")
	}
	fmt.Println(string(j))
}

func writeCheckFile(check datadog.SyntheticsTest) (string, error) {
	filePath := fmt.Sprintf("./data/api-%s.yml", getHostName(*check.Config.Request.Url))
	jsonData, err := json.Marshal(check)
	yamDat, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return "", fmt.Errorf("could not marshal check data to Yaml due to error %s", err)
	}
	err = ioutil.WriteFile(filePath, yamDat, 0744)
	if err != nil {
		return "", fmt.Errorf("could not write %s to disk due to error %s", filePath, err)
	}
	return filePath, nil
}

//todo, check for errors here
func getHostName(path string) string {
	u, _ := url.Parse(path)
	host := u.Hostname
	return host()
}

func protocol(port int) string {
	if port == 443 {
		return "https"
	}
	return "http"
}

// Stringp returns a pointer to the string value passed in.
func Stringp(v string) *string {
	return &v
}

// Intp returns a pointer to the string value passed in.
func Intp(v int) *int {
	return &v
}

func convertTags(tags []pingdom.CheckResponseTag) ([]string, error) {
	var tagsSet []string
	for _, tag := range tags {
		tagsSet = append(tagsSet, tag.Name)
	}
	return tagsSet, nil
}

func pingdomToSynthetics(check pingdom.CheckResponse) (datadog.SyntheticsTest, error) {
	var s datadog.SyntheticsTest
	s.Type = Stringp("api")
	so := datadog.SyntheticsOptions{TickEvery: Intp(check.Resolution * 60), MinFailureDuration: Intp(60 * check.SendNotificationWhenDown)}
	headers := map[string]string(check.Type.HTTP.RequestHeaders)
	uri := fmt.Sprintf("%s://%s%s", protocol(check.Type.HTTP.Port), check.Hostname, check.Type.HTTP.Url)
	tags, err := convertTags(check.Tags)
	if err != nil {
		return s, fmt.Errorf("could not process tags from Pingdom check. %s", err)
	}
	s.Name = &check.Name
	s.Tags = tags
	s.Locations = getLocations(check.ProbeFilters)
	//todo, translate assertions from the shouldcontain logic in the pingdom response
	sr := datadog.SyntheticsRequest{Url: &uri,
		Method:  Stringp("GET"),
		Timeout: Intp(check.ResponseTimeThreshold / 1000),
		Headers: headers,
	}
	s.Message = Stringp(fmt.Sprintf("The %s did not respond with expected data", check.Name))
	s.Options = &so
	s.Status = Stringp("paused")
	sc := datadog.SyntheticsConfig{Request: &sr, Assertions: []datadog.SyntheticsAssertion{defaultAssertion}}
	s.Config = &sc
	return s, nil
}

//fmt.Printf("%d | %s | %s\n", check.ID, check.Name, check.Status)
func getLocations(probeFilters []string) []string {
	//when no region is defined we set a pair of default regions, this is in line with the Pingdom default
	if len(probeFilters) == 0 {
		return []string{"aws:us-east-2", "aws:eu-west-2"}
	}
	var result []string
	for _, probe := range probeFilters {
		result = append(result, getRegion(probe))
	}
	return result
}

func getRegion(pingdomRegion string) string {
	reg, ok := RegionMap[pingdomRegion]
	if !ok {
		log.Printf("Pingdom Region %s does not match expected region set, defaulted to aws:us-east-2", pingdomRegion)
		return "aws:us-east-2"
	}
	return reg
}

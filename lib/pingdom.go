package lib

import (
	"fmt"
	"os"

	"github.com/russellcardullo/go-pingdom/pingdom"
)

var (
	PingdomUser         = os.Getenv("PINGDOM_USER")
	PingdomPassword     = os.Getenv("PINGDOM_PASSWORD")
	PingdomAPIKey       = os.Getenv("PINGDOM_API_KEY")
	PingdomAccountEmail = os.Getenv("PINGDOM_ACCOUNT_EMAIL")
)

func ListPingdomChecks() {
	client := pingdom.NewMultiUserClient(PingdomUser, PingdomPassword, PingdomAPIKey, PingdomAccountEmail)
	checks, err := client.Checks.List()
	if err != nil {
		fmt.Printf("Could not get checks: %s", err)
	}
	fmt.Printf("%+v", checks[0])
}

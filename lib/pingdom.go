package lib

import (
	"fmt"
	"os"
<<<<<<< HEAD
	"strconv"
=======
>>>>>>> edbe8078564d4d452eb18b9a16ff1d74074d1bca

	"github.com/russellcardullo/go-pingdom/pingdom"
)

var (
	PingdomUser         = os.Getenv("PINGDOM_USER")
	PingdomPassword     = os.Getenv("PINGDOM_PASSWORD")
	PingdomAPIKey       = os.Getenv("PINGDOM_API_KEY")
	PingdomAccountEmail = os.Getenv("PINGDOM_ACCOUNT_EMAIL")
)

<<<<<<< HEAD
// ListPingdomChecks lists checks from pingdom
=======
>>>>>>> edbe8078564d4d452eb18b9a16ff1d74074d1bca
func ListPingdomChecks() {
	client := pingdom.NewMultiUserClient(PingdomUser, PingdomPassword, PingdomAPIKey, PingdomAccountEmail)
	checks, err := client.Checks.List()
	if err != nil {
		fmt.Printf("Could not get checks: %s", err)
	}
<<<<<<< HEAD
	for _, check := range checks {
		fmt.Printf("%d | %s | %s\n", check.ID, check.Name, check.Status)
	}
}

func GetPingdomCheck(id string) {
	client := pingdom.NewMultiUserClient(PingdomUser, PingdomPassword, PingdomAPIKey, PingdomAccountEmail)
	checkID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("Check ID can only be an integer.")
		return
	}
	check, err := client.Checks.Read(checkID)
	if err != nil {
		fmt.Printf("Could not get check %s: %s", id, err)
	}
	fmt.Printf("%d | %s | %s\n", check.ID, check.Name, check.Status)
=======
	fmt.Printf("%+v", checks[0])
>>>>>>> edbe8078564d4d452eb18b9a16ff1d74074d1bca
}

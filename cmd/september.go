package cmd

import (
	"fmt"
	"math"
	"time"

	"github.com/binaryfigments/crtsh"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(septemberCmd)
	septemberCmd.Flags().String("domain", "", "the domain to look for.")
	septemberCmd.Flags().Int("timeout", 10, "the time-out to use.")

	septemberCmd.MarkFlagRequired("domain")
}

var septemberCmd = &cobra.Command{
	Use:   "september",
	Short: "Get certificates that are valid 398 days after september.",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		timeout, _ := cmd.Flags().GetInt("timeout")
		timeout2 := time.Duration(timeout) * time.Second

		certs := crtsh.Get(domain, timeout2, 30)
		if certs.Error == true {
			fmt.Println(certs.ErrorMessage)
		}

		for _, cert := range certs.Certificates {
			// Skip expired certificates.
			if cert.Expired == true {
				continue
			}

			// 1 september 2020
			// 398 days max
			dateSept := setDate(2020, 9, 1)
			dateMax := dateSept.AddDate(0, 0, 398)

			if cert.NotAfter.Before(dateMax) {
				continue
			}

			// Get first certificate
			first := cert.NameValue[0]

			// Date calcs
			currentTime := time.Now()
			days := cert.NotAfter.Sub(currentTime).Hours() / 24
			days = math.Floor(days)

			// Days after september + 398 days
			days2 := cert.NotAfter.Sub(dateMax).Hours() / 24
			days2 = math.Floor(days2)

			// fmt.Println(first)
			// fmt.Println(cert.NotAfter.String())
			// fmt.Println(days)
			// fmt.Println(days2)
			// fmt.Println(first, ",", cert.NotAfter.String(), ",", days, ",", days2)

			fmt.Printf("%s, %s, %g, %g\n", first, cert.NotAfter.String(), days, days2)
		}
	},
}

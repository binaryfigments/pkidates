package cmd

import (
	"fmt"
	"math"
	"time"

	"github.com/binaryfigments/crtsh"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(expiringCmd)
	expiringCmd.Flags().String("domain", "", "the domain to look for.")
	expiringCmd.Flags().Int("days", 30, "days before expiring.")
	expiringCmd.Flags().Int("timeout", 10, "the time-out to use.")
	expiringCmd.Flags().Bool("verify", false, "verify with the online certificate, not working with wildcards.")

	expiringCmd.MarkFlagRequired("domain")
}

var expiringCmd = &cobra.Command{
	Use:   "expiring",
	Short: "Get get certificates that are expiring in XX days.",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		days, _ := cmd.Flags().GetInt("days")
		timeout, _ := cmd.Flags().GetInt("timeout")
		timeout2 := time.Duration(timeout) * time.Second

		certs := crtsh.Get(domain, timeout2, days)
		if certs.Error == true {
			fmt.Println(certs.ErrorMessage)
		}

		for _, cert := range certs.Certificates {
			// Skip expired certificates.
			if cert.Expired == true {
				continue
			}

			// Only certs that are expirering within 30 days
			if cert.Replace == true {
				continue
			}

			// Get first certificate
			first := cert.NameValue[0]

			// Date calcs
			currentTime := time.Now()
			days := cert.NotAfter.Sub(currentTime).Hours() / 24
			days = math.Floor(days)

			fmt.Printf("%s, %s, %g\n", first, cert.NotAfter.String(), days)

		}
	},
}

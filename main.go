package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	// Set up Viper and Cobra
	viper.SetDefault("prometheus.address", "http://localhost:9090")
	rootCmd := &cobra.Command{
		Use:   "pquery",
		Short: "Look up and print out Prometheus query results",
		Run: func(cmd *cobra.Command, args []string) {
			// Set up connection to Prometheus API
			address := viper.GetString("prometheus.address")
			client, err := prometheus.NewClient(prometheus.Config{
				Address: address,
			})
			if err != nil {
				log.Fatal(err)
			}

			// Check that a metric was provided as input
			if len(args) == 0 {
				log.Fatal("No metric provided")
			}
			metric := args[0]

			// Execute query to retrieve metrics data
			query := fmt.Sprintf("%s", metric)
			value, err := client.Query(context.Background(), query, time.Now())
			if err != nil {
				log.Fatal(err)
			}

			// Print query results
			fmt.Println("Query results:")
			for _, row := range value.Matrix.Values {
				fmt.Printf("  %v\n", row)
			}
		},
	}
	rootCmd.PersistentFlags().String("prometheus.address", "", "Address of the Prometheus server")
	viper.BindPFlag("prometheus.address", rootCmd.PersistentFlags().Lookup("prometheus.address"))

	// Parse command-line arguments and execute command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

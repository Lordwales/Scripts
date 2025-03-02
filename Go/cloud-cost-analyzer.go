package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
)

// AWS Cost Report Structure
type CostReport struct {
	Service string
	Cost    string
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Error loading AWS config: %v", err)
	}

	client := costexplorer.NewFromConfig(cfg)
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02") // Last 7 days
	endDate := time.Now().Format("2006-01-02")

	resp, err := client.GetCostAndUsage(context.TODO(), &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: &startDate,
			End:   &endDate,
		},
		Granularity: "DAILY",
		Metrics:     []string{"BlendedCost"},
		GroupBy: []costexplorer.GroupDefinition{
			{Type: "DIMENSION", Key: "SERVICE"},
		},
	})
	if err != nil {
		log.Fatalf("Error fetching cost data: %v", err)
	}

	var reports []CostReport
	for _, result := range resp.ResultsByTime {
		for _, group := range result.Groups {
			reports = append(reports, CostReport{
				Service: *group.Keys[0],
				Cost:    *group.Metrics["BlendedCost"].Amount,
			})
		}
	}

	// Convert to JSON
	output, _ := json.MarshalIndent(reports, "", "  ")
	fmt.Println(string(output))
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package stepFunctioncmd

import (
	"github.com/Appkube-awsx/awsx-stepFunctions/authenticater"
	"github.com/Appkube-awsx/awsx-stepFunctions/client"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/spf13/cobra"
	"log"
)

// getConfigDataCmd represents the getConfigData command
var GetCostDataCmd = &cobra.Command{
	Use:   "getCostData",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		vaultUrl := cmd.Parent().PersistentFlags().Lookup("vaultUrl").Value.String()
		accountNo := cmd.Parent().PersistentFlags().Lookup("accountId").Value.String()
		region := cmd.Parent().PersistentFlags().Lookup("zone").Value.String()
		acKey := cmd.Parent().PersistentFlags().Lookup("accessKey").Value.String()
		secKey := cmd.Parent().PersistentFlags().Lookup("secretKey").Value.String()
		crossAccountRoleArn := cmd.Parent().PersistentFlags().Lookup("crossAccountRoleArn").Value.String()
		externalId := cmd.Parent().PersistentFlags().Lookup("externalId").Value.String()

		authFlag := authenticater.AuthenticateData(vaultUrl, accountNo, region, acKey, secKey, crossAccountRoleArn, externalId)

		if authFlag {
			table, _ := cmd.Flags().GetString("table")
			getDynamoDbCostDetail(region, crossAccountRoleArn, acKey, secKey, table, externalId)
		}
	},
}

func getDynamoDbCostDetail(region string, crossAccountRoleArn string, accessKey string, secretKey string, table string, externalId string) (*costexplorer.GetCostAndUsageOutput, error) {
	log.Println("Getting cost data")
	costClient := client.GetCostClient(region, crossAccountRoleArn, accessKey, secretKey, externalId)

	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String("2023-01-01"),
			End:   aws.String("2023-02-01"),
		},
		Metrics: []*string{
			aws.String("UnblendedCost"),
		},
		Granularity: aws.String("MONTHLY"),
		Filter: &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key: aws.String("SERVICE"),
				Values: []*string{
					aws.String("Amazon DynamoDB"),
				},
			},
		},
	}

	costData, err := costClient.GetCostAndUsage(input)
	if err != nil {
		log.Fatalln("Error: in getting cost data", err)
	}

	log.Println(costData)
	return costData, err
}

func init() {

}

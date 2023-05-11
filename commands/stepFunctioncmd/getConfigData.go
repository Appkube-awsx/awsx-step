package stepFunctioncmd

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-stepFunctions/authenticater"
	"github.com/Appkube-awsx/awsx-stepFunctions/client"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/spf13/cobra"
)

// getConfigDataCmd represents the getConfigData command
var GetConfigDataCmd = &cobra.Command{
	Use:   "getConfigData",
	Short: "get detailed data about state machine",
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
			arn, _ := cmd.Flags().GetString("arn")
			getStateMachineDetail(region, crossAccountRoleArn, acKey, secKey, arn, externalId)
		}
	},
}

func getStateMachineDetail(region string, crossAccountRoleArn string, accessKey string, secretKey string, arn string, externalId string) (*sfn.DescribeStateMachineOutput, error) {
	log.Println("Getting state machine data")
	sfnClient := client.GetClient(region, crossAccountRoleArn, accessKey, secretKey, externalId)

	input := &sfn.DescribeStateMachineInput{
		StateMachineArn: aws.String(arn),
	}

	sfnMachineData, err := sfnClient.DescribeStateMachine(input)
	if err != nil {
		log.Fatalln("Error: in getting step machine data", err)
	}

	log.Println(sfnMachineData)
	return sfnMachineData, err
}

func init() {
	GetConfigDataCmd.Flags().StringP("arn", "a", "", "state machine arn")

	if err := GetConfigDataCmd.MarkFlagRequired("arn"); err != nil {
		fmt.Println(err)
	}
}

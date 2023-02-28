/*
Copyright © 2023 Manoj Sharma manoj.sharma@synectiks.com
*/
package commands

import (
	"github.com/Appkube-awsx/awsx-stepFunctions/authenticater"
	"github.com/Appkube-awsx/awsx-stepFunctions/client"
	"github.com/Appkube-awsx/awsx-stepFunctions/commands/stepFunctioncmd"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/spf13/cobra"
	"log"
)

// AwsxCloudElementsCmd represents the base command when called without any subcommands
var AwsxStateMachineCmd = &cobra.Command{
	Use:   "stepFunctions",
	Short: "get stepFunctions Details command gets resource counts",
	Long:  `get stepFunctions Details command gets resource counts details of an AWS account`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Command stepFunctions started")
		vaultUrl := cmd.PersistentFlags().Lookup("vaultUrl").Value.String()
		accountNo := cmd.PersistentFlags().Lookup("accountId").Value.String()
		region := cmd.PersistentFlags().Lookup("zone").Value.String()
		acKey := cmd.PersistentFlags().Lookup("accessKey").Value.String()
		secKey := cmd.PersistentFlags().Lookup("secretKey").Value.String()
		crossAccountRoleArn := cmd.PersistentFlags().Lookup("crossAccountRoleArn").Value.String()
		externalId := cmd.PersistentFlags().Lookup("externalId").Value.String()

		authFlag := authenticater.AuthenticateData(vaultUrl, accountNo, region, acKey, secKey, crossAccountRoleArn, externalId)

		if authFlag {
			getStateMachineList(region, crossAccountRoleArn, acKey, secKey, externalId)
		}
	},
}

func getStateMachineList(region string, crossAccountRoleArn string, accessKey string, secretKey string, externalId string) (*sfn.ListStateMachinesOutput, error) {
	log.Println("Getting stepFunctions list summary")
	sfnClient := client.GetClient(region, crossAccountRoleArn, accessKey, secretKey, externalId)

	input := &sfn.ListStateMachinesInput{}
	machineList, err := sfnClient.ListStateMachines(input)
	if err != nil {
		log.Fatalln("Error: in getting state machine list", err)
	}
	log.Println(machineList)
	return machineList, err
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := AwsxStateMachineCmd.Execute()
	if err != nil {
		log.Fatal("There was some error while executing the CLI: ", err)
		return
	}
}

func init() {
	AwsxStateMachineCmd.AddCommand(stepFunctioncmd.GetConfigDataCmd)
	AwsxStateMachineCmd.AddCommand(stepFunctioncmd.GetCostDataCmd)
	AwsxStateMachineCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxStateMachineCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxStateMachineCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxStateMachineCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxStateMachineCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxStateMachineCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxStateMachineCmd.PersistentFlags().String("externalId", "", "aws external id auth")
}

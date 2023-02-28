/*
Copyright © 2023 Manoj Sharma manoj.sharma@synectiks.com
*/
package cmd

import (
	"github.com/Appkube-awsx/awsx-getLambdas/client"
	"github.com/Appkube-awsx/awsx-getLambdas/vault"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/cobra"
	"log"
)

// AwsxCloudElementsCmd represents the base command when called without any subcommands
var AwsxCloudElementsCmd = &cobra.Command{
	Use:   "getLambdas",
	Short: "getLambdas command gets resource counts",
	Long:  `getLambdas command gets resource counts details of an AWS account`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Command getElementDetails started")
		vaultUrl, _ := cmd.Flags().GetString("vaultUrl")
		accountNo, _ := cmd.Flags().GetString("accountId")
		region, _ := cmd.Flags().GetString("zone")
		acKey, _ := cmd.Flags().GetString("accessKey")
		secKey, _ := cmd.Flags().GetString("secretKey")
		crossAccountRoleArn, _ := cmd.Flags().GetString("crossAccountRoleArn")

		if vaultUrl != "" && accountNo != "" {
			if region == "" {
				log.Fatalln("Zone not provided. Program exit")
				return
			}
			log.Println("Getting account details")
			data, err := vault.GetAccountDetails(vaultUrl, accountNo)
			if err != nil {
				log.Println("Error in calling the account details api. \n", err)
				return
			}
			if data.AccessKey == "" || data.SecretKey == "" || data.CrossAccountRoleArn == "" {
				log.Println("Account details not found.")
				return
			}
			getLambdaList(region, data.CrossAccountRoleArn, data.AccessKey, data.SecretKey)
		} else if region != "" && acKey != "" && secKey != "" {
			getLambdaList(region, crossAccountRoleArn, acKey, secKey)
		} else {
			log.Fatal("AWS credentials like accesskey/secretkey/region/crossAccountRoleArn not provided. Program exit")
			return
		}
	},
}

func getLambdaList(region string, crossAccountRoleArn string, accessKey string, secretKey string) (*lambda.ListFunctionsOutput, error) {
	log.Println("Getting lambda list summary")
	lambdaClient := client.GetClient(region, crossAccountRoleArn, accessKey, secretKey)

	input := &lambda.ListFunctionsInput{}
	lambdaList, err := lambdaClient.ListFunctions(input)
	if err != nil {
		log.Fatalln("Error: in getting lambda list", err)
	}
	log.Println(lambdaList)
	return lambdaList, err
}

//func GetConfig(region string, crossAccountRoleArn string, accessKey string, secretKey string) *configservice.GetDiscoveredResourceCountsOutput {
//	return getLambdaList(region, crossAccountRoleArn, accessKey, secretKey)
//}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := AwsxCloudElementsCmd.Execute()
	if err != nil {
		log.Fatal("There was some error while executing the CLI: ", err)
		return
	}
}

func init() {
	AwsxCloudElementsCmd.Flags().String("vaultUrl", "", "vault end point")
	AwsxCloudElementsCmd.Flags().String("accountId", "", "aws account number")
	AwsxCloudElementsCmd.Flags().String("zone", "", "aws region")
	AwsxCloudElementsCmd.Flags().String("accessKey", "", "aws access key")
	AwsxCloudElementsCmd.Flags().String("secretKey", "", "aws secret key")
	AwsxCloudElementsCmd.Flags().String("crossAccountRoleArn", "", "aws cross account role arn")
}

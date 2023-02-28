package client

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-getLambdas/awssession"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sts"
	"log"
)

func GetClient(region string, crossAccountRoleArn string, accessKey string, secretKey string) *lambda.Lambda {

	roleArn := crossAccountRoleArn
	sessionName := "assume_role_session_name"
	externalId := "DJ6@a8hzG@xkFwSvLmkSR5SN"

	return assumeRole(roleArn, sessionName, externalId, accessKey, secretKey, region)
}

func assumeRole(roleArn string, sessionName string, externalId string, accesskey string, secretKey string, region string) *lambda.Lambda {
	sess, err := awssession.GetSessionByCreds(region, accesskey, secretKey, "")

	if err != nil {
		fmt.Printf("failed to create aws session, %v\n", err)
		log.Fatal(err)
	}

	svc := sts.New(sess)

	assumeRoleInput := sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(sessionName),
		DurationSeconds: aws.Int64(60 * 60 * 1),
	}

	if externalId != "nil" {
		fmt.Println("Trying to fetch external id to assume new role")
		assumeRoleInput.ExternalId = aws.String(externalId)
	}

	result, err := svc.AssumeRole(&assumeRoleInput)

	if err != nil {
		fmt.Printf("failed to assume role, %v\n", err)
		log.Fatal(err)
	}
	fmt.Println("Assume role output: ", result)

	awsSession, err := awssession.GetSessionByCreds("us-east-1", *result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken)

	if err != nil {
		fmt.Printf("failed to assume role, %v\n", err)
		log.Fatal(err)
	}

	lambdaClient := lambda.New(awsSession)

	return lambdaClient
}

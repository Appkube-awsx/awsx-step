package client

import (
	"fmt"

	"github.com/Appkube-awsx/awsx-stepFunctions/awssession"
	"github.com/aws/aws-sdk-go/service/sfn"

	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
)

func GetClient(region string, crossAccountRoleArn string, accessKey string, secretKey string, externalId string) *sfn.SFN {

	sessionName := "assume_role_session_name"
	return assumeRole(crossAccountRoleArn, sessionName, externalId, accessKey, secretKey, region)
}


func assumeRole(roleArn string, sessionName string, externalId string, accesskey string, secretKey string, region string) *sfn.SFN {
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
	//fmt.Println("Assume role output: ", result)

	awsSession, err := awssession.GetSessionByCreds("us-east-1", *result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken)

	if err != nil {
		fmt.Printf("failed to assume role, %v\n", err)
		log.Fatal(err)
	}

	sfnClient := sfn.New(awsSession)

	return sfnClient
}

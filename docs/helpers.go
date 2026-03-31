package infra

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/tidwall/gjson"
)

func GetAWSRegion() string {
	if os.Getenv("AWS_REGION")!= "" {
		return os.Getenv("AWS_REGION")
	}
	return "us-west-2"
}

func GetAWSSTSClient() *sts.STS {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(GetAWSRegion())}, nil)
	if err!= nil {
		panic(err)
	}
	return sts.New(sess)
}

func GetAWSAccountID() string {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(GetAWSRegion())}, nil))
	stsClient := sts.New(sess)
	resp, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err!= nil {
		panic(err)
	}
	return *resp.Account
}

func GetTerraformStateFile() string {
	stateFilePath := os.Getenv("TERRAFORM_STATE_FILE")
	if stateFilePath!= "" {
		return stateFilePath
	}
	return "terraform.tfstate"
}

func GetTerraformBackendConfig() map[string]string {
	backendConfig := make(map[string]string)
	stateFile := GetTerraformStateFile()
	backendConfig["backend"] = "s3"
	backendConfig["bucket"] = "my-terraform-state-bucket"
	backendConfig["key"] = stateFile
	backendConfig["region"] = GetAWSRegion()
	return backendConfig
}

func GetTerraformVariablesFromJSONFile(filePath string) map[string]interface{} {
	jsonFile, err := os.Open(filePath)
	if err!= nil {
		fmt.Println(err)
		return nil
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err!= nil {
		fmt.Println(err)
		return nil
	}
	var variables map[string]interface{}
	json.Unmarshal([]byte(jsonData), &variables)
	return variables
}

func GetTerraformVariablesFromJSONString(jsonString string) map[string]interface{} {
	var variables map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &variables)
	if err!= nil {
		fmt.Println(err)
		return nil
	}
	return variables
}

func GetTerraformPlanOutput(planFilePath string) string {
	planFile, err := os.Open(planFilePath)
	if err!= nil {
		fmt.Println(err)
		return ""
	}
	defer planFile.Close()
	jsonData, err := ioutil.ReadAll(planFile)
	if err!= nil {
		fmt.Println(err)
		return ""
	}
	var plan map[string]interface{}
	json.Unmarshal([]byte(jsonData), &plan)
	planOutput := gjson.Get(string(jsonData), "resource_changes.#").String()
	return planOutput
}
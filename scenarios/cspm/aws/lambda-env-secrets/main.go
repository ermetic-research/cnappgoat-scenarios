package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Define your secrets
		secret1 := "administrator"
		secret2 := "123156asd!@%!#^3a"

		// Create IAM role for the lambda function
		lambdaRole, err := iam.NewRole(ctx, "CNAPPgoat-lambda-role", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Action": "sts:AssumeRole",
					"Principal": {
						"Service": "lambda.amazonaws.com"
					},
					"Effect": "Allow",
					"Sid": ""
				}]
			}`),
		})
		if err != nil {
			return err
		}

		// Create Lambda function
		lambdaSecrets, err := lambda.NewFunction(ctx, "CNAPPgoat-lambda-env-secrets", &lambda.FunctionArgs{
			Handler: pulumi.String("index.handler"),
			Role:    lambdaRole.Arn,
			Runtime: pulumi.String("nodejs18.x"),
			Environment: &lambda.FunctionEnvironmentArgs{
				Variables: pulumi.StringMap{
					"username": pulumi.String(secret1),
					"password": pulumi.String(secret2),
				},
			},
			Code: pulumi.NewFileArchive("./app.zip"),
			Tags: pulumi.StringMap{
				"Name":      pulumi.String("CNAPPgoat-lambda-env-secrets"),
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("CNAPPgoat-lambda-role", lambdaRole.Arn)
		ctx.Export("CNAPPgoat-lambda-env-secrets", lambdaSecrets.Arn)
		return nil
	})
}

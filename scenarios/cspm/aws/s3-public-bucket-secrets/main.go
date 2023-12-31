package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a new S3 bucket
		bucket, err := s3.NewBucketV2(
			ctx,
			"CNAPPgoat-public-bucket",
			&s3.BucketV2Args{
				Tags: pulumi.StringMap{
					"Cnappgoat": pulumi.String("true"),
				},
			})

		if err != nil {
			return err
		}
		publicAccessBlock, err := s3.NewBucketPublicAccessBlock(ctx, "CNAPPgoat-public-bucket-access-block", &s3.BucketPublicAccessBlockArgs{
			Bucket:                bucket.ID(),
			BlockPublicAcls:       pulumi.Bool(false),
			BlockPublicPolicy:     pulumi.Bool(false),
			IgnorePublicAcls:      pulumi.Bool(false),
			RestrictPublicBuckets: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		_, err = s3.NewBucketPolicy(ctx, "CNAPPgoat-public-bucket-policy", &s3.BucketPolicyArgs{
			Bucket: bucket.ID(),
			Policy: bucket.ID().ApplyT(func(id pulumi.String) (pulumi.String, error) {
				return pulumi.String(fmt.Sprintf(`{
                    "Version": "2012-10-17",
                    "Statement": [
                      {
                        "Effect": "Allow",
                        "Principal": "*",
                        "Action": "s3:GetObject",
                        "Resource": "arn:aws:s3:::%s/*"
                      }
                    ]
                  }`, id)), nil
			}).(pulumi.StringOutput),
			// depends on the public bucket access block
		}, pulumi.DependsOn([]pulumi.Resource{publicAccessBlock}))
		if err != nil {
			return err
		}

		// Upload a secret file to the bucket
		bucketObject, err := s3.NewBucketObject(ctx, "CNAPPgoat-public-bucket-secret-object", &s3.BucketObjectArgs{
			Bucket:      bucket.ID(),
			Key:         pulumi.String("CNAPPgoatSecret"),
			Source:      pulumi.NewFileAsset("secret.txt"),
			ContentType: pulumi.String("text/plain"),
		})
		if err != nil {
			return err
		}
		ctx.Export("CNAPPgoat-public-bucket", bucket.Arn)
		ctx.Export("object-key", bucketObject.Key)
		return nil
	})
}

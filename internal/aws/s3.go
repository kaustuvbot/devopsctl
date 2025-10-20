package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
)

// CheckS3PublicBuckets identifies S3 buckets that are publicly accessible.
// Severity: CRITICAL
func CheckS3PublicBuckets(ctx context.Context, client S3Client) ([]reporter.CheckResult, error) {
	var results []reporter.CheckResult

	bucketsOut, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		if isPermissionError(err) {
			return results, nil
		}
		return nil, fmt.Errorf("ListBuckets: %w", err)
	}

	for _, bucket := range bucketsOut.Buckets {
		name := *bucket.Name
		isPublic, err := isBucketPublic(ctx, client, name)
		if err != nil {
			continue
		}
		if isPublic {
			results = append(results, reporter.CheckResult{
				CheckName:      "s3-public-bucket",
				Severity:       "CRITICAL",
				ResourceID:     name,
				Message:        fmt.Sprintf("S3 bucket %q is publicly accessible", name),
				Recommendation: "Enable S3 Block Public Access settings for the bucket and account",
			})
		}
	}
	return results, nil
}

// CheckS3Encryption checks for S3 buckets without server-side encryption configured.
// Severity: HIGH
func CheckS3Encryption(ctx context.Context, client S3Client) ([]reporter.CheckResult, error) {
	var results []reporter.CheckResult

	bucketsOut, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		if isPermissionError(err) {
			return results, nil
		}
		return nil, fmt.Errorf("ListBuckets: %w", err)
	}

	for _, bucket := range bucketsOut.Buckets {
		name := *bucket.Name
		_, err := client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{Bucket: &name})
		if err != nil {
			if isNoSuchBucket(err) || isPermissionError(err) {
				continue
			}
			// NoSuchServerSideEncryptionConfiguration means no encryption configured
			results = append(results, reporter.CheckResult{
				CheckName:      "s3-no-encryption",
				Severity:       "HIGH",
				ResourceID:     name,
				Message:        fmt.Sprintf("S3 bucket %q has no server-side encryption configured", name),
				Recommendation: "Enable SSE-S3 or SSE-KMS encryption on the bucket",
			})
		}
	}
	return results, nil
}

// CheckS3Versioning checks for S3 buckets without versioning enabled.
// Severity: LOW
func CheckS3Versioning(ctx context.Context, client S3Client) ([]reporter.CheckResult, error) {
	var results []reporter.CheckResult

	bucketsOut, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		if isPermissionError(err) {
			return results, nil
		}
		return nil, fmt.Errorf("ListBuckets: %w", err)
	}

	for _, bucket := range bucketsOut.Buckets {
		name := *bucket.Name
		verOut, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{Bucket: &name})
		if err != nil {
			continue
		}
		if verOut.Status != s3types.BucketVersioningStatusEnabled {
			results = append(results, reporter.CheckResult{
				CheckName:      "s3-versioning-disabled",
				Severity:       "LOW",
				ResourceID:     name,
				Message:        fmt.Sprintf("S3 bucket %q does not have versioning enabled", name),
				Recommendation: "Enable versioning for data protection and point-in-time recovery",
			})
		}
	}
	return results, nil
}

func isBucketPublic(ctx context.Context, client S3Client, bucket string) (bool, error) {
	pabOut, err := client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{
		Bucket: &bucket,
	})
	if err == nil && pabOut.PublicAccessBlockConfiguration != nil {
		cfg := pabOut.PublicAccessBlockConfiguration
		if boolVal(cfg.BlockPublicAcls) &&
			boolVal(cfg.BlockPublicPolicy) &&
			boolVal(cfg.IgnorePublicAcls) &&
			boolVal(cfg.RestrictPublicBuckets) {
			return false, nil
		}
	}

	aclOut, err := client.GetBucketAcl(ctx, &s3.GetBucketAclInput{Bucket: &bucket})
	if err != nil {
		return false, err
	}
	for _, grant := range aclOut.Grants {
		if grant.Grantee != nil && grant.Grantee.Type == s3types.TypeGroup {
			uri := ""
			if grant.Grantee.URI != nil {
				uri = *grant.Grantee.URI
			}
			if strings.Contains(uri, "AllUsers") || strings.Contains(uri, "AuthenticatedUsers") {
				return true, nil
			}
		}
	}
	return false, nil
}

func boolVal(b *bool) bool { return b != nil && *b }

func isNoSuchBucket(err error) bool {
	return err != nil && strings.Contains(err.Error(), "NoSuchBucket")
}

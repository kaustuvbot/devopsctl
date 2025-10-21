package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type mockS3Client struct {
	listBucketsOutput          *s3.ListBucketsOutput
	listBucketsErr             error
	getBucketAclOutput         *s3.GetBucketAclOutput
	getBucketAclErr            error
	getPublicAccessBlockOutput *s3.GetPublicAccessBlockOutput
	getPublicAccessBlockErr    error
	getBucketEncryptionOutput  *s3.GetBucketEncryptionOutput
	getBucketEncryptionErr     error
	getBucketVersioningOutput  *s3.GetBucketVersioningOutput
	getBucketVersioningErr     error
}

func (m *mockS3Client) ListBuckets(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.listBucketsOutput, m.listBucketsErr
}
func (m *mockS3Client) GetBucketAcl(_ context.Context, _ *s3.GetBucketAclInput, _ ...func(*s3.Options)) (*s3.GetBucketAclOutput, error) {
	return m.getBucketAclOutput, m.getBucketAclErr
}
func (m *mockS3Client) GetPublicAccessBlock(_ context.Context, _ *s3.GetPublicAccessBlockInput, _ ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return m.getPublicAccessBlockOutput, m.getPublicAccessBlockErr
}
func (m *mockS3Client) GetBucketEncryption(_ context.Context, _ *s3.GetBucketEncryptionInput, _ ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return m.getBucketEncryptionOutput, m.getBucketEncryptionErr
}
func (m *mockS3Client) GetBucketVersioning(_ context.Context, _ *s3.GetBucketVersioningInput, _ ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return m.getBucketVersioningOutput, m.getBucketVersioningErr
}

func TestCheckS3PublicBuckets_PublicACL(t *testing.T) {
	allUsersURI := "http://acs.amazonaws.com/groups/global/AllUsers"
	mock := &mockS3Client{
		listBucketsOutput:       &s3.ListBucketsOutput{Buckets: []s3types.Bucket{{Name: aws.String("public-bucket")}}},
		getPublicAccessBlockErr: fmt.Errorf("NoSuchPublicAccessBlockConfiguration"),
		getBucketAclOutput: &s3.GetBucketAclOutput{
			Grants: []s3types.Grant{{
				Grantee: &s3types.Grantee{Type: s3types.TypeGroup, URI: &allUsersURI},
			}},
		},
	}
	results, err := CheckS3PublicBuckets(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Severity != "CRITICAL" {
		t.Errorf("expected CRITICAL for public bucket, got %v", results)
	}
}

func TestCheckS3PublicBuckets_BlockEnabled(t *testing.T) {
	mock := &mockS3Client{
		listBucketsOutput: &s3.ListBucketsOutput{Buckets: []s3types.Bucket{{Name: aws.String("safe-bucket")}}},
		getPublicAccessBlockOutput: &s3.GetPublicAccessBlockOutput{
			PublicAccessBlockConfiguration: &s3types.PublicAccessBlockConfiguration{
				BlockPublicAcls:       aws.Bool(true),
				BlockPublicPolicy:     aws.Bool(true),
				IgnorePublicAcls:      aws.Bool(true),
				RestrictPublicBuckets: aws.Bool(true),
			},
		},
	}
	results, err := CheckS3PublicBuckets(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for blocked bucket, got %d", len(results))
	}
}

func TestCheckS3Encryption_Missing(t *testing.T) {
	mock := &mockS3Client{
		listBucketsOutput:      &s3.ListBucketsOutput{Buckets: []s3types.Bucket{{Name: aws.String("unencrypted")}}},
		getBucketEncryptionErr: fmt.Errorf("ServerSideEncryptionConfigurationNotFoundError"),
	}
	results, err := CheckS3Encryption(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Severity != "HIGH" {
		t.Errorf("expected HIGH for missing encryption, got %v", results)
	}
}

func TestCheckS3Encryption_Present(t *testing.T) {
	mock := &mockS3Client{
		listBucketsOutput:     &s3.ListBucketsOutput{Buckets: []s3types.Bucket{{Name: aws.String("encrypted")}}},
		getBucketEncryptionOutput: &s3.GetBucketEncryptionOutput{},
	}
	results, err := CheckS3Encryption(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for encrypted bucket, got %d", len(results))
	}
}

func TestCheckS3Versioning_Disabled(t *testing.T) {
	mock := &mockS3Client{
		listBucketsOutput:         &s3.ListBucketsOutput{Buckets: []s3types.Bucket{{Name: aws.String("no-versioning")}}},
		getBucketVersioningOutput: &s3.GetBucketVersioningOutput{Status: s3types.BucketVersioningStatusSuspended},
	}
	results, err := CheckS3Versioning(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Severity != "LOW" {
		t.Errorf("expected LOW for disabled versioning, got %v", results)
	}
}

func TestCheckS3Versioning_Enabled(t *testing.T) {
	mock := &mockS3Client{
		listBucketsOutput:         &s3.ListBucketsOutput{Buckets: []s3types.Bucket{{Name: aws.String("versioned")}}},
		getBucketVersioningOutput: &s3.GetBucketVersioningOutput{Status: s3types.BucketVersioningStatusEnabled},
	}
	results, err := CheckS3Versioning(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for versioned bucket, got %d", len(results))
	}
}

package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestCheckEBSEncryption_Unencrypted(t *testing.T) {
	mock := &mockEC2Client{
		describeVolumesOutput: &ec2.DescribeVolumesOutput{
			Volumes: []ec2types.Volume{{
				VolumeId:  aws.String("vol-0123456789"),
				Encrypted: aws.Bool(false),
				State:     ec2types.VolumeStateInUse,
			}},
		},
	}
	results, err := CheckEBSEncryption(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Severity != "HIGH" {
		t.Errorf("expected HIGH for unencrypted volume, got %v", results)
	}
}

func TestCheckEBSEncryption_Encrypted(t *testing.T) {
	mock := &mockEC2Client{
		describeVolumesOutput: &ec2.DescribeVolumesOutput{
			Volumes: []ec2types.Volume{{
				VolumeId:  aws.String("vol-encrypted"),
				Encrypted: aws.Bool(true),
				State:     ec2types.VolumeStateInUse,
			}},
		},
	}
	results, err := CheckEBSEncryption(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for encrypted volume, got %d", len(results))
	}
}

func TestCheckEBSUnattached_Available(t *testing.T) {
	mock := &mockEC2Client{
		describeVolumesOutput: &ec2.DescribeVolumesOutput{
			Volumes: []ec2types.Volume{{
				VolumeId:  aws.String("vol-orphan"),
				Encrypted: aws.Bool(true),
				State:     ec2types.VolumeStateAvailable,
			}},
		},
	}
	results, err := CheckEBSUnattached(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Severity != "LOW" {
		t.Errorf("expected LOW for unattached volume, got %v", results)
	}
}

func TestCheckEBSUnattached_InUse(t *testing.T) {
	mock := &mockEC2Client{
		describeVolumesOutput: &ec2.DescribeVolumesOutput{
			Volumes: []ec2types.Volume{{
				VolumeId:  aws.String("vol-in-use"),
				Encrypted: aws.Bool(true),
				State:     ec2types.VolumeStateInUse,
			}},
		},
	}
	results, _ := CheckEBSUnattached(context.Background(), mock)
	if len(results) != 0 {
		t.Errorf("expected no results for in-use volume, got %d", len(results))
	}
}

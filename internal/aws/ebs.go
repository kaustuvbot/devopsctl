package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/kaustuvbot/devopsctl/internal/reporter"
)

// CheckEBSEncryption checks for unencrypted EBS volumes.
// Severity: HIGH
func CheckEBSEncryption(ctx context.Context, client EC2Client) ([]reporter.CheckResult, error) {
	var results []reporter.CheckResult

	vols, err := describeAllVolumes(ctx, client)
	if err != nil {
		return nil, err
	}
	for _, vol := range vols {
		if vol.Encrypted == nil || !*vol.Encrypted {
			results = append(results, reporter.CheckResult{
				CheckName:      "ebs-unencrypted",
				Severity:       "HIGH",
				ResourceID:     *vol.VolumeId,
				Message:        fmt.Sprintf("EBS volume %q is not encrypted", *vol.VolumeId),
				Recommendation: "Enable EBS encryption by default in your AWS account settings",
			})
		}
	}
	return results, nil
}

// CheckEBSUnattached checks for EBS volumes not attached to any instance.
// Severity: LOW
func CheckEBSUnattached(ctx context.Context, client EC2Client) ([]reporter.CheckResult, error) {
	var results []reporter.CheckResult

	vols, err := describeAllVolumes(ctx, client)
	if err != nil {
		return nil, err
	}
	for _, vol := range vols {
		if vol.State == ec2types.VolumeStateAvailable {
			results = append(results, reporter.CheckResult{
				CheckName:      "ebs-unattached",
				Severity:       "LOW",
				ResourceID:     *vol.VolumeId,
				Message:        fmt.Sprintf("EBS volume %q is not attached to any instance", *vol.VolumeId),
				Recommendation: "Delete unused EBS volumes to reduce costs",
			})
		}
	}
	return results, nil
}

func describeAllVolumes(ctx context.Context, client EC2Client) ([]ec2types.Volume, error) {
	out, err := client.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{})
	if err != nil {
		if isPermissionError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("DescribeVolumes: %w", err)
	}
	return out.Volumes, nil
}

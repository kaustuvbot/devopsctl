package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type mockEC2Client struct {
	describeSecurityGroupsOutput *ec2.DescribeSecurityGroupsOutput
	describeSecurityGroupsErr    error
	describeVolumesOutput        *ec2.DescribeVolumesOutput
	describeVolumesErr           error
}

func (m *mockEC2Client) DescribeSecurityGroups(_ context.Context, _ *ec2.DescribeSecurityGroupsInput, _ ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error) {
	return m.describeSecurityGroupsOutput, m.describeSecurityGroupsErr
}
func (m *mockEC2Client) DescribeVolumes(_ context.Context, _ *ec2.DescribeVolumesInput, _ ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error) {
	return m.describeVolumesOutput, m.describeVolumesErr
}

func TestCheckSecurityGroups_SSHOpen(t *testing.T) {
	fromPort := int32(22)
	toPort := int32(22)
	cidr := "0.0.0.0/0"
	proto := "tcp"
	mock := &mockEC2Client{
		describeSecurityGroupsOutput: &ec2.DescribeSecurityGroupsOutput{
			SecurityGroups: []ec2types.SecurityGroup{{
				GroupId:   aws.String("sg-12345"),
				GroupName: aws.String("web-sg"),
				IpPermissions: []ec2types.IpPermission{{
					FromPort:   &fromPort,
					ToPort:     &toPort,
					IpProtocol: &proto,
					IpRanges:   []ec2types.IpRange{{CidrIp: &cidr}},
				}},
			}},
		},
	}
	results, err := CheckSecurityGroups(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Severity != "CRITICAL" {
		t.Errorf("expected CRITICAL for SSH open, got %v", results)
	}
	if results[0].CheckName != "sg-ssh-open" {
		t.Errorf("expected check name 'sg-ssh-open', got %s", results[0].CheckName)
	}
}

func TestCheckSecurityGroups_AllTraffic(t *testing.T) {
	cidr := "0.0.0.0/0"
	proto := "-1"
	mock := &mockEC2Client{
		describeSecurityGroupsOutput: &ec2.DescribeSecurityGroupsOutput{
			SecurityGroups: []ec2types.SecurityGroup{{
				GroupId:   aws.String("sg-99999"),
				GroupName: aws.String("open-sg"),
				IpPermissions: []ec2types.IpPermission{{
					IpProtocol: &proto,
					IpRanges:   []ec2types.IpRange{{CidrIp: &cidr}},
				}},
			}},
		},
	}
	results, err := CheckSecurityGroups(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].CheckName != "sg-all-ports-open" {
		t.Errorf("expected sg-all-ports-open result, got %v", results)
	}
}

func TestCheckSecurityGroups_NoPublicAccess(t *testing.T) {
	fromPort := int32(443)
	toPort := int32(443)
	cidr := "10.0.0.0/8"
	proto := "tcp"
	mock := &mockEC2Client{
		describeSecurityGroupsOutput: &ec2.DescribeSecurityGroupsOutput{
			SecurityGroups: []ec2types.SecurityGroup{{
				GroupId:   aws.String("sg-private"),
				GroupName: aws.String("private-sg"),
				IpPermissions: []ec2types.IpPermission{{
					FromPort:   &fromPort,
					ToPort:     &toPort,
					IpProtocol: &proto,
					IpRanges:   []ec2types.IpRange{{CidrIp: &cidr}},
				}},
			}},
		},
	}
	results, _ := CheckSecurityGroups(context.Background(), mock)
	if len(results) != 0 {
		t.Errorf("expected no results for private CIDR, got %d", len(results))
	}
}

func TestCheckSecurityGroups_PermissionError(t *testing.T) {
	mock := &mockEC2Client{
		describeSecurityGroupsErr: fmt.Errorf("AccessDenied: User is not authorized to perform ec2:DescribeSecurityGroups"),
	}
	results, err := CheckSecurityGroups(context.Background(), mock)
	// Permission errors should be handled gracefully, returning empty results
	if err != nil {
		t.Errorf("expected no error for permission denied, got %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results on permission error, got %d", len(results))
	}
}

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/kaustuvbot/devopsctl/internal/reporter"
)

// CheckSecurityGroups checks for overly permissive security group rules.
// Severity: CRITICAL for port 22 open to 0.0.0.0/0 or all-traffic rules.
func CheckSecurityGroups(ctx context.Context, client EC2Client) ([]reporter.CheckResult, error) {
	var results []reporter.CheckResult

	sgsOut, err := client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		if isPermissionError(err) {
			return results, nil
		}
		return nil, fmt.Errorf("DescribeSecurityGroups: %w", err)
	}

	for _, sg := range sgsOut.SecurityGroups {
		sgID := *sg.GroupId
		sgName := ""
		if sg.GroupName != nil {
			sgName = *sg.GroupName
		}
		for _, perm := range sg.IpPermissions {
			for _, ipRange := range perm.IpRanges {
				if ipRange.CidrIp == nil || *ipRange.CidrIp != "0.0.0.0/0" {
					continue
				}
				// All-traffic rule: protocol -1 means all
				if perm.IpProtocol != nil && *perm.IpProtocol == "-1" {
					results = append(results, reporter.CheckResult{
						CheckName:      "sg-all-ports-open",
						Severity:       "CRITICAL",
						ResourceID:     sgID,
						Message:        fmt.Sprintf("Security group %q (%s) allows all traffic from 0.0.0.0/0", sgName, sgID),
						Recommendation: "Restrict security group rules to specific ports and CIDR ranges",
					})
					break
				}
				// SSH port 22 open
				if perm.FromPort != nil && perm.ToPort != nil &&
					*perm.FromPort <= 22 && *perm.ToPort >= 22 {
					results = append(results, reporter.CheckResult{
						CheckName:      "sg-ssh-open",
						Severity:       "CRITICAL",
						ResourceID:     sgID,
						Message:        fmt.Sprintf("Security group %q (%s) allows SSH (port 22) from 0.0.0.0/0", sgName, sgID),
						Recommendation: "Restrict SSH access to known IP ranges or use AWS Systems Manager Session Manager",
					})
				}
			}
		}
	}
	return results, nil
}

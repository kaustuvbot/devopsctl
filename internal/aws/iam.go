package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
)

// CheckIAMUsersMFA checks for IAM users without MFA enabled.
// Severity: HIGH
func CheckIAMUsersMFA(ctx context.Context, client IAMClient) ([]reporter.CheckResult, error) {
	var results []reporter.CheckResult

	users, err := paginateUsers(ctx, client)
	if err != nil {
		if isPermissionError(err) {
			return results, nil
		}
		return nil, fmt.Errorf("iam-mfa: %w", err)
	}

	for _, user := range users {
		mfaOut, err := client.ListMFADevices(ctx, &iam.ListMFADevicesInput{
			UserName: user.UserName,
		})
		if err != nil {
			if isPermissionError(err) {
				continue
			}
			return nil, fmt.Errorf("ListMFADevices for %s: %w", *user.UserName, err)
		}
		if len(mfaOut.MFADevices) == 0 {
			results = append(results, reporter.CheckResult{
				CheckName:      "iam-mfa-disabled",
				Severity:       "HIGH",
				ResourceID:     *user.UserName,
				Message:        fmt.Sprintf("IAM user %q has no MFA device enabled", *user.UserName),
				Recommendation: "Enable MFA for all IAM users",
			})
		}
	}
	return results, nil
}

// CheckIAMAccessKeyAge checks for active access keys older than the threshold.
// Severity: MEDIUM (>keyAgeDays days), HIGH (>120 days).
func CheckIAMAccessKeyAge(ctx context.Context, client IAMClient, keyAgeDays int) ([]reporter.CheckResult, error) {
	var results []reporter.CheckResult

	users, err := paginateUsers(ctx, client)
	if err != nil {
		if isPermissionError(err) {
			return results, nil
		}
		return nil, fmt.Errorf("iam-key-age: %w", err)
	}

	now := time.Now()
	for _, user := range users {
		keysOut, err := client.ListAccessKeys(ctx, &iam.ListAccessKeysInput{
			UserName: user.UserName,
		})
		if err != nil {
			if isPermissionError(err) {
				continue
			}
			return nil, fmt.Errorf("ListAccessKeys for %s: %w", *user.UserName, err)
		}
		for _, key := range keysOut.AccessKeyMetadata {
			if key.Status == iamtypes.StatusTypeInactive {
				continue
			}
			ageDays := int(now.Sub(*key.CreateDate).Hours() / 24)
			if ageDays > keyAgeDays {
				severity := "MEDIUM"
				if ageDays > 120 {
					severity = "HIGH"
				}
				results = append(results, reporter.CheckResult{
					CheckName:      "iam-old-access-key",
					Severity:       severity,
					ResourceID:     *key.AccessKeyId,
					Message:        fmt.Sprintf("Access key for %q is %d days old", *user.UserName, ageDays),
					Recommendation: "Rotate access keys regularly; delete unused keys",
				})
			}
		}
	}
	return results, nil
}

const adminPolicyARN = "arn:aws:iam::aws:policy/AdministratorAccess"

// CheckIAMAdminUsers detects IAM users with AdministratorAccess attached directly or via group.
// Severity: CRITICAL
func CheckIAMAdminUsers(ctx context.Context, client IAMClient) ([]reporter.CheckResult, error) {
	var results []reporter.CheckResult

	users, err := paginateUsers(ctx, client)
	if err != nil {
		if isPermissionError(err) {
			return results, nil
		}
		return nil, fmt.Errorf("iam-admin: %w", err)
	}

	for _, user := range users {
		isAdmin, err := userHasAdminAccess(ctx, client, *user.UserName)
		if err != nil {
			continue
		}
		if isAdmin {
			results = append(results, reporter.CheckResult{
				CheckName:      "iam-admin-access",
				Severity:       "CRITICAL",
				ResourceID:     *user.UserName,
				Message:        fmt.Sprintf("IAM user %q has AdministratorAccess policy", *user.UserName),
				Recommendation: "Apply least-privilege; remove AdministratorAccess from regular users",
			})
		}
	}
	return results, nil
}

// paginateUsers returns all IAM users, handling truncated responses.
func paginateUsers(ctx context.Context, client IAMClient) ([]iamtypes.User, error) {
	var users []iamtypes.User
	var marker *string
	for {
		out, err := client.ListUsers(ctx, &iam.ListUsersInput{Marker: marker})
		if err != nil {
			return nil, err
		}
		users = append(users, out.Users...)
		if !out.IsTruncated {
			break
		}
		marker = out.Marker
	}
	return users, nil
}

func userHasAdminAccess(ctx context.Context, client IAMClient, userName string) (bool, error) {
	policiesOut, err := client.ListAttachedUserPolicies(ctx, &iam.ListAttachedUserPoliciesInput{
		UserName: &userName,
	})
	if err == nil {
		for _, p := range policiesOut.AttachedPolicies {
			if *p.PolicyArn == adminPolicyARN {
				return true, nil
			}
		}
	}

	groupsOut, err := client.ListGroupsForUser(ctx, &iam.ListGroupsForUserInput{
		UserName: &userName,
	})
	if err != nil {
		return false, nil
	}
	for _, grp := range groupsOut.Groups {
		grpPolicies, err := client.ListAttachedGroupPolicies(ctx, &iam.ListAttachedGroupPoliciesInput{
			GroupName: grp.GroupName,
		})
		if err != nil {
			continue
		}
		for _, p := range grpPolicies.AttachedPolicies {
			if *p.PolicyArn == adminPolicyARN {
				return true, nil
			}
		}
	}
	return false, nil
}

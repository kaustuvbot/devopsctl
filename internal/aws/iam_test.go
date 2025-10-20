package aws

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type mockIAMClient struct {
	listUsersOutput                 *iam.ListUsersOutput
	listUsersErr                    error
	listMFADevicesOutput            *iam.ListMFADevicesOutput
	listMFADevicesErr               error
	listAccessKeysOutput            *iam.ListAccessKeysOutput
	listAccessKeysErr               error
	listAttachedUserPoliciesOutput  *iam.ListAttachedUserPoliciesOutput
	listAttachedUserPoliciesErr     error
	listGroupsForUserOutput         *iam.ListGroupsForUserOutput
	listGroupsForUserErr            error
	listAttachedGroupPoliciesOutput *iam.ListAttachedGroupPoliciesOutput
	listAttachedGroupPoliciesErr    error
}

func (m *mockIAMClient) ListUsers(_ context.Context, _ *iam.ListUsersInput, _ ...func(*iam.Options)) (*iam.ListUsersOutput, error) {
	return m.listUsersOutput, m.listUsersErr
}
func (m *mockIAMClient) ListMFADevices(_ context.Context, _ *iam.ListMFADevicesInput, _ ...func(*iam.Options)) (*iam.ListMFADevicesOutput, error) {
	return m.listMFADevicesOutput, m.listMFADevicesErr
}
func (m *mockIAMClient) ListAccessKeys(_ context.Context, _ *iam.ListAccessKeysInput, _ ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error) {
	return m.listAccessKeysOutput, m.listAccessKeysErr
}
func (m *mockIAMClient) ListAttachedUserPolicies(_ context.Context, _ *iam.ListAttachedUserPoliciesInput, _ ...func(*iam.Options)) (*iam.ListAttachedUserPoliciesOutput, error) {
	return m.listAttachedUserPoliciesOutput, m.listAttachedUserPoliciesErr
}
func (m *mockIAMClient) ListGroupsForUser(_ context.Context, _ *iam.ListGroupsForUserInput, _ ...func(*iam.Options)) (*iam.ListGroupsForUserOutput, error) {
	return m.listGroupsForUserOutput, m.listGroupsForUserErr
}
func (m *mockIAMClient) ListAttachedGroupPolicies(_ context.Context, _ *iam.ListAttachedGroupPoliciesInput, _ ...func(*iam.Options)) (*iam.ListAttachedGroupPoliciesOutput, error) {
	return m.listAttachedGroupPoliciesOutput, m.listAttachedGroupPoliciesErr
}

func TestCheckIAMUsersMFA_NoMFA(t *testing.T) {
	mock := &mockIAMClient{
		listUsersOutput:      &iam.ListUsersOutput{Users: []iamtypes.User{{UserName: aws.String("alice")}}},
		listMFADevicesOutput: &iam.ListMFADevicesOutput{MFADevices: []iamtypes.MFADevice{}},
	}
	results, err := CheckIAMUsersMFA(context.Background(), mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Severity != "HIGH" {
		t.Errorf("expected HIGH severity, got %s", results[0].Severity)
	}
	if results[0].ResourceID != "alice" {
		t.Errorf("expected resource 'alice', got %s", results[0].ResourceID)
	}
}

func TestCheckIAMUsersMFA_WithMFA(t *testing.T) {
	mock := &mockIAMClient{
		listUsersOutput:      &iam.ListUsersOutput{Users: []iamtypes.User{{UserName: aws.String("bob")}}},
		listMFADevicesOutput: &iam.ListMFADevicesOutput{MFADevices: []iamtypes.MFADevice{{SerialNumber: aws.String("arn:aws:iam::123:mfa/token")}}},
	}
	results, err := CheckIAMUsersMFA(context.Background(), mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for user with MFA, got %d", len(results))
	}
}

func TestCheckIAMUsersMFA_PermissionError(t *testing.T) {
	mock := &mockIAMClient{
		listUsersErr: fmt.Errorf("AccessDenied: not authorized"),
	}
	results, err := CheckIAMUsersMFA(context.Background(), mock)
	if err != nil {
		t.Errorf("expected graceful skip on permission error, got: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results on permission error, got %d", len(results))
	}
}

func TestCheckIAMAccessKeyAge_Medium(t *testing.T) {
	created := time.Now().AddDate(0, 0, -100)
	mock := &mockIAMClient{
		listUsersOutput: &iam.ListUsersOutput{Users: []iamtypes.User{{UserName: aws.String("alice")}}},
		listAccessKeysOutput: &iam.ListAccessKeysOutput{
			AccessKeyMetadata: []iamtypes.AccessKeyMetadata{{
				AccessKeyId: aws.String("AKIAIOSFODNN7EXAMPLE"),
				CreateDate:  &created,
				Status:      iamtypes.StatusTypeActive,
				UserName:    aws.String("alice"),
			}},
		},
	}
	results, err := CheckIAMAccessKeyAge(context.Background(), mock, 90)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Severity != "MEDIUM" {
		t.Errorf("expected MEDIUM for 100-day key, got %s", results[0].Severity)
	}
}

func TestCheckIAMAccessKeyAge_High(t *testing.T) {
	created := time.Now().AddDate(0, 0, -130)
	mock := &mockIAMClient{
		listUsersOutput: &iam.ListUsersOutput{Users: []iamtypes.User{{UserName: aws.String("alice")}}},
		listAccessKeysOutput: &iam.ListAccessKeysOutput{
			AccessKeyMetadata: []iamtypes.AccessKeyMetadata{{
				AccessKeyId: aws.String("AKIAIOSFODNN7EXAMPLE"),
				CreateDate:  &created,
				Status:      iamtypes.StatusTypeActive,
				UserName:    aws.String("alice"),
			}},
		},
	}
	results, _ := CheckIAMAccessKeyAge(context.Background(), mock, 90)
	if len(results) != 1 || results[0].Severity != "HIGH" {
		t.Errorf("expected HIGH for 130-day key, got %v", results)
	}
}

func TestCheckIAMAccessKeyAge_Fresh(t *testing.T) {
	created := time.Now().AddDate(0, 0, -30)
	mock := &mockIAMClient{
		listUsersOutput: &iam.ListUsersOutput{Users: []iamtypes.User{{UserName: aws.String("alice")}}},
		listAccessKeysOutput: &iam.ListAccessKeysOutput{
			AccessKeyMetadata: []iamtypes.AccessKeyMetadata{{
				AccessKeyId: aws.String("AKIAIOSFODNN7EXAMPLE"),
				CreateDate:  &created,
				Status:      iamtypes.StatusTypeActive,
				UserName:    aws.String("alice"),
			}},
		},
	}
	results, _ := CheckIAMAccessKeyAge(context.Background(), mock, 90)
	if len(results) != 0 {
		t.Errorf("expected no results for fresh key, got %d", len(results))
	}
}

func TestCheckIAMAdminUsers_DirectPolicy(t *testing.T) {
	mock := &mockIAMClient{
		listUsersOutput: &iam.ListUsersOutput{Users: []iamtypes.User{{UserName: aws.String("admin-user")}}},
		listAttachedUserPoliciesOutput: &iam.ListAttachedUserPoliciesOutput{
			AttachedPolicies: []iamtypes.AttachedPolicy{{
				PolicyArn:  aws.String(adminPolicyARN),
				PolicyName: aws.String("AdministratorAccess"),
			}},
		},
		listGroupsForUserOutput: &iam.ListGroupsForUserOutput{},
	}
	results, err := CheckIAMAdminUsers(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Severity != "CRITICAL" {
		t.Errorf("expected CRITICAL for admin user, got %v", results)
	}
}

func TestCheckIAMAdminUsers_NoAdmin(t *testing.T) {
	mock := &mockIAMClient{
		listUsersOutput: &iam.ListUsersOutput{Users: []iamtypes.User{{UserName: aws.String("regular-user")}}},
		listAttachedUserPoliciesOutput: &iam.ListAttachedUserPoliciesOutput{
			AttachedPolicies: []iamtypes.AttachedPolicy{{
				PolicyArn:  aws.String("arn:aws:iam::aws:policy/ReadOnlyAccess"),
				PolicyName: aws.String("ReadOnlyAccess"),
			}},
		},
		listGroupsForUserOutput: &iam.ListGroupsForUserOutput{},
	}
	results, _ := CheckIAMAdminUsers(context.Background(), mock)
	if len(results) != 0 {
		t.Errorf("expected no results for non-admin user, got %d", len(results))
	}
}

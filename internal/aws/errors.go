package aws

import "strings"

// isPermissionError returns true if an AWS API error is authorization-related.
// Callers skip checks gracefully on these errors rather than propagating them.
func isPermissionError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "AccessDenied") ||
		strings.Contains(msg, "UnauthorizedOperation") ||
		strings.Contains(msg, "AuthFailure") ||
		strings.Contains(msg, "NoCredentialProviders")
}

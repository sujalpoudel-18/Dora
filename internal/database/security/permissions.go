package security

type PermissionRule struct {
	ID          string
	Description string
}

func CheckPrivileges(role string) []PermissionRule {
	var rules []PermissionRule
	if role == "SUPERUSER" || role == "root" {
		rules = append(rules, PermissionRule{
			ID:          "DB-PERM-001",
			Description: "Account possesses superuser or administrative privileges",
		})
	}
	return rules
}

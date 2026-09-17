package security

type ConfigAudit struct {
	Setting string
	Status  string
}

func AuditBaseline() []ConfigAudit {
	return []ConfigAudit{
		{Setting: "remote_root_login", Status: "disabled"},
	}
}

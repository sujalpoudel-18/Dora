package security

type EncryptionStatus struct {
	TLSEnforced bool
	Cipher      string
}

func EvaluateEncryption(sslEnabled bool) EncryptionStatus {
	return EncryptionStatus{
		TLSEnforced: sslEnabled,
	}
}

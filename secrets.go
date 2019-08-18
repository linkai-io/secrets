package secrets

// Secrets interface for securely or locally accessing secrets
type Secrets interface {
	// GetSecureParameter retrieves the value from our secure parameter store
	GetSecureParameter(key string) ([]byte, error)
	// SetSecureParameter sets the value for the key in our secure parameter store
	SetSecureParameter(key, value string) error
}

package secrets

import (
	"fmt"
)

// SecretsCache for accessing cached/stored secrets
type SecretsCache struct {
	Region      string
	Environment string
	secrets     Secrets
}

// NewSecretsCache returns an instance for acquiring the secrets from either local env vars or AWS
func NewSecretsCache(env, region string) *SecretsCache {
	s := &SecretsCache{Environment: env, Region: region}
	if s.Environment != "local" {
		s.secrets = NewAWSSecrets(region)
	} else {
		s.secrets = NewEnvSecrets()
	}
	return s
}

// GetSecureString allows caller to provide the full key to return a string value
func (s *SecretsCache) GetSecureString(key string) (string, error) {
	data, err := s.secrets.GetSecureParameter(key)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Password retrieves the password from our secrets implementation.
// Additional methods can be added to SecretsCache for ensuring uniform
// access.
func (s *SecretsCache) Password() (string, error) {
	data, err := s.secrets.GetSecureParameter(fmt.Sprintf("/some/%s/password", s.Environment))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

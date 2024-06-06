package auth

// SecretKeeper the interface to get the secret
type SecretKeeper interface {
	GetSecret(key string) (string, error)
}

type DefaultProvider struct {
	AppKey    string
	AppSecret string
}

func (dp DefaultProvider) GetSecret(key string) (string, error) {
	return dp.AppSecret, nil
}

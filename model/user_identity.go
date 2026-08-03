package model

type UserIdentity struct {
	BaseModel
	UserID          uint64 `json:"user_id" column:"user_id"`
	Provider        string `json:"provider" column:"provider"`
	ProviderSubject string `json:"provider_subject" column:"provider_subject"`
	Email           string `json:"email" column:"email"`
}

func (identity *UserIdentity) TableName() string {
	return "user_identities"
}

func (identity *UserIdentity) Create() error {
	return DB.Self.Create(identity).Error
}

func GetUserIdentity(provider, providerSubject string) (*UserIdentity, error) {
	identity := &UserIdentity{}
	d := DB.Self.Where("provider = ? AND provider_subject = ?", provider, providerSubject).First(identity)
	return identity, d.Error
}

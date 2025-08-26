package service

type Validator struct {
}

func (v *Validator) ValidateSupportAsset() bool {
	return true
}

func (v *Validator) ValidateMarketPrice() bool {
	return true
}

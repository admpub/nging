package lib

type DefaultValidator struct {
	NowVersions []string
}

func (v *DefaultValidator) Validate(data *LicenseData) error {
	if err := data.CheckExpiration(); err != nil {
		return err
	}
	if err := data.CheckVersion(v.NowVersions...); err != nil {
		return err
	}
	if err := data.CheckMAC(); err != nil {
		return err
	}
	return nil
}

package read

type ReadConfig struct {
	UsePreview  bool
	SpaceID     string
	Environment string
	AccessToken string
	Locale      string
	// e.g. 200 items per page
}

package processor

// Config holds the configuration for the file processor
type Config struct {
	LicenseText string
	Input       string
	Skip        string
	Prompt      bool
	DryRun      bool
	Verbose     bool
	IgnoreFail  bool   // If true will return 0 on a fail
	PresetStyle string // Header/Footer style
	PreferMulti bool
}

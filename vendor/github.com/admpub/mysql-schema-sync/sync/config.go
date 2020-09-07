package sync

type Config struct {
	Sync            bool
	Drop            bool
	SourceDSN       string
	DestDSN         string
	AlterIgnore     string
	Tables          string
	SkipTables      string
	MailTo          string
	SQLPreprocessor func(string) string `json:"-" xml:"-"`
}

type EmailConfig struct {
	On       bool
	SMTPHost string
	From     string
	Password string
	To       string
}

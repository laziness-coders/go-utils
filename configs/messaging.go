package configs

// TelegramConfig represents Telegram bot configuration.
type TelegramConfig struct {
	BotToken        string `mapstructure:"BOT_TOKEN"`
	ChannelID       int64  `mapstructure:"CHANNEL_ID"`
	MessageThreadID int64  `mapstructure:"MESSAGE_THREAD_ID"`
	Enabled         bool   `mapstructure:"ENABLED"`
}

// EmailConfig represents email/SMTP configuration.
type EmailConfig struct {
	SMTPHost string `mapstructure:"SMTP_HOST"`
	SMTPPort int    `mapstructure:"SMTP_PORT"`
	From     string `mapstructure:"FROM"`
	Password string `mapstructure:"PASSWORD"`
	UseTLS   bool   `mapstructure:"USE_TLS"`
	Enabled  bool   `mapstructure:"ENABLED"`
}

package dto

type TwilioWebhook struct {
	Id                string
	Status            string
	MessageStatus     string
	ChannelToAddress  string
	To                string
	ChannelPrefix     string
	MessageSid        string
	AccountSid        string
	From              string
	ApiVersion        string
	ChannelInstallSid string
}

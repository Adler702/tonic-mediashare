package structs

type Config struct {
	Port                int    `json:"port"`
	BaseURL             string `json:"baseurl"`
	MongoURL            string `json:"mongourl"`
	DiscordClientId     string `json:"discordclientid"`
	DiscordClientSecret string `json:"discordclientsecret"`
	State               string `json:"discordsecretstate"`
}

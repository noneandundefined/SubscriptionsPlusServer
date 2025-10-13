package machinelearning

import "strings"

type SubscriptionImage struct {
	Keywords []string
	Image    string
}

var SubscriptionImages = []SubscriptionImage{
	// Яндекс Плюс
	{Keywords: []string{
		"yandex plus", "yandex+", "яндекс плюс", "яндекс+", "yandex music", "яндекс музыка",
		"yandex подписка", "яндекс подписка", "yplus", "yandex", "yandex premium",
		"яндекс премиум", "yandex pl", "яндекс pl", "яндексплюс", "yandexplus",
	}, Image: "yplus.png"},

	// Сбер Прайм
	{Keywords: []string{
		"sber", "sber prime", "sberprime", "сбер прайм", "сберпрайм",
		"sber подписка", "сбер подписка", "сбербанк подписка", "сбербанк",
	}, Image: "sber.jpg"},

	// Subscription Plus
	{Keywords: []string{
		"subscription plus", "sub plus", "subscription+", "subplus",
		"подписка плюс", "subscriptionplus", "sub plus app",
	}, Image: "subscriptionplus.png"},

	// Tinkoff / TBank
	{Keywords: []string{
		"tbank", "t-pro", "t pro", "t-bank", "tinkoff", "тинькоф", "тинькофф",
		"tinkoff pro", "tinkoff black", "tbank подписка", "tinkoff подписка",
	}, Image: "tbank.jpg"},

	// Telegram
	{Keywords: []string{
		"telegram", "telegram premium", "tg premium", "tg+", "tg",
		"телеграм премиум", "тг премиум", "tg подписка", "telegram подписка",
	}, Image: "tg.png"},

	// Spotify
	{Keywords: []string{
		"spotify", "spoti", "spoti fy", "спотифай", "споттифай",
		"spotify premium", "spotify подписка", "споти премиум",
	}, Image: "spotify.jpg"},

	// Netflix
	{Keywords: []string{
		"netflix", "нетфликс", "нетфликс подписка", "netflix premium", "netflix ultra",
		"net flix", "net flix подписка",
	}, Image: "netflix.jpg"},

	// YouTube Premium
	{Keywords: []string{
		"youtube", "youtube premium", "ютуб премиум", "ютуб+", "yt premium",
		"youtube+", "ютуб подписка", "youtube music",
	}, Image: "youtube.png"},

	// Apple One / Music / TV+
	{Keywords: []string{
		"apple", "apple one", "apple music", "apple tv+", "apple tv plus",
		"apple подписка", "эпл музыка", "эпл подписка",
	}, Image: "apple.png"},

	// Google One
	{Keywords: []string{
		"google one", "гугл one", "гугл подписка", "google диск", "google storage", "google облако",
	}, Image: "googleone.jpg"},

	// Microsoft 365 / Office 365
	{Keywords: []string{
		"microsoft 365", "office 365", "microsoft подписка", "офис 365", "microsoft one drive",
	}, Image: "office365.jpg"},

	// Adobe
	{Keywords: []string{
		"adobe", "adobe cc", "creative cloud", "фотошоп подписка", "adobe подписка",
	}, Image: "adobe.jpg"},

	// Discord Nitro
	{Keywords: []string{
		"discord", "discord nitro", "nitro", "дискорд нитро", "дискорд подписка",
	}, Image: "discord.png"},

	// Steam / Xbox / PS Plus
	{Keywords: []string{"steam", "steam subscription", "steam подписка", "steam premium"}, Image: "steam.png"},
	{Keywords: []string{"xbox", "xbox game pass", "gamepass", "геймпас", "xbox подписка"}, Image: "xbox.jpg"},
	{Keywords: []string{"playstation", "ps plus", "пс плюс", "плейстейшн плюс", "ps подписка"}, Image: "psplus.jpg"},

	// IVI / Okko / Kinopoisk
	{Keywords: []string{"ivi", "иви", "ivi подписка"}, Image: "ivi.jpg"},
	{Keywords: []string{"okko", "окко", "okko подписка", "окко подписка"}, Image: "okko.png"},
	{Keywords: []string{"kinopoisk", "кино поиск", "кинопоиск", "кинопоиск подписка"}, Image: "kinopoisk.jpg"},

	// VPN
	{Keywords: []string{
		"vpn", "v p n", "впн", "прокси", "proxy", "vpn service", "vpn подписка",
		"vpn premium", "windscribe", "nordvpn", "surfshark", "expressvpn", "vpn plus",
	}, Image: "vpn.jpg"},

	// Cloud / VDS / VPS
	{Keywords: []string{
		"vds", "виртуальный сервер", "дедик", "dedicated", "dedicated server",
		"vdsina", "timeweb", "reg.ru vps", "reg vps",
	}, Image: "vds.jpg"},
	{Keywords: []string{"vps", "сервер", "reg.ru vps", "reg vps"}, Image: "vps.jpg"},

	// AI / ChatGPT / Claude / Midjourney
	{Keywords: []string{"chatgpt", "gpt plus", "chat gpt+", "openai plus", "чатгпт подписка", "openai"}, Image: "chatgpt.png"},
	{Keywords: []string{"midjourney", "mid journey", "midjourney подписка"}, Image: "midjourney.jpg"},
	{Keywords: []string{"claude", "anthropic", "claude ai", "claude подписка"}, Image: "claude.jpg"},
}

func (nlp *NLPBuilder) GetSubscriptionImage(input string) string {
	if input == "" {
		return ""
	}

	processed := strings.ToLower(input)
	// tokens := nlp.Preprocess(processed)

	for _, sub := range SubscriptionImages {
		for _, kw := range sub.Keywords {
			if strings.Contains(processed, strings.ToLower(kw)) {
				return sub.Image
			}
		}
	}

	return ""
}

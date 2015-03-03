package def

type Config struct {
	Producers *Producers `toml:"producers"`
}

type Producers struct {
	Github *GithubConfig `toml:"github"`
}

type GithubConfig struct {
	AccessToken string `toml:"token"`
}

type Event struct {
	Type string

	Summary string
	Message string
}

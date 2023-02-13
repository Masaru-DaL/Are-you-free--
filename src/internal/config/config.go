package config

import "github.com/jinzhu/configor"

var Config = struct {
	DB struct {
		SQLDriver string `yaml:"driver"`
		DBPath    string `yaml:"path"`
		FilePath  string `yaml:"file"`
	}

	FilePath struct {
		CSS   string `yaml:"CSS"`
		JS    string `yaml:"JS"`
		Image string `yaml:"image"`
		Icon  string `yaml:"icon"`
	}

	AUTH struct {
		Session              string `yaml:"session"`
		SessionKey           string `yaml:"sessionKey"`
		SessionPath          string `yaml:"sessionPath"`
		SessionExpirationSec int    `yaml:"sessionExpirationSec"`
		SessionExpirationDay int    `yaml:"sessionExpirationDay"`
	}
}{}

/*
	 config.yamlを読み込む関数,
		読み込んだ後は構造体にアクセスして環境変数を読み込む
*/
func LoadConfigForYaml() {
	err := configor.New(&configor.Config{
		AutoReload: true,
	}).Load(&Config, "internal/config/config.yaml")

	if err != nil {
		panic("ERROR: config cannot be loaded")
	}
}

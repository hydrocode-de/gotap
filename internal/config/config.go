package config

import "github.com/spf13/viper"

var v *viper.Viper

func Init() {
	v = viper.New()

	setupDefaults()

	v.SetEnvPrefix("TAP")
	v.AutomaticEnv()
}

func GetViper() *viper.Viper {
	return v
}

func setupDefaults() {
	v.SetDefault("spec_file", "/src/tool.yml")
	v.SetDefault("input_file", "/in/inputs.json")
	v.SetDefault("citation_file", "/src/CITATION.cff")
	v.SetDefault("license_file", "/src/LICENSE")
}

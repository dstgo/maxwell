package cfgx

import "github.com/spf13/viper"

// LoadConfigFile load config from file
func LoadConfigFile(file string) (map[string]any, error) {
	v := viper.New()
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v.AllSettings(), nil
}

// LoadConfigAndMapTo load config from file and map to specific struct
func LoadConfigAndMapTo(file string, mapTo any) error {
	v := viper.New()
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	return v.Unmarshal(mapTo)
}

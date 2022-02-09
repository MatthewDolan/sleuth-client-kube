package config

import "github.com/spf13/viper"

type App struct {
	Sleuth Sleuth

	Kubernetes Kubernetes
}

type Sleuth struct {
	OrganizationSlug string

	APIKey string
}

type Kubernetes struct {
	Environment string

	Annotations KubernetesAnnotations
}

type KubernetesAnnotations struct {
	DeployedAtKey string

	DeploymentSlugKey string

	SHAKey string
}

func Load(configPathOpt []string) (*App, error) {
	v := viper.New()
	if len(configPathOpt) > 0 {
		for _, configPath := range configPathOpt {
			viper.AddConfigPath(configPath)
		}
	} else {
		viper.AddConfigPath(".")
	}
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var client App
	if err := v.Unmarshal(&client); err != nil {
		return nil, err
	}

	return &client, nil
}

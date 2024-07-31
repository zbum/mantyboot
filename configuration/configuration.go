package configuration

import "gopkg.in/yaml.v3"

type Configuration[T any] struct {
	payload *T
}

func NewConfiguration[T any]() *Configuration[T] {
	return &Configuration[T]{}
}

func (c *Configuration[T]) Load(input string) error {

	buf := []byte(input)
	var config T
	err := yaml.Unmarshal(buf, &config)
	if err != nil {
		return err
	}

	c.payload = &config
	return nil
}

func (c *Configuration[T]) GetConfiguration() *T {
	return c.payload
}

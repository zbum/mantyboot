package configuration

import (
	"embed"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"os"
	"strings"
)

type Configuration[T any] struct {
	embedDir embed.FS
	profile  string
	payload  *T
}

func NewConfiguration[T any](embedDir embed.FS, profile string) *Configuration[T] {
	return &Configuration[T]{
		embedDir: embedDir,
		profile:  profile,
	}
}

func (c *Configuration[T]) GetConfiguration() *T {
	return c.payload
}

func (c *Configuration[T]) Load() (*T, error) {
	readers := c.findConfigurationFiles(c.profile)
	for _, reader := range readers {
		bytes, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		c.parse(bytes)
	}

	return c.payload, nil
}

func (c *Configuration[T]) parse(input []byte) error {

	if c.payload == nil {
		var config T
		c.payload = &config
	}
	err := yaml.Unmarshal(input, c.payload)
	if err != nil {
		return err
	}

	return nil
}

func (c *Configuration[T]) findConfigurationFiles(profile string) []io.Reader {
	var files []io.Reader

	firstCandidate := c.findFirstCandidate(profile)
	if firstCandidate != nil {
		files = append(files, firstCandidate)
	}

	secondCandidate := c.findSecondCandidate(profile)
	if secondCandidate != nil {
		files = append(files, secondCandidate)
	}

	thirdCandidate := c.findThirdCandidate(profile)
	if thirdCandidate != nil {
		files = append(files, thirdCandidate)
	}

	return files
}

func (c *Configuration[T]) findFirstCandidate(profile string) io.Reader {
	var property io.Reader
	fs.WalkDir(c.embedDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, "application-"+profile+".yaml") {
			property, err = c.embedDir.Open(path)
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	})
	return property
}

func (c *Configuration[T]) findSecondCandidate(profile string) *os.File {
	wd, err := os.Getwd()
	if err != nil {
		return nil
	}

	property, err := os.Open(wd + "/application-" + profile + ".yaml")

	if err != nil {
		return nil
	}

	return property
}

func (c *Configuration[T]) findThirdCandidate(profile string) *os.File {
	wd, err := os.Getwd()
	if err != nil {
		return nil
	}

	property, err := os.Open(wd + "/config/application-" + profile + ".yaml")

	if err != nil {
		return nil
	}

	return property
}

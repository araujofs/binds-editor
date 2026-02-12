package files

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
)

type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (f File) FilterValue() string {
	return f.Name
}

func (f File) Description() string {
	return f.Path
}

func (f File) Title() string {
	return f.Name
}

type Configuration struct {
	Files []*File `json:"files"`
}

// searchSlice tries to find a *File inside a []*File slice
// based on the fieldName and fieldValue.
// Returns the index if it exists or -1 if it doesn't.
func searchSlice(s []*File, fieldName string, fieldValue string) int {
	for idx, file := range s {
		if (fieldName == "Name" && file.Name == fieldValue) || (fieldName == "Path" && file.Path == fieldValue) {
			return idx
		}
	}

	return -1
}

func (c *Configuration) fileExists(path *string, name *string) (int, int) {
	if name != nil {
		if nameExists := searchSlice(c.Files, "Name", *name); nameExists != -1 {
			return nameExists, 1
		}
	}

	if path != nil {
		if pathExists := searchSlice(c.Files, "Path", *path); pathExists != -1 {
			return pathExists, 2
		}
	}

	return -1, -1
}

func (c *Configuration) AddFile(path string, name string) error {
	_, exists := c.fileExists(&path, &name)

	if exists == 1 {
		return fmt.Errorf(`bindings file already exists with name "%s"`, name)
	}

	if exists == 2 {
		return fmt.Errorf(`bindings file already exists with path "%s"`, path)
	}

	c.Files = append(c.Files, &File{Name: name, Path: path})

	return nil
}

func (c *Configuration) RemoveFile(name string) error {
	idx, exists := c.fileExists(nil, &name)

	if exists == -1 {
		return fmt.Errorf(`bindings file with name "%s" doesn't exist`, name)
	}

	c.Files = slices.Delete(c.Files, idx, idx+1)

	return nil
}

func (c *Configuration) EditFile(oldName string, newName string, newPath string) error {
	idx, exists := c.fileExists(nil, &oldName)

	if exists == -1 {
		return fmt.Errorf(`bindings file with name "%s" doesn't exist`, oldName)
	}

	c.Files[idx] = &File{Name: newName, Path: newPath}

	return nil
}

func (c *Configuration) SaveConfiguration() error {
	return writeConfigurationFile(c)
}

func GetConfigData() *Configuration {
	var config Configuration

	file, err := os.Open(getConfigurationPath())
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		return &Configuration{
			Files: []*File{},
		}
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil && err != io.EOF {
		panic(err)
	}

	return &config
}

// writeConfiguration writes the Configuration struct to the
// configuration file as JSON or creates it, if it doesn't exist,
// and then writes.
func writeConfigurationFile(config *Configuration) error {
	if config == nil {
		return errors.New("configuration is nil")
	}

	file, err := cleanConfigurationFile()
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	defer file.Close()

	return encoder.Encode(*config)
}

func cleanConfigurationFile() (*os.File, error) {
	dir := filepath.Dir(getConfigurationPath())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.Create(getConfigurationPath())
	if err != nil {
		return nil, err
	}

	return file, nil
}

func getConfigurationPath() string {
	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".config", "binds-editor", "config.json")
}

package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// TODO: reconsider the necessity of parametrizing interface methods
// when default values can be passed straight to writeYaml and readYaml

type YamlData interface {
}

func writeYaml(filename string, yd YamlData) error {
	yamlData, err := yaml.Marshal(yd)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, yamlData, 0644)
}

func readYaml(filename string, yd YamlData) error {
	yamlData, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s not found and will be created if services are down\n", filename)
			NewServiceStatus()
			return nil
		}
		return err
	}

	err = yaml.UnmarshalStrict(yamlData, yd)
	if err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			return fmt.Errorf("error reading %s: file contains incompatible value type:\n %v", filename, err)
		} else if strings.Contains(err.Error(), "not found in type") {
			return fmt.Errorf("error reading %s: file contains unexpected key:\n %v", filename, err)
		} else {
			return err
		}
	}
	return nil
}

// func validateFields(filename string, yd YamlData) error {
// 	structType := reflect.TypeOf(yd).Elem()
// 	structFieldCount := structType.NumField()

// 	data := make(map[string]any, structFieldCount)

// 	if err := readYaml(filename, &data); err != nil {
// 		return err
// 	}

// 	// var missingKeys []string

// 	for i := 0; i < structFieldCount; i++ {
// 		field := structType.Field(i)
// 		yamlTag := field.Tag.Get("yaml")

// 		_, key := data[yamlTag]

// 		if !key {
// 			// missingKeys = append(missingKeys, yamlTag)
// 			return fmt.Errorf("missing field in '%s': '%s'", filename, yamlTag)
// 		}
// 	}

// 	// if len(missingKeys) > 0 {
// 	// 	return fmt.Errorf("missing fields in %s:\n%s", filename, strings.Join(missingKeys, "\n"))
// 	// }

// 	if err := readYaml(filename, yd); err != nil {
// 		return err
// 	}

// 	return nil
// }

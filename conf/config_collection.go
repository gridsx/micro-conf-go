package conf

import (
	"encoding/json"
	"strconv"
	"strings"
)

const defaultSize = 8

func (c *ConfigContext) GetIntList(k string) []int64 {
	v, err := c.get(k)
	if err != nil {
		return nil
	}
	if len(v) == 0 {
		return nil
	}
	arr := strings.Split(v, ",")
	result := make([]int64, 0, defaultSize)
	for _, str := range arr {
		x, pe := strconv.ParseInt(str, 10, 64)
		if pe != nil {
			continue
		}
		result = append(result, x)
	}
	return result
}

func (c *ConfigContext) GetFloatList(k string) []float64 {
	v, err := c.get(k)
	if err != nil {
		return nil
	}
	if len(v) == 0 {
		return nil
	}
	arr := strings.Split(v, ",")
	result := make([]float64, 0, defaultSize)
	for _, str := range arr {
		x, pe := strconv.ParseFloat(str, 64)
		if pe != nil {
			continue
		}
		result = append(result, x)
	}
	return result
}

func (c *ConfigContext) GetStringList(k string) []string {
	v, err := c.get(k)
	if err != nil {
		return nil
	}
	if len(v) == 0 {
		return nil
	}
	return strings.Split(v, ",")
}

func (c *ConfigContext) GetStringAnyMap(k string) map[string]any {
	v, err := c.get(k)
	if err != nil {
		return nil
	}
	if len(v) == 0 {
		return nil
	}
	result := make(map[string]any, defaultSize)
	if err := json.Unmarshal([]byte(v), &result); err != nil {
		return nil
	}
	return result
}

func (c *ConfigContext) GetStringMap(k string) map[string]string {
	v, err := c.get(k)
	if err != nil {
		return nil
	}
	if len(v) == 0 {
		return nil
	}
	result := make(map[string]string, defaultSize)
	if err := json.Unmarshal([]byte(v), &result); err != nil {
		return nil
	}
	return result
}

func (c *ConfigContext) GetStringIntMap(k string) map[string]int64 {
	v, err := c.get(k)
	if err != nil {
		return nil
	}
	if len(v) == 0 {
		return nil
	}
	result := make(map[string]int64, defaultSize)
	if err := json.Unmarshal([]byte(v), &result); err != nil {
		return nil
	}
	return result
}
func (c *ConfigContext) GetStringBoolMap(k string) map[string]bool {
	v, err := c.get(k)
	if err != nil {
		return nil
	}
	if len(v) == 0 {
		return nil
	}
	result := make(map[string]bool, defaultSize)
	if err := json.Unmarshal([]byte(v), &result); err != nil {
		return nil
	}
	return result
}
func (c *ConfigContext) GetStringFloatMap(k string) map[string]float64 {
	v, err := c.get(k)
	if err != nil {
		return nil
	}
	if len(v) == 0 {
		return nil
	}
	result := make(map[string]float64, defaultSize)
	if err := json.Unmarshal([]byte(v), &result); err != nil {
		return nil
	}
	return result
}

func (c *ConfigContext) GetObject(k string, o any) error {
	v, err := c.get(k)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(v), &o); err != nil {
		return err
	}
	return nil
}

package conf

import "strconv"

func (c *ConfigContext) GetBoolDefault(k string, x bool) bool {
	v, err := c.get(k)
	if err != nil {
		return x
	}
	b, parseErr := strconv.ParseBool(v)
	if parseErr != nil {
		return x
	}
	return b
}

func (c *ConfigContext) GetBool(k string) (bool, error) {
	v, err := c.get(k)
	if err != nil {
		return false, err
	}
	b, parseErr := strconv.ParseBool(v)
	if parseErr != nil {
		return false, err
	}
	return b, nil
}

func (c *ConfigContext) GetStringDefault(k, s string) string {
	v, err := c.get(k)
	if err != nil {
		return s
	}
	return v
}

func (c *ConfigContext) GetString(k string) (string, error) {
	v, err := c.get(k)
	if err != nil {
		return "", err
	}
	return v, nil
}

func (c *ConfigContext) GetIntDefault(k string, n int64) int64 {
	v, err := c.get(k)
	if err != nil {
		return n
	}
	r, parseErr := strconv.ParseInt(v, 10, 64)
	if parseErr != nil {
		return n
	}
	return r
}

func (c *ConfigContext) GetInt(k string) (int64, error) {
	v, err := c.get(k)
	if err != nil {
		return 0, err
	}
	r, parseErr := strconv.ParseInt(v, 10, 64)
	if parseErr != nil {
		return 0, parseErr
	}
	return r, nil
}

func (c *ConfigContext) GetFloatDefault(k string, f float64) float64 {
	v, err := c.get(k)
	if err != nil {
		return f
	}
	r, parseErr := strconv.ParseFloat(v, 64)
	if parseErr != nil {
		return f
	}
	return r
}

func (c *ConfigContext) GetFloat(k string) (float64, error) {
	v, err := c.get(k)
	if err != nil {
		return 0, err
	}
	r, parseErr := strconv.ParseFloat(v, 64)
	if parseErr != nil {
		return 0, parseErr
	}
	return r, nil
}

package mc

func (m *MicroClient) GetBoolDefault(k string, x bool) bool {
	return m.context.GetBoolDefault(k, x)
}

func (m *MicroClient) GetBool(k string) (bool, error) {
	return m.context.GetBool(k)
}

func (m *MicroClient) GetStringDefault(k, s string) string {
	return m.context.GetStringDefault(k, s)
}

func (m *MicroClient) GetString(k string) (string, error) {
	return m.context.GetString(k)
}

func (m *MicroClient) GetIntDefault(k string, n int64) int64 {
	return m.context.GetIntDefault(k, n)
}

func (m *MicroClient) GetInt(k string) (int64, error) {
	return m.context.GetInt(k)
}

func (m *MicroClient) GetFloatDefault(k string, f float64) float64 {
	return m.context.GetFloatDefault(k, f)
}

func (m *MicroClient) GetFloat(k string) (float64, error) {
	return m.context.GetFloat(k)
}

func (m *MicroClient) GetIntList(k string) []int64 {
	return m.context.GetIntList(k)
}

func (m *MicroClient) GetFloatList(k string) []float64 {
	return m.context.GetFloatList(k)
}

func (m *MicroClient) GetStringList(k string) []string {
	return m.context.GetStringList(k)
}

func (m *MicroClient) GetAnyMap(k string) map[string]any {
	return m.context.GetStringAnyMap(k)
}

func (m *MicroClient) GetMap(k string) map[string]string {
	return m.context.GetStringMap(k)
}

func (m *MicroClient) GetIntMap(k string) map[string]int64 {
	return m.context.GetStringIntMap(k)
}
func (m *MicroClient) GetBoolMap(k string) map[string]bool {
	return m.context.GetStringBoolMap(k)
}
func (m *MicroClient) GetFloatMap(k string) map[string]float64 {
	return m.context.GetStringFloatMap(k)
}

func (m *MicroClient) GetObject(k string, o any) error {
	return m.context.GetObject(k, o)
}

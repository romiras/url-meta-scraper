package options

type Option interface {
	Key() string
	Value() interface{}
}

type keyVal struct {
	key   string
	value interface{}
}

func (kv *keyVal) Key() string {
	return kv.key
}
func (kv *keyVal) Value() interface{} {
	return kv.value
}

func MetricName(metric string) Option {
	return &keyVal{
		key:   "metric_name",
		value: metric,
	}
}
func OptionsToMap(options []Option) map[string]interface{} {
	res := make(map[string]interface{})
	for _, o := range options {
		res[o.Key()] = o.Value()
	}
	return res
}

func GetOptionOrDefault(options map[string]interface{}, key string, defaultVal interface{}) interface{} {
	if val, ok := options[key]; ok && val != nil {
		return val
	}
	return defaultVal
}

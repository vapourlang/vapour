package cache

var cache = make(map[string]interface{})

func Set(key string, value interface{}) {
	cache[key] = value
}

func Get(key string) (interface{}, bool) {
	v, ok := cache[key]
	return v, ok
}

func Has(key string) bool {
	_, ok := cache[key]
	return ok
}

func Clear() {
	cache = make(map[string]interface{})
}

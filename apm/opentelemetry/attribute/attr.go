package attribute

import (
	"encoding/json"

	"go.opentelemetry.io/otel/attribute"
)

func parseAttrToKV(attr interface{}) []attribute.KeyValue {
	kv := make([]attribute.KeyValue, 0)
	b, err := json.Marshal(attr)
	if err != nil {
		return kv
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		return kv
	}

	for k, v := range m {
		// TODO: add more type if needed
		switch v.(type) {
		case string:
			kv = append(kv, attribute.String(k, v.(string)))

		case map[string]interface{}:
			bi, err := json.Marshal(v)
			if err != nil {
				continue
			}

			kv = append(kv, attribute.String(k, string(bi)))
		}
	}

	return kv
}

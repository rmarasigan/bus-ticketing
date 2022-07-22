package kvp

// Attribute contains a Key-Value Pair. It's a data representation
// in computing systems and applications.
type Attribute struct {
	Key   string
	Value interface{}
}

// KeyValue returns the key and value of KVP.
func (kv Attribute) KeyValue() (string, interface{}) {
	return kv.Key, kv.Value
}

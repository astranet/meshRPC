package db

import (
	"fmt"
	"reflect"
	"time"

	"github.com/globalsign/mgo/bson"
)

// GeneratePartialSet returns a $set bson map for a given model. This should be used with service models
// that use pointer for the variables and with omitempty tags.
func GeneratePartialSet(model interface{}) (bson.M, error) {
	updates := bson.M{}
	set := map[string]interface{}{}

	bytes, err := bson.Marshal(model)
	if err != nil {
		err = fmt.Errorf("GeneratePartialSet 1: %v", err)
		return updates, err
	}

	err = bson.Unmarshal(bytes, set)
	if err != nil {
		err = fmt.Errorf("GeneratePartialSet 2: %v", err)
		return updates, err
	}

	updates["$set"] = Flatten(set, false)

	return updates, err
}

// SetUpdatedOn sets the "updon" property on a $set bson map.
func SetUpdatedOn(updates *bson.M) {
	if (*updates)["$set"] == nil {
		(*updates)["$set"] = bson.M{"updon": time.Now()}
	} else if _, ok := (*updates)["$set"].(bson.M); ok {
		(*updates)["$set"].(bson.M)["updon"] = time.Now()
	} else if _, ok := (*updates)["$set"].(map[string]interface{}); ok {
		(*updates)["$set"].(map[string]interface{})["updon"] = time.Now()
	}
}

func Flatten(object map[string]interface{}, doFlattenSlice bool) map[string]interface{} {
	result := make(map[string]interface{})

	for k, raw := range object {
		flatten(result, k, reflect.ValueOf(raw), doFlattenSlice)
	}

	return result
}

func flatten(result map[string]interface{}, prefix string, v reflect.Value, doFlattenSlice bool) {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		flattenMap(result, prefix, v, doFlattenSlice)
	case reflect.Slice:
		if doFlattenSlice {
			flattenSlice(result, prefix, v, doFlattenSlice)
		} else {
			result[prefix] = v.Interface()
		}
	default:
		result[prefix] = v.Interface()
	}
}

func flattenMap(result map[string]interface{}, prefix string, v reflect.Value, doFlattenSlice bool) {
	for _, k := range v.MapKeys() {
		if k.Kind() == reflect.Interface {
			k = k.Elem()
		}

		if k.Kind() != reflect.String {
			panic(fmt.Sprintf("%s: map key is not string: %s", prefix, k))
		}

		flatten(result, fmt.Sprintf("%s.%s", prefix, k.String()), v.MapIndex(k), doFlattenSlice)
	}
}

func flattenSlice(result map[string]interface{}, prefix string, v reflect.Value, doFlattenSlice bool) {
	prefix = prefix + "."
	for i := 0; i < v.Len(); i++ {
		flatten(result, fmt.Sprintf("%s%d", prefix, i), v.Index(i), doFlattenSlice)
	}
}

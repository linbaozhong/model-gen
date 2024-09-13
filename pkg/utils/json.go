package utils

//
//func StringMarshal(v string) string {
//	return "\"" + strings.ReplaceAll(v, "\"", "\\\"") + "\""
//}
//func IntMarshal(v int64) string {
//	return strconv.FormatInt(v, 10)
//}
//func UintMarshal(v uint64) string {
//	return strconv.FormatUint(v, 10)
//}
//func FloatMarshal(v float64) string {
//	return strconv.FormatFloat(v, 'f', -1, 64)
//}
//func TimeMarshal(v time.Time) string {
//	return "\"" + v.Format(time.RFC3339Nano) + "\""
//}
//
////JSONValue redis存储序列化数据转换方法
//func JSONValue(i interface{}) string {
//	switch v := i.(type) {
//	case string:
//		return "\"" + strings.ReplaceAll(v, "\"", "\\\"") + "\""
//	case types.BigUint:
//		return v.String()
//	case types.Money:
//		return strconv.FormatFloat((float64(v) / 100), 'f', -1, 64)
//	case int:
//		return strconv.FormatInt(int64(v), 10)
//	case int8:
//		return strconv.FormatInt(int64(v), 10)
//	case int16:
//		return strconv.FormatInt(int64(v), 10)
//	case int32:
//		return strconv.FormatInt(int64(v), 10)
//	case int64:
//		return strconv.FormatInt(v, 10)
//	case uint:
//		return strconv.FormatUint(uint64(v), 10)
//	case uint8:
//		return strconv.FormatUint(uint64(v), 10)
//	case uint16:
//		return strconv.FormatUint(uint64(v), 10)
//	case uint32:
//		return strconv.FormatUint(uint64(v), 10)
//	case uint64:
//		return strconv.FormatUint(v, 10)
//	case float32:
//		return strconv.FormatFloat(float64(v), 'f', -1, 64)
//	case float64:
//		return strconv.FormatFloat(v, 'f', -1, 64)
//	case time.Time:
//		return "\"" + v.Format(time.RFC3339Nano) + "\""
//	case bool:
//		return strconv.FormatBool(v)
//	case []byte:
//		return "\"" + base64.StdEncoding.EncodeToString(v) + "\""
//	default:
//		rv := reflect.Indirect(reflect.ValueOf(v))
//		switch rv.Kind() {
//		case reflect.Struct:
//			if iter, ok := rv.Interface().(json.Marshaler); ok {
//				b, e := iter.MarshalJSON()
//				if e != nil {
//					panic(e)
//				}
//				return string(b)
//			}
//		case reflect.Slice:
//			var buf bytes.Buffer
//			buf.WriteByte('[')
//			for i := 0; i < rv.Len(); i++ {
//				if i > 0 {
//					buf.WriteByte(',')
//				}
//				buf.WriteString(JSONValue(rv.Index(i).Interface()))
//				rv.Index(i).Interface()
//			}
//			buf.WriteByte(']')
//			return buf.String()
//		case reflect.Map:
//			var buf bytes.Buffer
//			buf.WriteByte('{')
//			iter := rv.MapRange()
//			var i = 0
//			for iter.Next() {
//				if i > 0 {
//					buf.WriteByte(',')
//				}
//				i++
//				buf.WriteString("\"" + iter.Key().String() + "\":" + JSONValue(iter.Value().Interface()))
//			}
//			buf.WriteByte('}')
//			return buf.String()
//		}
//		s, e := JSON.MarshalToString(v)
//		if e != nil {
//			panic(e)
//		}
//		return s
//	}
//}

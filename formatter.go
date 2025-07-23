package logman

// func formatJSON(msg Message) (string, error) {
// 	s := `{"Fields":`
// 	keys := msg.Fields()
// 	switch len(keys) {
// 	default:
// 		s += "{"
// 		for _, key := range keys {
// 			switch val := msg.Value(key).(type) {
// 			case string:
// 				fmt.Println("string", val)
// 				s += fmt.Sprintf(`"%v":"%v",`, key, val)
// 			default:
// 				s += fmt.Sprintf(`"%v":%v,`, key, val)
// 			}

// 		}
// 		s = strings.TrimSuffix(s, ",") + "}"

// 	case 0:
// 		s += `null`
// 	}
// 	s += `}`
// 	return s, nil
// }

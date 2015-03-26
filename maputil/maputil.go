package maputil

import (
    "strings"
    "strconv"
    "sort"
    "github.com/shutterstock/go-stockutil/stringutil"
    _ "log"
    "reflect"
)


func StringKeys(input map[string]interface{}) []string {
    keys := make([]string, 0)

    for k, _ := range input {
        keys = append(keys, k)
    }

    return keys
}

func MapValues(input map[string]interface{}) []interface{} {
    values := make([]interface{}, 0)

    for _, value := range input {
        values = append(values, value)
    }

    return values
}

// Take a flat (non-nested) map keyed with fields joined on fieldJoiner and return a
// deeply-nested map
//
func DiffuseMap(data map[string]interface{}, fieldJoiner string) (map[string]interface{}, error) {
    output     := make(map[string]interface{})

//  get the list of keys and sort them because order in a map is undefined
    dataKeys := StringKeys(data)
    sort.Strings(dataKeys)

//  for each data item
    for _, key := range dataKeys {
        value, _ := data[key]
        keyParts := strings.Split(key, fieldJoiner)

        output = DeepSet(output, keyParts, value).(map[string]interface{})
    }

    return output, nil
}


// Take a deeply-nested map and return a flat (non-nested) map with keys whose intermediate tiers are joined with fieldJoiner
//
func CoalesceMap(data map[string]interface{}, fieldJoiner string) (map[string]interface{}, error) {
    return deepGetValues([]string{}, fieldJoiner, data), nil
}


func deepGetValues(keys []string, joiner string, data interface{}) map[string]interface{} {
    rv := make(map[string]interface{})
    
    if data != nil {
        switch reflect.TypeOf(data).Kind() {
        case reflect.Map:
            for k, v := range data.(map[string]interface{}){
                newKey := keys
                newKey = append(newKey, k)

                for kk, vv := range deepGetValues(newKey, joiner, v) {
                    rv[kk] = vv
                }
            }

        case reflect.Slice, reflect.Array:
            for i, value := range data.([]interface{}) {
                newKey := keys
                newKey = append(newKey, strconv.Itoa(i))

                for k, v := range deepGetValues(newKey, joiner, value){
                    rv[k] = v
                }
            }

        default:
            rv[strings.Join(keys, joiner)] = data
        }
    }

    return rv
}


func DeepGet(data interface{}, path []string, fallback interface{}) interface{} {
    current := data

    for i := 0; i < len(path); i++ {
        part := path[i]

        switch current.(type) {
    //  arrays
        case []interface{}:
            currentAsArray := current.([]interface{})

            if stringutil.IsInteger(part) {
                if partIndex, err := strconv.Atoi(part); err == nil {
                    if partIndex < len(currentAsArray) {
                        if value := currentAsArray[partIndex]; value != nil {
                            current = value
                            continue
                        }
                    }
                }
            }

            return fallback

    //  maps
        case map[string]interface{}:
            currentAsMap := current.(map[string]interface{})

            if value, ok := currentAsMap[part]; !ok {
                return fallback
            }else{
                current = value
            }
        }
    }

    return current
}



func DeepSet(data interface{}, path []string, value interface{}) interface{} {
    if len(path) == 0 {
        return data
    }

    var first = path[0]
    var rest    = make([]string, 0)

    if len(path) > 1 {
        rest = path[1:]
    }

//  Leaf Nodes
//    this is where the value we're setting actually gets set/appended
    if len(rest) == 0 {
        switch data.(type) {
        // ARRAY
        case []interface{}:
            return append(data.([]interface{}), value)

        // MAP
        case map[string]interface{}:
            dataMap := data.(map[string]interface{})
            dataMap[first] = value

            return dataMap
        }

    }else{
    //  Array Embedding
    //    this is where keys that are actually array indices get processed
    //  ================================
    //  is `first' numeric (an array index)
        if stringutil.IsInteger(rest[0]) {
            switch data.(type) {
            case map[string]interface{}:
              dataMap := data.(map[string]interface{})

          //  is the value at `first' in the map isn't present or isn't an array, create it
          //  -------->
              curVal, _ := dataMap[first]
 
              switch curVal.(type) {
              case []interface{}:
              default:
                  dataMap[first] = make([]interface{}, 0)
                  curVal, _ = dataMap[first]
              }
          //  <--------|


          //  recurse into our cool array and do awesome stuff with it
              dataMap[first] = DeepSet(curVal.([]interface{}), rest, value).([]interface{})
              return dataMap
            default:
              // log.Printf("WHAT %s/%s", first, rest)
            }


    //  Intermediate Map Processing
    //    this is where branch nodes get created and populated via recursion
    //    depending on the data type of the input `data', non-existent maps
    //    will be created and either set to `data[first]' (the map)
    //    or appended to `data[first]' (the array)
    //  ================================
        }else{
            switch data.(type) {
        //  handle arrays of maps
            case []interface{}:
                dataArray := data.([]interface{})

                if curIndex, err := strconv.Atoi(first); err == nil {
                    if curIndex >= len(dataArray) {
                        for add := len(dataArray); add <= curIndex; add++ {
                            dataArray = append(dataArray, make(map[string]interface{}))
                        }
                    }

                    if curIndex < len(dataArray) {
                        dataArray[curIndex] = DeepSet(dataArray[curIndex], rest, value)
                        return dataArray
                    }
                }

        //  handle good old fashioned maps-of-maps
            case map[string]interface{}:
                dataMap := data.(map[string]interface{})

            //  is the value at `first' in the map isn't present or isn't a map, create it
            //  -------->
                curVal, _ := dataMap[first]

                switch curVal.(type) {
                case map[string]interface{}:
                default:
                    dataMap[first] = make(map[string]interface{})
                    curVal, _ = dataMap[first]
                }
            //  <--------|

                dataMap[first] = DeepSet(dataMap[first], rest, value)
                return dataMap
            }
        }
    }

    return data
}

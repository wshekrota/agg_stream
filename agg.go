package main

import "strings"
import "bufio"
import "os"
import "encoding/json"
import "fmt"
import "strconv"

/*
   Author: Walt Shekrota wshekrota@icloud.com
   Name: agg

   Description:
   Binary data in json format 
   Input:
		{"id":121509,"market":5773,"price":1.234,"volume":1234.56,"is_buy":true}
		{"id":121510,"market":5774,"price":2.345,"volume":2345.67,"is_buy":false}
		{"id":121511,"market":5775,"price":3.456,"volume":3456.78,"is_buy":true}

   Output:
		Maps will be nested to maintain data agg.
		Format: map[string]map[string]float64{} (seemed easier for me to make it all float)
		Then you loop through the completed map encoding the data back to json.
		
*/

type Ret struct {
	slice []byte
	err   error
}

func main() {

	var output = map[string]map[string]float64{}

	// Read json from stdin passed from null logstash pipe
	//

	reader := bufio.NewScanner(os.Stdin)
	var retData Ret

	// input loop
	//
	for {
		if reader.Scan() {
			retData = Ret{reader.Bytes(), reader.Err()}
		} else {
			break
		}

		// Ignore jibber in the stream (non json)
		//
		if retData.slice[0:1][0] != byte('{') { continue }

		// Decoded data
		var Hash map[string]interface{}

		// ignore bad json
		if err := json.Unmarshal(retData.slice, &Hash); err != nil {
			fmt.Printf("input Unmarshal error: %v", err)
			continue
		}
		// Valid decode count and look for market key
		// If key then update else init
		//
		var market = fmt.Sprintf("%f",Hash["market"])
		market = strings.Split(market,".")[0]

		var buy float64 = 0.0
		if strconv.FormatBool(Hash["is_buy"].(bool))=="true" { buy = 1.0 }

		// Does market key exist in map?
		// yes?
		if _, ok := output[market]; ok {

			// could later delete items like count,count_buy and total_price if you don't want to encode those
			output[market]["count"]+=1.0
			if buy==1.0 {
				output[market]["count_buy"]++
			}
			output[market]["total_volume"] += Hash["volume"].(float64)
			output[market]["total_price"] += Hash["price"].(float64)
			output[market]["mean_price"] = output[market]["total_price"] / output[market]["count"]
			output[market]["mean_volume"] = output[market]["total_volume"] / output[market]["count"]
			output[market]["volume_weighted_average_price"] = output[market]["total_price"] * Hash["volume"].(float64) / output[market]["total_volume"]
			output[market]["percentage_buy"] = output[market]["count_buy"] / output[market]["count"] * 100
		} else {
			// No then initialize
			output[market] = map[string]float64{
				"count":                         1.,
				"count_buy":                     buy,
				"total_volume":                  Hash["volume"].(float64),
				"mean_price":                    Hash["price"].(float64),
				"mean_volume":                   Hash["volume"].(float64),
				"volume_weighted_average_price": Hash["price"].(float64),
				"percentage_buy":                buy}


		}

	}

	// re encode output map as JSON
	for key, element := range output {

		// Could put deletes here for the maps to make json output smaller
		// delete(output[key],"count")
		// delete(output[key],"count_buy")
		// delete(output[key],"total_price")
		line, _ := json.Marshal(element)
		fmt.Println(key,"=",string(line))
	}
}

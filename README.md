# agg_stream
### Aggregate a stream based on predetermined rules.
Stream is binary or byte and will decode from json to internal.

Will reveal following keys..

---

id 

market 

price 

volume 

is_buy

---

The key to aggregate on is market. Maintain a key | value for each unique market.

Value will be the map of aggregation for that market.

Keys in the per market map will be...
---

  total_volume +=
  
  total_price  +=
  
  mean_price   total_price / per market count
  
  mean_volume  total_volume / per market count
  
  volume_weighted_average_price total_price * volume / total_volume
  
  percentage_buy count_buy / count * 100
  
  count count += 1
  
  count_buy if buy { count_buy += 1 }
  
---

So overall structure is 

### map[string]map[string]float64{}.

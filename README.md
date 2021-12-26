# Logbot - simple telegram bot for logging arbitrary events

Supports writing to csv file or influxdb.

## Example
Oleg Vasilev, [12/26/21 10:21 PM]
test1

logbot test, [12/26/21 10:21 PM]
Item test1 received

logbot test, [12/26/21 10:21 PM]
option

Oleg Vasilev, [12/26/21 10:21 PM]
opt1

logbot test, [12/26/21 10:21 PM]
Option opt1 received

logbot test, [12/26/21 10:21 PM]
26 Dec 21 22:21 MSK

Oleg Vasilev, [12/26/21 10:21 PM]
done

logbot test, [12/26/21 10:21 PM]
Recorded
item: test1
param: opt1
comment: ""
ts: 2021-12-26T22:21:37.240956472+03:00

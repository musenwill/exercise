package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "pong")
	}

	http.HandleFunc("/ping", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
echo "GET http://localhost:8080/ping" | vegeta attack -duration=60s -rate=1000 | tee results.bin | vegeta report
Requests      [total, rate, throughput]         60000, 1000.02, 1000.02
Duration      [total, attack, wait]             59.999s, 59.999s, 141.18µs
Latencies     [min, mean, 50, 90, 95, 99, max]  95.992µs, 223.937µs, 175.87µs, 229.764µs, 258.75µs, 631.559µs, 79.786ms
Bytes In      [total, mean]                     240000, 4.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:60000
Error Set:
*/

/*
echo "GET http://localhost:8080/ping" | vegeta attack -duration=60s -rate=2000 | tee results.bin | vegeta report
Requests      [total, rate, throughput]         120000, 2000.03, 2000.02
Duration      [total, attack, wait]             59.999s, 59.999s, 176.217µs
Latencies     [min, mean, 50, 90, 95, 99, max]  85.947µs, 173.261µs, 168.482µs, 207.72µs, 224.347µs, 285.976µs, 6.51ms
Bytes In      [total, mean]                     480000, 4.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:120000
Error Set:
*/

/*
echo "GET http://localhost:8080/ping" | vegeta attack -duration=60s -rate=4000 | tee results.bin | vegeta report
Requests      [total, rate, throughput]         240000, 4000.01, 4000.00
Duration      [total, attack, wait]             1m0s, 1m0s, 237.003µs
Latencies     [min, mean, 50, 90, 95, 99, max]  92.35µs, 310.954µs, 183.116µs, 257.419µs, 352.916µs, 1.906ms, 91.616ms
Bytes In      [total, mean]                     960000, 4.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:240000
Error Set:
*/

/*
echo "GET http://localhost:8080/ping" | vegeta attack -duration=60s -rate=8000 | tee results.bin | vegeta report
Requests      [total, rate, throughput]         480000, 8000.03, 8000.00
Duration      [total, attack, wait]             1m0s, 1m0s, 225.414µs
Latencies     [min, mean, 50, 90, 95, 99, max]  75.454µs, 192.017µs, 172.478µs, 233.918µs, 282.294µs, 679.689µs, 9.196ms
Bytes In      [total, mean]                     1920000, 4.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:480000
Error Set:
*/

/*
echo "GET http://localhost:8080/ping" | vegeta attack -duration=60s -rate=16000 | tee results.bin | vegeta report
Requests      [total, rate, throughput]         960000, 15999.95, 15999.89
Duration      [total, attack, wait]             1m0s, 1m0s, 219.718µs
Latencies     [min, mean, 50, 90, 95, 99, max]  60.168µs, 317.836µs, 199.193µs, 531.413µs, 968.511µs, 2.231ms, 25.028ms
Bytes In      [total, mean]                     3840000, 4.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:960000
Error Set:
*/

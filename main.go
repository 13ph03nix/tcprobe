package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type probeArgs []string

func (p *probeArgs) Set(val string) error {
	*p = append(*p, val)
	return nil
}

func (p probeArgs) String() string {
	return strings.Join(p, ",")
}

func main() {

	// concurrency flag
	var concurrency int
	flag.IntVar(&concurrency, "c", 128, "set the concurrency level")

	// additional probes
	var probes probeArgs
	flag.Var(&probes, "p", "add additional port probe")

	// skip default probes flag
	var skipDefault bool
	flag.BoolVar(&skipDefault, "s", false, "skip the address built-in port check")

	// timeout flag
	var to int
	flag.IntVar(&to, "t", 10000, "timeout (milliseconds)")

	flag.Parse()

	timeout := time.Duration(to * 1000000)

	addresses := make(chan string)
	output := make(chan string)

	// Tcp scan workers
	var addressesWG sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		addressesWG.Add(1)

		go func() {
			for address := range addresses {
				if isListening(address, timeout) {
					output <- address
					continue
				}
			}
			addressesWG.Done()
		}()
	}

	// Output worker
	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {
		for o := range output {
			fmt.Println(o)
		}
		outputWG.Done()
	}()

	// Close the output channel when the tcp scan workers are done
	go func() {
		addressesWG.Wait()
		close(output)
	}()

	// accept address on stdin
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		address := sc.Text()
		host, _, err := net.SplitHostPort(address)
		if err != nil {
			host = address

		} else if !skipDefault {
			// submit address built-in port check
			addresses <- address
		}

		// Adding port templates
		// https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml
		xlarge := []string{
			"0", "1", "5", "7", "9", "11", "13", "15", "17", "18", "19", "20",
			"21", "22", "23", "25", "26", "35", "37", "39", "41", "42", "43",
			"49", "53", "56", "57", "70", "79", "80", "81", "88", "101", "102",
			"107", "109", "110", "111", "113", "115", "117", "118", "119",
			"135", "137", "138", "139", "143", "152", "153", "156", "158",
			"161", "162", "170", "179", "194", "201", "209", "213", "218",
			"220", "259", "264", "308", "311", "318", "323", "366", "369",
			"371", "383", "384", "387", "389", "401", "411", "412", "427",
			"443", "444", "445", "464", "465", "475", "497", "500", "502",
			"512", "513", "514", "515", "520", "524", "530", "531", "532",
			"540", "542", "543", "544", "546", "547", "548", "554", "556",
			"563", "587", "591", "593", "604", "631", "636", "639", "646",
			"647", "648", "652", "665", "674", "691", "692", "695", "699",
			"700", "701", "702", "706", "711", "712", "720", "749", "782",
			"829", "860", "873", "901", "911", "981", "989", "990", "991",
			"992", "993", "995", "1025", "1026", "1029", "1058", "1059",
			"1080", "1099", "1109", "1140", "1176", "1182", "1198", "1214",
			"1223", "1241", "1248", "1270", "1311", "1313", "1337", "1352",
			"1387", "1414", "1431", "1433", "1434", "1494", "1512", "1521",
			"1524", "1526", "1533", "1547", "1677", "1716", "1723", "1755",
			"1761", "1762", "1763", "1764", "1765", "1766", "1767", "1768",
			"1863", "1935", "1970", "1971", "1972", "1984", "1994", "1998",
			"2000", "2002", "2031", "2053", "2073", "2074", "2082", "2083",
			"2086", "2087", "2095", "2096", "2161", "2181", "2200", "2219",
			"2220", "2222", "2301", "2369", "2370", "2381", "2404", "2447",
			"2483", "2484", "2546", "2593", "2598", "2612", "2710", "2735",
			"2809", "2948", "2949", "2967", "3000", "3001", "3002", "3003",
			"3004", "3006", "3007", "3025", "3050", "3074", "3128", "3260",
			"3268", "3269", "3300", "3305", "3306", "3333", "3386", "3389",
			"3396", "3689", "3690", "3702", "3724", "3784", "3868", "3872",
			"3899", "3900", "3945", "4000", "4007", "4089", "4093", "4111",
			"4224", "4226", "4662", "4664", "4894", "4899", "5000", "5001",
			"5003", "5031", "5050", "5051", "5060", "5061", "5104", "5106",
			"5107", "5110", "5121", "5176", "5190", "5222", "5223", "5269",
			"5351", "5402", "5405", "5421", "5432", "5495", "5498", "5500",
			"5501", "5517", "5555", "5556", "5631", "5666", "5667", "5800",
			"5814", "5900", "6000", "6005", "6050", "6051", "6100", "6110",
			"6111", "6112", "6129", "6346", "6347", "6444", "6445", "6502",
			"6522", "6566", "6600", "6619", "6665", "6666", "6667", "6668",
			"6669", "6679", "6697", "6699", "6881", "6882", "6883", "6884",
			"6885", "6886", "6887", "6888", "6889", "6890", "6891", "6892",
			"6893", "6894", "6895", "6896", "6897", "6898", "6899", "6900",
			"6901", "6902", "6903", "6904", "6905", "6906", "6907", "6908",
			"6909", "6910", "6911", "6912", "6913", "6914", "6915", "6916",
			"6917", "6918", "6919", "6920", "6921", "6922", "6923", "6924",
			"6925", "6926", "6927", "6928", "6929", "6930", "6931", "6932",
			"6933", "6934", "6935", "6936", "6937", "6938", "6939", "6940",
			"6941", "6942", "6943", "6944", "6945", "6946", "6947", "6948",
			"6949", "6950", "6951", "6952", "6953", "6954", "6955", "6956",
			"6957", "6958", "6959", "6960", "6961", "6962", "6963", "6964",
			"6965", "6966", "6967", "6968", "6969", "6970", "6971", "6972",
			"6973", "6974", "6975", "6976", "6977", "6978", "6979", "6980",
			"6981", "6982", "6983", "6984", "6985", "6986", "6987", "6988",
			"6989", "6990", "6991", "6992", "6993", "6994", "6995", "6996",
			"6997", "6998", "6999", "7000", "7001", "7002", "7005", "7006",
			"7010", "7025", "7047", "7171", "7306", "7307", "7670", "7777",
			"8000", "8002", "8008", "8009", "8010", "8074", "8080", "8086",
			"8087", "8090", "8118", "8200", "8220", "8291", "8294", "8443",
			"8500", "8881", "8882", "8888", "9000", "9001", "9009", "9043",
			"9060", "9100", "9119", "9535", "9800", "10024", "10025", "10050",
			"10051", "10113", "10114", "10115", "10116", "12345", "12975",
			"13720", "13721", "13724", "13782", "13783", "15000", "16000",
			"16080", "19226", "19638", "19813", "20720", "22347", "22350",
			"24554", "25999", "26000", "30564", "31337", "31456", "31457",
			"31458", "32245", "37777", "43594", "43595",
		}
		large := []string{
			"21", "22", "23", "25", "80", "110", "137", "138", "139", "143",
			"443", "445", "1433", "3306", "3389", "8080",
		}

		// submit any additional port probes
		for _, p := range probes {
			switch p {
			case "xlarge":
				for _, port := range xlarge {
					addresses <- fmt.Sprintf("%s:%s", host, port)
				}
			case "large":
				for _, port := range large {
					addresses <- fmt.Sprintf("%s:%s", host, port)
				}
			default:
				addresses <- fmt.Sprintf("%s:%s", host, p)
			}
		}
	}

	close(addresses)

	// check there were no errors reading stdin (unlikely)
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}

	// Wait until the output waitgroup is done
	outputWG.Wait()
}

// TCP Connect scan
func isListening(address string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

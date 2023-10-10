package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const Version = "23.08.21"

type Server struct {
	Remarks    string `json:"remarks"`
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
}

type Response struct {
	Servers []Server `json:"servers"`
}

// getRandomServers returns a random selection of servers from the given slice
func GetRandomServers(servers []Server, limit int) []Server {
	rand.Seed(time.Now().UnixNano())

	randomServers := make([]Server, 0, limit)
	indexes := rand.Perm(len(servers))

	for i := 0; i < limit; i++ {
		randomIndex := indexes[i]
		randomServers = append(randomServers, servers[randomIndex])
	}

	return randomServers
}

// LogInit initializes the logrus logger
func LogInit() {
	bytesWriter := &bytes.Buffer{}
	stdoutWriter := os.Stdout
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z",
		FullTimestamp:   true})
	log.SetOutput(io.MultiWriter(bytesWriter, stdoutWriter))
	log.SetLevel(log.InfoLevel)
}

func main() {
	// Define command line flags
	hlepFlag := flag.Bool("h", false, "--help")
	urlSubscription := flag.String("u", "http://www.subcriptionurl.com", "Subscription URL")
	serverLimit := flag.Int("l", 10, "Number of servers to retrieve (default 0)")
	serverProxyPort := flag.Int("p", 11080, "TCP proxy port")
	serverRedPort := flag.Int("r", 11081, "RED proxy port")
	proxyStrategy := flag.String("s", "fifo", "Proxy strategy (round|rand|fifo|hash)")
	proxyFailTimeout := flag.Int("t", 600, "Proxy fail timeout (in seconds)")
	proxyMaxFails := flag.Int("m", 1, "Max failed count for proxy")
	printVersion := flag.Bool("V", false, "print version")
	filterSubscribedKeywords := flag.String("f", "套餐|重置|剩余|更新", "Filter subscriptions containing keywords")
	outputPath := flag.String("o", "config.yml", "Output file path")

	if *hlepFlag {
		flag.Usage()
		return
	}

	if *printVersion {
		fmt.Println(fmt.Printf("Current version: %s", Version))
		return
	}

	flag.Parse()

	if *urlSubscription == "" {
		flag.Usage()
		log.Fatalln("please check if your parameter inputs are correct.")
	}

	LogInit()

	log.WithFields(log.Fields{
		"serverLimit":      *serverLimit,
		"serverProxyPort":  *serverProxyPort,
		"serverRedPort":    *serverRedPort,
		"proxyStrategy":    *proxyStrategy,
		"proxyFailTimeout": *proxyFailTimeout,
		"proxyMaxFails":    *proxyMaxFails,
		"outputPath":       *outputPath,
	}).Info("build gost config file")

	resp, err := http.Get(*urlSubscription)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response Response
	json.Unmarshal(body, &response)

	limit := *serverLimit
	if *serverLimit > len(response.Servers) {
		limit = len(response.Servers)
	}

	servers := GetRandomServers(response.Servers, limit)

	log.WithFields(log.Fields{"servers": len(servers)}).Info("subsciption success")

	nodes := make([]map[string]interface{}, len(servers))
outerLoop:
	for i, server := range servers {
		keys := strings.Split(*filterSubscribedKeywords, "|")
		for _, key := range keys {
			if strings.Contains(server.Remarks, key) {
				log.WithFields(log.Fields{"remarks": server.Remarks, "key": key}).Info("name contains key, igonre")
				continue outerLoop
			}
		}
		nodes[i] = map[string]interface{}{
			"name": server.Remarks,
			"addr": fmt.Sprintf("%s:%d", server.Server, server.ServerPort),
			"connector": map[string]interface{}{
				"type": "ss",
				"auth": map[string]string{
					"username": server.Method,
					"password": server.Password,
				},
			},
		}
		log.WithFields(log.Fields{"name": server.Remarks}).Info("add subscription success")
	}

	yamlData := map[string]interface{}{
		"services": []map[string]interface{}{
			{
				"name": strconv.Itoa(*serverRedPort),
				"addr": fmt.Sprintf(":%d", *serverRedPort),
				"handler": map[string]string{
					"type":  "red",
					"chain": "chain-0",
				},
				"listener": map[string]string{
					"type": "red",
				},
			},
			{
				"name": strconv.Itoa(*serverProxyPort),
				"addr": fmt.Sprintf(":%d", *serverProxyPort),
				"handler": map[string]string{
					"type":  "auto",
					"chain": "chain-0",
				},
				"listener": map[string]string{
					"type": "tcp",
				},
			},
		},
		"chains": []map[string]interface{}{
			{
				"name": "chain-0",
				"hops": []map[string]interface{}{{
					"name": "hop-0",
					"selector": map[string]interface{}{
						"strategy":    proxyStrategy,
						"maxFails":    strconv.Itoa(*proxyMaxFails),
						"failTimeout": fmt.Sprintf("%ds", *proxyFailTimeout),
					},
					"nodes": nodes,
				},
				},
			},
		},
	}

	file, err := os.Create(*outputPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(yamlData)
	if err != nil {
		log.Fatalln(err)
	}

}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type DnsRecordEntity struct {
	Data string `json:"data"`
	TTL  int    `json:"ttl"`
}

func putDNS(domain string, domainType string, domainValue string, domainNameArray []string, shopperID string, apiKey string, apiSecret string, proxyUrl string) {
	var dnsRecord = DnsRecordEntity{
		Data: domainValue,
		TTL:  1800,
	}
	var putData = [1]DnsRecordEntity{dnsRecord}
	var putDataJson, _ = json.Marshal(putData)

	// 创建一个自定义的HTTP Transport
	var tr *http.Transport
	if proxyUrl != "" {
		proxy, _ := url.Parse(proxyUrl)
		tr = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	} else {
		tr = &http.Transport{}
	}

	// 创建一个自定义的HTTP客户端
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	for _, v := range domainNameArray {
		domainUrl := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/%s/%s", domain, domainType, v)
		req, err := http.NewRequest("PUT", domainUrl, bytes.NewBuffer(putDataJson))
		if err != nil {
			log.Error(err)
			continue
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-Shopper-Id", shopperID)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", apiKey, apiSecret))

		// 使用自定义的HTTP客户端发送请求
		resp, err := client.Do(req)
		if err != nil {
			log.Error(err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.WithFields(log.Fields{
				"domain":     v,
				"domainIP":   domainValue,
				"domainType": domainType,
				"statusCode": resp.StatusCode,
			}).Warning("update failed")
			continue
		}
		log.WithFields(log.Fields{
			"domain":     v,
			"domainIP":   domainValue,
			"domainType": domainType,
			"statusCode": resp.StatusCode,
		}).Info("update success")
	}

}

func getIP(ipv6 bool) (string, error) {
	url := "https://4.ipw.cn"
	if ipv6 {
		url = "https://6.ipw.cn"
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//Convert the body to type string
	var ip = string(body)

	return ip, nil
}

func loggerInit() {
	bytesWriter := &bytes.Buffer{}
	stdoutWriter := os.Stdout
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z",
		FullTimestamp:   true})
	log.SetOutput(io.MultiWriter(bytesWriter, stdoutWriter))
	log.SetLevel(log.InfoLevel)
}

func main() {
	h := flag.Bool("help", false, "--help")
	flagDomain := flag.String("domain", "", "Domain name")
	flagType := flag.String("type", "A", "Domain name type, default is A type.")
	flagName := flag.String("name", "", "Sub domain name, separate multiple sub domains by comma(,).")
	flagRecord := flag.String("record", "", "Domain name corresponding value, if empty, automatically obtain IPV4/IPV6 value.")
	flagShopperID := flag.String("shopperid", "", "Godaddy shopper id for api.")
	flagKey := flag.String("key", "", "The key of godaddy api.")
	flagSecret := flag.String("secret", "", "The secret of godaddy api.")
	proxyUrl := flag.String("proxy", "", "Proxy for HTTP requests.")
	flag.CommandLine.SortFlags = false
	flag.Parse()

	if *h {
		flag.Usage()
		return
	}

	if *flagDomain == "" || *flagName == "" || *flagShopperID == "" || *flagKey == "" || *flagSecret == "" {
		flag.Usage()
		log.Fatalln("please check if your parameter inputs are correct.")
	}

	loggerInit()

	var nameArray []string = strings.Split(*flagName, ",")

	var domainValue string = *flagRecord
	if domainValue == "" {
		domainValue, _ = getIP(*flagType == "AAAA")
		if domainValue == "" {
			log.Fatalln("get ip value error")
		}
	}
	log.WithField("IP", domainValue).Info("preparing to update")

	putDNS(*flagDomain, *flagType, domainValue, nameArray, *flagShopperID, *flagKey, *flagSecret, *proxyUrl)

}

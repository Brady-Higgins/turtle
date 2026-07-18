package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// instance type : cost per hour
// 7/16/2026
var costs = map[string]float64{
	"t4g.nano":    0.0042,
	"t4g.micro":   0.0084,
	"t4g.small":   0.0168,
	"t4g.medium":  0.0336,
	"t4g.large":   0.0672,
	"t4g.xlarge":  0.1344,
	"t4g.2xlarge": 0.2688,
	"t3.nano":     0.0052,
	"t3.micro":    0.0104,
	"t3.small":    0.0208,
	"t3.medium":   0.0416,
	"t3.large":    0.0832,
	"t3.xlarge":   0.1664,
	"t3.2xlarge":  0.3328,
	"t3a.nano":    0.0047,
	"t3a.micro":   0.0094,
	"t3a.small":   0.0188,
	"t3a.medium":  0.0376,
	"t3a.large":   0.0752,
	"t3a.xlarge":  0.1504,
}

// reads main.tf file to find instance type
func FindServerType() (string, error) {
	if _, err := os.Stat("main.tf"); errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	file, err := os.Open("main.tf")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		if _, ok := costs[strings.Trim(word, "\"")]; ok {
			return strings.Trim(word, "\""), nil
		}
	}
	return "", errors.New("server instance type not found")
}

func ServerToCost(instanceType string) float64 {
	if _, ok := costs[instanceType]; ok {
		return costs[instanceType]
	}
	return -1
}

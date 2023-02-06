package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tinygg/gofaker"
)

type kafkaClient struct {
	Topic string
	P     *kafka.Producer
	C     *kafka.Consumer
}

type stringFileWriter struct {
	F *os.File
}

func (f stringFileWriter) Write(p []byte) (int, error) {
	return f.F.WriteString(string(p) + "\n")
}

func (k kafkaClient) Write(p []byte) (int, error) {
	err := k.P.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &k.Topic, Partition: kafka.PartitionAny},
		Key:            []byte(nil),
		Value:          p,
	}, nil)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

var kc kafkaClient
var fsw stringFileWriter

func GenData(templateFileName string, isKafka bool, isFile bool, topic string, limit uint64) {
	jsonFile, err := os.Open(templateFileName)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer jsonFile.Close()

	// Load a template
	bytes, _ := ioutil.ReadAll(jsonFile)
	template := string(bytes)
	var templateMap interface{}
	err = json.Unmarshal([]byte(template), &templateMap)
	if err != nil {
		log.Fatal(err)
	}

	if isKafka {
		//Kafka Producer
		p := kafkaProducer()
		defer p.Flush(1000)
		defer p.Close()

		kc.Topic = topic
		kc.P = p
	}

	if isFile {
		// Prepare file to write to
		_ = os.Mkdir("./data", os.ModePerm)
		f, err := os.Create("./data/data.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Sync()
		defer f.Close()
		fsw.F = f
	}

	startTime := time.Now()

	var count uint64 = 0
	for {
		templateMapCopy := CopyMap(templateMap.(map[string]interface{}))
		fakeData(templateMapCopy)
		ob, _ := json.Marshal(templateMapCopy)

		if limit != 0 {
			if count < limit {
				if isFile {
					fsw.Write(ob)
				}
				if isKafka {
					kc.Write(ob)
				}
			} else {
				break
			}
		} else {
			if isKafka {
				kc.Write(ob)
			}

			if count < 1000 && isFile {
				fsw.Write(ob)
			}

			if count >= 1000 && !isKafka {
				break
			}
		}
		count += 1
	}

	elapsed := time.Since(startTime)

	log.Printf("Data generation took: %s", elapsed)
	log.Printf("Generated %d data elements", count)
}

func kafkaProducer() *kafka.Producer {
	configFile := "./conf/kafka.properties"
	conf := ReadConfig(configFile)

	p, err := kafka.NewProducer(&conf)

	if err != nil {
		log.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	return p
}

func fakeData(m map[string]interface{}) {
	for k, v := range m {
		if _, ok := v.(map[string]interface{}); ok {
			fakeData(v.(map[string]interface{}))
		} else {
			if _, isString := v.(string); isString && strings.HasPrefix(v.(string), "fake:{") {
				str := v.(string)
				funcName := strings.TrimSpace(str[6:strings.IndexRune(str, '}')])

				if strings.HasPrefix(funcName, "number:") ||
					strings.HasPrefix(funcName, "latituderange:") ||
					strings.HasPrefix(funcName, "longituderange:") ||
					strings.HasPrefix(funcName, "float32range:") ||
					strings.HasPrefix(funcName, "float64range:") {

					params := getParams(funcName)
					funcName = funcName[:strings.IndexRune(funcName, ':')]
					info := gofaker.GetFuncLookup(funcName)
					val2, _ := info.Call(&params, info)
					m[k] = val2
				} else {

					info := gofaker.GetFuncLookup(funcName)
					val2, _ := info.Call(nil, info)
					m[k] = val2
				}

			}
		}
	}
}

func getParams(str string) map[string][]string {
	params := make(map[string][]string)
	min := make([]string, 1)
	min[0] = str[strings.IndexRune(str, ':')+1 : strings.IndexRune(str, ',')]
	max := make([]string, 1)
	max[0] = str[strings.IndexRune(str, ',')+1:]
	params["min"] = min
	params["max"] = max
	return params
}

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

func ReadConfig(configFile string) kafka.ConfigMap {

	m := make(map[string]kafka.ConfigValue)

	file, err := os.Open(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			parameter := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[parameter] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read file: %s", err)
		os.Exit(1)
	}

	return m

}

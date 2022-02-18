/**
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type Cadf struct {
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
}

func main() {
	// Checks that the minimal number of arguments is passed
	if len(os.Args) < 4 {
		log.Fatal("Usage: hostname username password topic")
		os.Exit(1)
	}

	// Hostname: hostname:port
	// Example: "localhost:9092"
	hostname := os.Args[1]

	username := os.Args[2]
	password := os.Args[3]

	// GroupId that is used to retreive the messages.
	groupId := "fhir-audit-consumer-group"

	// Kafka Topic:
	// Example: FHIR_AUDIT
	topic := os.Args[4]

	mechanism := plain.Mechanism{
		Username: username,
		Password: password,
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
		TLS: &tls.Config{
			// Intentionally insecure as it's dev-only
			InsecureSkipVerify: true,
		},
	}

	fmt.Printf("Hostname: '%s' \n", string(hostname))

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(hostname, ","),
		GroupID: groupId,
		Topic:   topic,
		Dialer:  dialer,
	})

	// Look backwards one hour in order to get previously sent messages
	startTime := time.Now().Add(-time.Hour)
	endTime := time.Now().Add(+time.Hour)
	r.SetOffsetAt(context.Background(), startTime)

	for {
		m, err := r.ReadMessage(context.Background())

		if err != nil {
			log.Fatal("failed to close reader:", err)
			break
		}
		if m.Time.After(endTime) {
			break
		}

		fmt.Printf("\n---------- \n message at offset %d: %s = \n %s", m.Offset, string(m.Key), string(m.Value))

		var result Cadf
		json.Unmarshal([]byte(m.Value), &result)
		attachments := result.Attachments
		for _, value := range attachments {
			content := value.Content
			audit, _ := b64.StdEncoding.DecodeString(string(content))
			log.Println("The audit details are: ", string(audit))
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Finished")
}

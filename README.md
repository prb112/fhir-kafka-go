# Demonstration: Streaming the FHIR Audit from the IBM FHIR Server with Go

The [IBM FHIR Server](https://ibm.github.io/FHIR) supports audit events for FHIR operations (CREATE-READ-UPDATE-DELETE-OPERATION) in [Cloud Auditing Data Federation (CADF)](https://www.dmtf.org/standards/cadf) and [HL7 FHIR AuditEvent](https://www.hl7.org/fhir/auditevent.html) and pushing the events to an [Apache Kafka](https://kafka.apache.org/) backend. You can read more about it in [another post I made](https://bastide.org/2021/04/21/question-how-to-setup-the-ibm-fhir-server-with-fhir-audit/).

The blog is located at [recipe-streaming-the-fhir-audit-from-the-ibm-fhir-server-with-go](https://bastide.org/2022/02/18/recipe-streaming-the-fhir-audit-from-the-ibm-fhir-server-with-go/).

You'll need golang installed. 

1. Clone the repository
git clone https://github.com/prb112/fhir-kafka-go.git

2. Build the Go guild 
go build

3. Run with Kafka and EventStreams

# Documentation

> Usage: hostname username password topic

It'll retrieve messages from one hour back, and up to one hour after start.

## Example Simple

> ./fhir-kafka-go broker-4-*****.kafka.svc02.us-east.eventstreams.cloud.ibm.com:9093,broker-0-*****.kafka.svc02.us-east.eventstreams.cloud.ibm.com:9093,broker-2-*****.kafka.svc02.us-east.eventstreams.cloud.ibm.com:9093,broker-3-*****.kafka.svc02.us-east.eventstreams.cloud.ibm.com:9093,broker-1-*****.kafka.svc02.us-east.eventstreams.cloud.ibm.com:9093,broker-5-*****.kafka.svc02.us-east.eventstreams.cloud.ibm.com:9093 "token" "Wbot2tK-****-" fhir-consumer-group FHIR_AUDIT

## Example More Complicated
./fhir-kafka-go $(cat ibm-creds.json | jq -r '.kafka_brokers_sasl | join(",")') "token" $(cat ibm-creds.json | jq -r '.password') FHIR_AUDIT
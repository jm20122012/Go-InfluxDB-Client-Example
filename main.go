package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading environment variables: ", err)
	}

	influxUrl := fmt.Sprintf("https://%s:8086", os.Getenv("INFLUXDB_SERVER_IP"))
	influxToken := os.Getenv("INFLUXDB_API_TOKEN")
	influxBucket := os.Getenv("INFLUXDB_BUCKET")
	influxOrg := os.Getenv("INFLUXDB_ORG")

	// Create a new client using an InfluxDB server base URL and an authentication token

	// This is the default client instantiation.  However, due to server using
	// a self-signed TLS certificate without using SANs, this throws an error
	// client := influxdb2.NewClient(influxUrl, influxToken)

	// To skip certificate verification, create the client like this:

	// Create a new TLS configuration with certificate verification disabled
	// This is not recommended though.  Only use for testing until a new
	// TLS certificate can be created with SANs
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	client := influxdb2.NewClientWithOptions(influxUrl, influxToken, influxdb2.DefaultOptions().SetTLSConfig(tlsConfig))

	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(influxOrg, influxBucket)

	// Create point using full params constructor
	// p := influxdb2.NewPoint("GoTestData",
	// 	map[string]string{"unit": "temperature"},
	// 	map[string]interface{}{"avg": 24.5, "max": 45.0},
	// 	time.Now())

	// // write point immediately
	// writeAPI.WritePoint(context.Background(), p)

	// Create point using fluent style
	p := influxdb2.NewPointWithMeasurement("GoTestData").
		AddTag("unit", "temperature").
		AddField("avg", 23.2).
		AddField("max", 45.0).
		SetTime(time.Now())

	err = writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		panic(err)
	}

	// Or write directly line protocol
	// line := fmt.Sprintf("stat,unit=temperature avg=%f,max=%f", 23.5, 45.0)
	// err = writeAPI.WriteRecord(context.Background(), line)
	// if err != nil {
	// 	panic(err)
	// }

	// Ensures background processes finishes
	client.Close()
}

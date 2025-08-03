package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

type Notifier struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Color     string `json:"color"`
}

type VisualizationData struct {
	Title    string    `json:"title"`
	Notifier []Notifier `json:"notifier"`
}

type DataModel struct {
	Data   VisualizationData `json:"data"`
	Signature string          `json:"signature"`
}

func generateColor() string {
	c, err := colorful.RandomColor()
	if err != nil {
		log.Fatal(err)
	}
	return c.Hex()
}

func generateNotifier(id string, message string) Notifier {
	return Notifier{
		ID:        id,
		Message:   message,
		Timestamp: time.Now().Unix(),
		Color:     generateColor(),
	}
}

func generateVisualizationData(title string, notifier []Notifier) VisualizationData {
	return VisualizationData{
		Title:    title,
		Notifier: notifier,
	}
}

func generateDataModel(title string, message string) DataModel {
	notifier := generateNotifier("1", message)
	data := generateVisualizationData(title, []Notifier{notifier})

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: "Notifier",
		},
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, &template, priv)
	if err != nil {
		log.Fatal(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, csr)
	if err != nil {
		log.Fatal(err)
	}

	return DataModel{
		Data:     data,
		Signature: fmt.Sprintf("%x", signature),
	}
}

func main() {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		model := generateDataModel("Secure Data Visualization", "New Data Available")
		fmt.Fprint(w, toJSON(model))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func toJSON(data interface{}) string {
	// implement your favorite JSON marshaller here
	return ""
}
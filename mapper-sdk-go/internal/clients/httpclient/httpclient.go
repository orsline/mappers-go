// Package httpclient implements HTTP client initialization and message processing
package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/gorilla/mux"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/internal/config"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/internal/httpadapter"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/pkg/di"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"time"
)

// HttpClient is structure used to init HttpClient
type HttpClient struct {
	IP             string
	Port           string
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	server         *http.Server
	restController *httpadapter.RestController
}

// NewHttpClient initializes a new Http client instance
func NewHttpClient(dic *di.Container) *HttpClient {
	return &HttpClient{
		IP:             "0.0.0.0",
		Port:           "1215",
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
		restController: httpadapter.NewRestController(mux.NewRouter(), dic),
	}
}

// Init is a method to construct HTTP server
func (hc *HttpClient) Init(c config.Config) error {
	hc.restController.InitRestRoutes()
	if c.Http.CaCert == "" {
		hc.server = &http.Server{
			Addr:         hc.IP + ":" + hc.Port,
			WriteTimeout: hc.WriteTimeout,
			ReadTimeout:  hc.ReadTimeout,
			Handler:      hc.restController.Router,
		}
	} else {
		// Enable two-way authentication http tls
		caCrtPath := c.Http.CaCert
		pool := x509.NewCertPool()
		crt, err := ioutil.ReadFile(caCrtPath)
		if err != nil {
			klog.Errorf("Failed to read cert %s:%v", caCrtPath, err)
			return err
		}
		pool.AppendCertsFromPEM(crt)
		hc.server = &http.Server{
			Addr:         hc.IP + ":" + hc.Port,
			WriteTimeout: hc.WriteTimeout,
			ReadTimeout:  hc.ReadTimeout,
			Handler:      hc.restController.Router,
			TLSConfig: &tls.Config{
				ClientCAs: pool,
				// check client certificate file
				ClientAuth: tls.RequireAndVerifyClientCert,
			},
		}
	}
	klog.V(1).Info("HttpServer Start......")
	go func() {
		_, err := hc.Receive(c)
		if err != nil {
			klog.Errorf("Http Receive error:%v", err)
		}
	}()
	return nil
}

// UnInit is a method to close http server
func (hc *HttpClient) UnInit() {
	err := hc.server.Close()
	if err != nil {
		klog.Error("Http server close err:", err.Error())
		return
	}
}

// Send no messages need to be sent
func (hc *HttpClient) Send(message interface{}) error {
	return nil
}

// Receive http server start listen
func (hc *HttpClient) Receive(c config.Config) (interface{}, error) {
	if c.Http.CaCert == "" && c.Http.Cert == "" && c.Http.PrivateKey == "" {
		err := hc.server.ListenAndServe()
		if err != nil {
			return nil, err
		}
	} else if c.Http.Cert != "" && c.Http.PrivateKey != "" {
		serverCrtPath := c.Http.Cert
		serverKeyPath := c.Http.PrivateKey
		err := hc.server.ListenAndServeTLS(serverCrtPath,
			serverKeyPath)
		if err != nil {
			klog.Error("HTTP Server exited...")
			return "", err
		}
	} else {
		err := errors.New("the certificate file provided is incomplete or does not match")
		return "", err
	}
	return "", nil
}

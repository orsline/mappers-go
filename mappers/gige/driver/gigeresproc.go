package driver

import (
	"crypto/tls"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
	"runtime"
	"strconv"
	"sync"
)

const (
	mqttInitFlag = false
)

type MqttProtocolConfig struct {
	ProtocolName       string `json:"protocolName"`
	ProtocolConfigData `json:"configData"`
}

type ProtocolConfigData struct {
}

type MqttProtocolCommonConfig struct {
	CommonCustomizedValues `json:"customizedValues"`
}

type CommonCustomizedValues struct {
	IP    string `json:"IP"`
	Port  int    `json:"port"`
	Topic string `json:"topic"`
}

var ErrConfigCert = errors.New("both certification and private key must be provided")

// MQTT Realize the structure of random number
type MQTT struct {
	mutex                sync.Mutex
	rotocolConfig        MqttProtocolConfig
	protocolCommonConfig MqttProtocolCommonConfig
	Client               MqttClient
}

// Config is the modbus mapper configuration.
type Config struct {
	Mqtt      Mqtt   `yaml:"mqtt,omitempty"`
	HTTP      HTTP   `yaml:"http,omitempty"`
	Configmap string `yaml:"configmap"`
}

// Mqtt is the Mqtt configuration.
type Mqtt struct {
	ServerAddress string `yaml:"server,omitempty"`
	ServerName    string `yaml:"servername"`
	Username      string `yaml:"username,omitempty"`
	Password      string `yaml:"password,omitempty"`
	ClientID      string `yaml:"clientId"`
	Cert          string `yaml:"certification,omitempty"`
	PrivateKey    string `yaml:"privatekey,omitempty"`
	CaCert        string `yaml:"caCert,omitempty"`
}

// HTTP is the HTTP configuration
type HTTP struct {
	CaCert     string `yaml:"caCert,omitempty"`
	Cert       string `yaml:"certification,omitempty"`
	PrivateKey string `yaml:"privatekey,omitempty"`
}

// MqttClient is parameters for Mqtt client.
type MqttClient struct {
	Qos        byte
	Retained   bool
	IP         string
	User       string
	Passwd     string
	Cert       string
	PrivateKey string
	Client     mqtt.Client
}

// newTLSConfig new TLS configuration.
// Only one side check. Mqtt broker check the cert from client.
func newTLSConfig(certFile string, privateKey string) (*tls.Config, error) {
	// Import client certificate/key pair

	cert, err := tls.LoadX509KeyPair(certFile, privateKey)
	if err != nil {
		return nil, err
	}
	// Create tls.Config with desired tls properties
	return &tls.Config{
		// ClientAuth = whether to requests cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}, nil
}

// Connect used to the Mqtt server.
func (mc *MqttClient) Connect() error {
	opts := mqtt.NewClientOptions().AddBroker(mc.IP).SetClientID("").SetCleanSession(true)
	if mc.Cert != "" {
		tlsConfig, err := newTLSConfig(mc.Cert, mc.PrivateKey)
		if err != nil {
			return err
		}
		opts.SetTLSConfig(tlsConfig)
		opts.SetUsername(mc.User)
		opts.SetPassword(mc.Passwd)
	} else {
		opts.SetUsername(mc.User)
		opts.SetPassword(mc.Passwd)
	}
	mc.Client = mqtt.NewClient(opts)
	// The token is used to indicate when actions have completed.
	if tc := mc.Client.Connect(); tc.Wait() && tc.Error() != nil {
		return tc.Error()
	}
	mc.Qos = 0          // At most 1 time
	mc.Retained = false // Not retained
	return nil
}

// Publish a Mqtt message.
func (mc *MqttClient) Publish(topic string, payload interface{}) error {
	if tc := mc.Client.Publish(topic, mc.Qos, mc.Retained, payload); tc.Wait() && tc.Error() != nil {
		return tc.Error()
	}
	return nil
}

func (d *MQTT) newClient() {

	server := d.protocolCommonConfig.IP + ":" + strconv.Itoa(d.protocolCommonConfig.Port)
	d.Client = MqttClient{
		IP:         server,
		User:       "",
		Passwd:     "",
		Cert:       "",
		PrivateKey: "",
	}
	err := d.Client.Connect()
	if err != nil {
		fmt.Printf("Failed to init mqtt client ,error:%v\n", err)
		os.Exit(1)
	}
	fmt.Println(server, " connect Successful")
	fmt.Println("Subscribe Successful")
}

// Parse to parse the configuration file. If failed, return error.
func (c *Config) Parse() error {
	var level klog.Level
	var loglevel string
	var configFile string
	// -config-file /home/xxx
	sysType := runtime.GOOS
	defaultConfigFile := "../res/config.yaml"
	pflag.StringVar(&loglevel, "v", "1", "log level")
	pflag.StringVar(&configFile, "config-file", defaultConfigFile, "Config file name")
	pflag.StringVar(&c.Mqtt.ServerAddress, "mqtt-address", c.Mqtt.ServerAddress, "MQTT broker address")
	pflag.StringVar(&c.Mqtt.Username, "mqtt-username", c.Mqtt.Username, "username")
	pflag.StringVar(&c.Mqtt.Password, "mqtt-password", c.Mqtt.Password, "password")
	pflag.StringVar(&c.Mqtt.Cert, "mqtt-certification", c.Mqtt.Cert, "certification file path")
	pflag.StringVar(&c.Mqtt.PrivateKey, "mqtt-privatekey", c.Mqtt.PrivateKey, "private key file path")
	pflag.Parse()
	cf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return errors.New("config.yaml not found," + err.Error())
	}
	if err = yaml.Unmarshal(cf, c); err != nil {
		return errors.New("yaml.Unmarshal error:," + err.Error())
	}
	if err = level.Set(loglevel); err != nil {
		return errors.New("set loglevel error:," + err.Error())
	}
	if c.Mqtt.Cert != "" && c.Mqtt.PrivateKey == "" {
		klog.V(1).Info("The PrivateKey path is empty,", ErrConfigCert.Error())
	} else if c.Mqtt.Cert == "" && c.Mqtt.PrivateKey != "" {
		klog.V(1).Info("The CertPath is empty,", ErrConfigCert.Error())
	} else if c.Mqtt.Cert == "" && c.Mqtt.PrivateKey == "" {
		klog.V(1).Info("The connection is not secure,if you want to be secure,", ErrConfigCert.Error())
	}
	return nil
}

func (gigEClient *GigEVisionDevice) InitMqttClient() {

}

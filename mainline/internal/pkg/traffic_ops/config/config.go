package config

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/apache/trafficcontrol/internal/pkg/traffic_ops/riaksvc"
	"github.com/apache/trafficcontrol/pkg/log"
	"github.com/apache/trafficcontrol/pkg/rfc"
	"github.com/apache/trafficcontrol/pkg/util"
	"github.com/basho/riak-go-client"
)

// Config reflects the structure of the cdn.conf file
type Config struct {
	URL               *url.URL `json:"-"`
	CertPath          string   `json:"-"`
	KeyPath           string   `json:"-"`
	ConfigHypnotoad   `json:"hypnotoad"`
	ConfigTrafficOps  `json:"traffic_ops"`
	ConfigTO          *ConfigTO   `json:"to"`
	SMTP              *ConfigSMTP `json:"smtp"`
	ConfigPortal      `json:"portal"`
	ConfigLetsEncrypt `json:"lets_encrypt"`
	DB                ConfigDatabase `json:"db"`
	Secrets           []string       `json:"secrets"`
	// NOTE: don't care about any other fields for now..
	RiakAuthOptions  *riak.AuthOptions
	RiakEnabled      bool
	ConfigLDAP       *ConfigLDAP
	LDAPEnabled      bool
	LDAPConfPath     string `json:"ldap_conf_location"`
	ConfigInflux     *ConfigInflux
	InfluxEnabled    bool
	InfluxDBConfPath string `json:"influxdb_conf_path"`
	Version          string
	UseIMS           bool `json:"use_ims"`
}

// ConfigHypnotoad carries http setting for hypnotoad (mojolicious) server
type ConfigHypnotoad struct {
	Listen []string `json:"listen"`
	// NOTE: don't care about any other fields for now..
}

// ConfigTrafficOps carries settings specific to traffic_ops server
type ConfigTrafficOps struct {
	// Deprecated in 5.0
	Insecure bool `json:"insecure"`
	// end deprecated
	Port                     string                     `json:"port"`
	ProxyTimeout             int                        `json:"proxy_timeout"`
	ProxyKeepAlive           int                        `json:"proxy_keep_alive"`
	ProxyTLSTimeout          int                        `json:"proxy_tls_timeout"`
	ProxyReadHeaderTimeout   int                        `json:"proxy_read_header_timeout"`
	ReadTimeout              int                        `json:"read_timeout"`
	RequestTimeout           int                        `json:"request_timeout"`
	ReadHeaderTimeout        int                        `json:"read_header_timeout"`
	WriteTimeout             int                        `json:"write_timeout"`
	IdleTimeout              int                        `json:"idle_timeout"`
	LogLocationError         string                     `json:"log_location_error"`
	LogLocationWarning       string                     `json:"log_location_warning"`
	LogLocationInfo          string                     `json:"log_location_info"`
	LogLocationDebug         string                     `json:"log_location_debug"`
	LogLocationEvent         string                     `json:"log_location_event"`
	MaxDBConnections         int                        `json:"max_db_connections"`
	DBMaxIdleConnections     int                        `json:"db_max_idle_connections"`
	DBConnMaxLifetimeSeconds int                        `json:"db_conn_max_lifetime_seconds"`
	BackendMaxConnections    map[string]int             `json:"backend_max_connections"`
	DBQueryTimeoutSeconds    int                        `json:"db_query_timeout_seconds"`
	Plugins                  []string                   `json:"plugins"`
	PluginConfig             map[string]json.RawMessage `json:"plugin_config"`
	PluginSharedConfig       map[string]interface{}     `json:"plugin_shared_config"`
	ProfilingEnabled         bool                       `json:"profiling_enabled"`
	ProfilingLocation        string                     `json:"profiling_location"`
	RiakPort                 *uint                      `json:"riak_port"`
	WhitelistedOAuthUrls     []string                   `json:"whitelisted_oauth_urls"`
	OAuthClientSecret        string                     `json:"oauth_client_secret"`
	RoutingBlacklist         `json:"routing_blacklist"`
	SupportedDSMetrics       []string    `json:"supported_ds_metrics"`
	TLSConfig                *tls.Config `json:"tls_config"`

	// CRConfigUseRequestHost is whether to use the client request host header in the CRConfig. If false, uses the tm.url parameter.
	// This defaults to false. Traffic Ops used to always use the host header, setting this true will resume that legacy behavior.
	// See https://github.com/apache/trafficcontrol/issues/2224
	// Deprecated: will be removed in the next major version.
	CRConfigUseRequestHost bool `json:"crconfig_snapshot_use_client_request_host"`
	// CRConfigEmulateOldPath is whether to emulate the legacy CRConfig request path when generating a new CRConfig. This primarily exists in the event a tool relies on the legacy path '/tools/write_crconfig'.
	// Deprecated: will be removed in the next major version.
	CRConfigEmulateOldPath bool `json:"crconfig_emulate_old_path"`
}

// RoutingBlacklist contains the list of route IDs that will be handled by TO-Perl, a list of route IDs that are disabled,
// and whether or not to ignore unknown routes.
type RoutingBlacklist struct {
	IgnoreUnknownRoutes bool  `json:"ignore_unknown_routes"`
	PerlRoutes          []int `json:"perl_routes"`
	DisabledRoutes      []int `json:"disabled_routes"`
}

// ConfigTO contains information to identify Traffic Ops in a network sense.
type ConfigTO struct {
	BaseURL               *rfc.URL          `json:"base_url"`
	EmailFrom             *rfc.EmailAddress `json:"email_from"`
	NoAccountFoundMessage *string           `json:"no_account_found_msg"`
}

// ConfigPortal contains information that can direct users to a friendly UI
type ConfigPortal struct {
	BaseURL          rfc.URL          `json:"base_url"`
	DocsURL          rfc.URL          `json:"docs_url"`
	EmailFrom        rfc.EmailAddress `json:"email_from"`
	PasswdResetPath  string           `json:"pass_reset_path"`
	UserRegisterPath string           `json:"user_register_path"`
}

// ConfigSMTP contains configuration information for connecting to and authenticating with an SMTP
// server.
type ConfigSMTP struct {
	Address  string `json:"address"`
	Enabled  bool   `json:"enabled"`
	Password string `json:"password"`
	User     string `json:"user"`
}

// ConfigLetsEncrypt contains configuration information for integration with the Let's Encrypt certificate authority.
type ConfigLetsEncrypt struct {
	Email                     string `json:"user_email,omitempty"`
	SendExpEmail              bool   `json:"send_expiration_email"`
	ConvertSelfSigned         bool   `json:"convert_self_signed"`
	RenewDaysBeforeExpiration int    `json:"renew_days_before_expiration"`
	Environment               string `json:"environment"`
}

// ConfigDatabase reflects the structure of the database.conf file
type ConfigDatabase struct {
	Description string `json:"description"`
	DBName      string `json:"dbname"`
	Hostname    string `json:"hostname"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Port        string `json:"port"`
	Type        string `json:"type"`
	SSL         bool   `json:"ssl"`
}

type ConfigLDAP struct {
	AdminPass       string `json:"admin_pass"`
	SearchBase      string `json:"search_base"`
	AdminDN         string `json:"admin_dn"`
	Host            string `json:"host"`
	SearchQuery     string `json:"search_query"`
	Insecure        bool   `json:"insecure"`
	LDAPTimeoutSecs int    `json:"ldap_timeout_secs"`
}

type ConfigInflux struct {
	User        string `json:"user"`
	Password    string `json:"password"`
	DSDBName    string `json:"deliveryservice_stats_db_name"`
	CacheDBName string `json:"cache_stats_db_name"`
	Secure      *bool  `json:"secure"`
}

// NewFakeConfig returns a fake Config struct with just enough data to view Routes.
func NewFakeConfig() Config {
	c := Config{}
	c.URL, _ = url.Parse("http://example.com")
	c.Secrets = append(c.Secrets, "foo")
	c.BackendMaxConnections = make(map[string]int, 1)
	c.BackendMaxConnections["mojolicious"] = 42
	return c
}

const DefaultLDAPTimeoutSecs = 60
const DefaultDBQueryTimeoutSecs = 20

// ErrorLog - critical messages
func (c Config) ErrorLog() log.LogLocation {
	return log.LogLocation(c.LogLocationError)
}

// WarningLog - warning messages
func (c Config) WarningLog() log.LogLocation {
	return log.LogLocation(c.LogLocationWarning)
}

// InfoLog - information messages
func (c Config) InfoLog() log.LogLocation { return log.LogLocation(c.LogLocationInfo) }

// DebugLog - troubleshooting messages
func (c Config) DebugLog() log.LogLocation {
	return log.LogLocation(c.LogLocationDebug)
}

// EventLog - access.log high level transactions
func (c Config) EventLog() log.LogLocation {
	return log.LogLocation(c.LogLocationEvent)
}

const BlockStartup = true
const AllowStartup = false

func LoadCdnConfig(cdnConfPath string) (Config, error) {
	// load json from cdn.conf
	confBytes, err := ioutil.ReadFile(cdnConfPath)
	if err != nil {
		return Config{}, fmt.Errorf("reading CDN conf '%s': %v", cdnConfPath, err)
	}

	cfg := Config{}
	err = json.Unmarshal(confBytes, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshalling '%s': %v", cdnConfPath, err)
	}
	if cfg.SMTP == nil {
		cfg.SMTP = &ConfigSMTP{}
	}
	return cfg, nil
}

// LoadConfig - reads the config file into the Config struct

func LoadConfig(cdnConfPath string, dbConfPath string, riakConfPath string, appVersion string) (Config, []error, bool) {
	// load cdn.conf
	cfg, err := LoadCdnConfig(cdnConfPath)
	if err != nil {
		return Config{}, []error{fmt.Errorf("Loading cdn config from '%s': %v", cdnConfPath, err)}, BlockStartup
	}
	cfg.Version = appVersion

	// load json from database.conf
	dbConfBytes, err := ioutil.ReadFile(dbConfPath)
	if err != nil {
		return Config{}, []error{fmt.Errorf("reading db conf '%s': %v", dbConfPath, err)}, BlockStartup
	}
	err = json.Unmarshal(dbConfBytes, &cfg.DB)
	if err != nil {
		return Config{}, []error{fmt.Errorf("unmarshalling '%s': %v", dbConfPath, err)}, BlockStartup
	}
	cfg, err = ParseConfig(cfg)
	if err != nil {
		return Config{}, []error{fmt.Errorf("parsing config '%s': %v", cdnConfPath, err)}, BlockStartup
	}

	if riakConfPath != "" {
		cfg.RiakEnabled, cfg.RiakAuthOptions, err = riaksvc.GetRiakConfig(riakConfPath)
		if err != nil {
			return Config{}, []error{fmt.Errorf("parsing config '%s': %v", riakConfPath, err)}, BlockStartup
		}
	}
	// check for and load ldap.conf
	if cfg.LDAPConfPath != "" {
		cfg.LDAPEnabled, cfg.ConfigLDAP, err = GetLDAPConfig(cfg.LDAPConfPath)
		if err != nil {
			cfg.LDAPEnabled = false
			return cfg, []error{fmt.Errorf("parsing ldap config '%s': %v", cfg.LDAPConfPath, err)}, BlockStartup
		}
	} else { // ldap config location not specified in cdn.conf, check in directory with cdn.conf for backwards compatibility with perl.
		confDir := filepath.Dir(cdnConfPath)
		genericLDAPConfPath := filepath.Join(confDir, "ldap.conf")
		if _, err := os.Stat(genericLDAPConfPath); !os.IsNotExist(err) { // ldap.conf exists and we should error if it is not readable/parseable.
			cfg.LDAPEnabled, cfg.ConfigLDAP, err = GetLDAPConfig(genericLDAPConfPath)
			if err != nil { // no config or unparseable, do not enable LDAP
				cfg.LDAPEnabled = false
				return cfg, []error{err}, BlockStartup
			}
		} else {
			cfg.LDAPEnabled = false // no ldap.conf, disable and allow startup
		}
	}

	idbPath := cfg.InfluxDBConfPath
	if idbPath == "" {
		mojoMode := os.Getenv("MOJO_MODE")

		if cwd, err := os.Getwd(); mojoMode != "" && err != nil {
			idbPath = filepath.Join(cwd, "conf", mojoMode, "influxdb.conf")
		} else {
			idbPath = filepath.Join(filepath.Dir(cdnConfPath), "influxdb.conf")
		}
	}

	if _, err = os.Stat(idbPath); err != nil {
		if os.IsNotExist(err) {
			cfg.InfluxEnabled = false
		} else {
			return cfg, []error{err}, BlockStartup
		}
	} else if cfg.InfluxEnabled, cfg.ConfigInflux, err = GetInfluxConfig(idbPath); err != nil {
		return cfg, []error{err}, BlockStartup
	}

	return cfg, []error{}, AllowStartup
}

// GetCertPath - extracts path to cert .cert file
func (c Config) GetCertPath() string {
	v, ok := c.URL.Query()["cert"]
	if ok {
		return v[0]
	}
	return ""
}

// GetKeyPath - extracts path to cert .key file
func (c Config) GetKeyPath() string {
	v, ok := c.URL.Query()["key"]
	if ok {
		return v[0]
	}
	return ""
}

const (
	MojoliciousConcurrentConnectionsDefault = 12 // MojoliciousConcurrentConnectionsDefault
	DBMaxIdleConnectionsDefault             = 10 // if this is higher than MaxDBConnections it will be automatically adjusted below it by the db/sql library
	DBConnMaxLifetimeSecondsDefault         = 60
)

// ParseConfig validates required fields, and parses non-JSON types
func ParseConfig(cfg Config) (Config, error) {
	missings := ""
	if cfg.Port == "" {
		missings += "port, "
	}
	if len(cfg.Secrets) == 0 {
		missings += "secrets, "
	}
	if cfg.LogLocationError == "" {
		cfg.LogLocationError = log.LogLocationNull
	}
	if cfg.LogLocationWarning == "" {
		cfg.LogLocationWarning = log.LogLocationNull
	}
	if cfg.LogLocationInfo == "" {
		cfg.LogLocationInfo = log.LogLocationNull
	}
	if cfg.LogLocationDebug == "" {
		cfg.LogLocationDebug = log.LogLocationNull
	}
	if cfg.LogLocationEvent == "" {
		cfg.LogLocationEvent = log.LogLocationNull
	}
	if cfg.BackendMaxConnections == nil {
		cfg.BackendMaxConnections = make(map[string]int)
	}
	if cfg.BackendMaxConnections["mojolicious"] == 0 {
		cfg.BackendMaxConnections["mojolicious"] = MojoliciousConcurrentConnectionsDefault
	}
	if cfg.DBMaxIdleConnections == 0 {
		cfg.DBMaxIdleConnections = DBMaxIdleConnectionsDefault
	}
	if cfg.DBConnMaxLifetimeSeconds == 0 {
		cfg.DBConnMaxLifetimeSeconds = DBConnMaxLifetimeSecondsDefault
	}
	if cfg.DBQueryTimeoutSeconds == 0 {
		cfg.DBQueryTimeoutSeconds = DefaultDBQueryTimeoutSecs
	}

	invalidTOURLStr := ""
	var err error
	if len(cfg.Listen) < 1 {
		missings += `"listen", `
	} else {
		listen := cfg.Listen[0]
		if cfg.URL, err = url.Parse(listen); err != nil {
			invalidTOURLStr = fmt.Sprintf("invalid Traffic Ops URL '%s': %v", listen, err)
		}
		cfg.KeyPath = cfg.GetKeyPath()
		cfg.CertPath = cfg.GetCertPath()

		newURL := url.URL{Scheme: cfg.URL.Scheme, Host: cfg.URL.Host, Path: cfg.URL.Path}
		cfg.URL = &newURL
	}

	if cfg.ConfigTO == nil {
		missings += "to, "
	} else {
		if cfg.ConfigTO.BaseURL == nil || cfg.ConfigTO.BaseURL.String() == "" {
			missings += "to.base_url, "
		}
		if cfg.ConfigTO.EmailFrom == nil || cfg.ConfigTO.EmailFrom.String() == "<@>" {
			missings += "to.email_from, "
		}
		if cfg.ConfigTO.NoAccountFoundMessage == nil || *cfg.ConfigTO.NoAccountFoundMessage == "" {
			missings += "to.no_account_found_msg, "
		}
	}

	if len(missings) > 0 {
		missings = "missing fields: " + missings[:len(missings)-2] // strip final `, `
	}

	errStr := missings
	if errStr != "" && invalidTOURLStr != "" {
		errStr += "; "
	}
	errStr += invalidTOURLStr
	if errStr != "" {
		return Config{}, fmt.Errorf(errStr)
	}

	if err := ValidateRoutingBlacklist(cfg.RoutingBlacklist); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func ValidateRoutingBlacklist(blacklist RoutingBlacklist) error {
	seenPerlIDs := make(map[int]struct{}, len(blacklist.PerlRoutes))
	for _, id := range blacklist.PerlRoutes {
		if _, found := seenPerlIDs[id]; !found {
			seenPerlIDs[id] = struct{}{}
		} else {
			return fmt.Errorf("route ID %d is listed multiple times in perl_routes", id)
		}
	}
	seenDisabledIDs := make(map[int]struct{}, len(blacklist.DisabledRoutes))
	for _, id := range blacklist.DisabledRoutes {
		if _, foundInPerl := seenPerlIDs[id]; foundInPerl {
			return fmt.Errorf("route ID %d cannot be listed in both perl_routes and disabled_routes", id)
		} else if _, found := seenDisabledIDs[id]; !found {
			seenDisabledIDs[id] = struct{}{}
		} else {
			return fmt.Errorf("route ID %d is listed multiple times in disabled_routes", id)
		}
	}
	return nil
}

func GetLDAPConfig(LDAPConfPath string) (bool, *ConfigLDAP, error) {
	LDAPConfBytes, err := ioutil.ReadFile(LDAPConfPath)
	if err != nil {
		return false, nil, fmt.Errorf("reading LDAP conf '%v': %v", LDAPConfPath, err)
	}
	LDAPconf, err := getLDAPConf(string(LDAPConfBytes))
	if err != nil {
		return false, LDAPconf, fmt.Errorf("parsing LDAP conf '%v': %v", LDAPConfBytes, err)
	}
	if strings.TrimSpace(LDAPconf.AdminPass) == "" {
		return false, LDAPconf, fmt.Errorf("LDAP conf missing admin_pass field")
	}
	if strings.TrimSpace(LDAPconf.SearchBase) == "" {
		return false, LDAPconf, fmt.Errorf("LDAP conf missing search_base field")
	}
	if strings.TrimSpace(LDAPconf.AdminDN) == "" {
		return false, LDAPconf, fmt.Errorf("LDAP conf missing admin_dn field")
	}
	if strings.TrimSpace(LDAPconf.Host) == "" {
		return false, LDAPconf, fmt.Errorf("LDAP conf missing host field")
	}
	if strings.TrimSpace(LDAPconf.SearchQuery) == "" {
		return false, LDAPconf, fmt.Errorf("LDAP conf missing search_query field")
	}

	return true, LDAPconf, nil
}

func GetInfluxConfig(path string) (bool, *ConfigInflux, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return false, nil, fmt.Errorf("reading InfluxDB configuration from '%s': %v", path, err)
	}

	c := ConfigInflux{}
	if err := json.Unmarshal(raw, &c); err != nil {
		return false, nil, fmt.Errorf("parsing InfluxDB configuration from '%s': %v", path, err)
	}

	if c.User == "" {
		return false, &c, errors.New("InfluxDB configuration missing a username")
	}

	if c.Password == "" {
		return false, &c, errors.New("InfluxDB configuration missing a password")
	}

	if c.DSDBName == "" {
		log.Warnln("InfluxDB configuration does not specify a DS stats DB name - falling back on 'deliveryservice_stats'")
		c.DSDBName = "deliveryservice_stats"
	}

	if c.CacheDBName == "" {
		log.Warnln("InfluxDB configuration does not specify a Cache Stats DB name - falling back on 'cache_stats'")
		c.CacheDBName = "cache_stats"
	}

	if c.Secure == nil {
		log.Warnln("InfluxDB configuration does not specify 'secure', defaulting to 'false'")
		c.Secure = util.BoolPtr(false)
	}

	return true, &c, nil
}

func getLDAPConf(s string) (*ConfigLDAP, error) {
	ldapConf := ConfigLDAP{LDAPTimeoutSecs: DefaultLDAPTimeoutSecs} //if the field is not set in the config we use the default instead of 0
	err := json.Unmarshal([]byte(s), &ldapConf)
	return &ldapConf, err
}

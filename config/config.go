package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/DryPeng/clashT/adapter"
	"github.com/DryPeng/clashT/adapter/outbound"
	"github.com/DryPeng/clashT/adapter/outboundgroup"
	"github.com/DryPeng/clashT/adapter/provider"
	"github.com/DryPeng/clashT/component/auth"
	"github.com/DryPeng/clashT/component/fakeip"
	"github.com/DryPeng/clashT/component/trie"
	C "github.com/DryPeng/clashT/constant"
	providerTypes "github.com/DryPeng/clashT/constant/provider"
	"github.com/DryPeng/clashT/dns"
	"github.com/DryPeng/clashT/log"
	R "github.com/DryPeng/clashT/rule"
	T "github.com/DryPeng/clashT/tunnel"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// General config
type General struct {
	LegacyInbound
	Controller
	Authentication []string     `json:"authentication"`
	Mode           T.TunnelMode `json:"mode"`
	LogLevel       log.LogLevel `json:"log-level"`
	IPv6           bool         `json:"ipv6"`
	Interface      string       `json:"-"`
	RoutingMark    int          `json:"-"`
}

// Controller
type Controller struct {
	ExternalController string `json:"-"`
	ExternalUI         string `json:"-"`
	Secret             string `json:"-"`
}

type LegacyInbound struct {
	Port        int    `json:"port"`
	SocksPort   int    `json:"socks-port"`
	RedirPort   int    `json:"redir-port"`
	TProxyPort  int    `json:"tproxy-port"`
	MixedPort   int    `json:"mixed-port"`
	AllowLan    bool   `json:"allow-lan"`
	BindAddress string `json:"bind-address"`
}

// DNS config
type DNS struct {
	Enable            bool             `yaml:"enable"`
	IPv6              bool             `yaml:"ipv6"`
	NameServer        []dns.NameServer `yaml:"nameserver"`
	Fallback          []dns.NameServer `yaml:"fallback"`
	FallbackFilter    FallbackFilter   `yaml:"fallback-filter"`
	Listen            string           `yaml:"listen"`
	EnhancedMode      C.DNSMode        `yaml:"enhanced-mode"`
	DefaultNameserver []dns.NameServer `yaml:"default-nameserver"`
	FakeIPRange       *fakeip.Pool
	Hosts             *trie.DomainTrie
	NameServerPolicy  map[string]dns.NameServer
	SearchDomains     []string
}

// FallbackFilter config
type FallbackFilter struct {
	GeoIP     bool         `yaml:"geoip"`
	GeoIPCode string       `yaml:"geoip-code"`
	IPCIDR    []*net.IPNet `yaml:"ipcidr"`
	Domain    []string     `yaml:"domain"`
}

// Profile config
type Profile struct {
	StoreSelected bool `yaml:"store-selected"`
	StoreFakeIP   bool `yaml:"store-fake-ip"`
}

// Experimental config
type Experimental struct {
	UDPFallbackMatch bool `yaml:"udp-fallback-match"`
}

// Config is clash config manager
type Config struct {
	General      *General
	DNS          *DNS
	Experimental *Experimental
	Hosts        *trie.DomainTrie
	Profile      *Profile
	Inbounds     []C.Inbound
	Rules        []C.Rule
	Users        []auth.AuthUser
	Proxies      map[string]C.Proxy
	Providers    map[string]providerTypes.ProxyProvider
	Tunnels      []Tunnel
}

type RawDNS struct {
	Enable            bool              `yaml:"enable"`
	IPv6              *bool             `yaml:"ipv6"`
	UseHosts          bool              `yaml:"use-hosts"`
	NameServer        []string          `yaml:"nameserver"`
	Fallback          []string          `yaml:"fallback"`
	FallbackFilter    RawFallbackFilter `yaml:"fallback-filter"`
	Listen            string            `yaml:"listen"`
	EnhancedMode      C.DNSMode         `yaml:"enhanced-mode"`
	FakeIPRange       string            `yaml:"fake-ip-range"`
	FakeIPFilter      []string          `yaml:"fake-ip-filter"`
	DefaultNameserver []string          `yaml:"default-nameserver"`
	NameServerPolicy  map[string]string `yaml:"nameserver-policy"`
	SearchDomains     []string          `yaml:"search-domains"`
}

type RawFallbackFilter struct {
	GeoIP     bool     `yaml:"geoip"`
	GeoIPCode string   `yaml:"geoip-code"`
	IPCIDR    []string `yaml:"ipcidr"`
	Domain    []string `yaml:"domain"`
}

type tunnel struct {
	Network []string `yaml:"network"`
	Address string   `yaml:"address"`
	Target  string   `yaml:"target"`
	Proxy   string   `yaml:"proxy"`
}

type Tunnel tunnel

// UnmarshalYAML implements yaml.Unmarshaler
func (t *Tunnel) UnmarshalYAML(unmarshal func(any) error) error {
	var tp string
	if err := unmarshal(&tp); err != nil {
		var inner tunnel
		if err := unmarshal(&inner); err != nil {
			return err
		}

		*t = Tunnel(inner)
		return nil
	}

	// parse udp/tcp,address,target,proxy
	parts := lo.Map(strings.Split(tp, ","), func(s string, _ int) string {
		return strings.TrimSpace(s)
	})
	if len(parts) != 4 {
		return fmt.Errorf("invalid tunnel config %s", tp)
	}
	network := strings.Split(parts[0], "/")

	// validate network
	for _, n := range network {
		switch n {
		case "tcp", "udp":
		default:
			return fmt.Errorf("invalid tunnel network %s", n)
		}
	}

	// validate address and target
	address := parts[1]
	target := parts[2]
	for _, addr := range []string{address, target} {
		if _, _, err := net.SplitHostPort(addr); err != nil {
			return fmt.Errorf("invalid tunnel target or address %s", addr)
		}
	}

	*t = Tunnel(tunnel{
		Network: network,
		Address: address,
		Target:  target,
		Proxy:   parts[3],
	})
	return nil
}

type RawConfig struct {
	Port               int          `yaml:"port"`
	SocksPort          int          `yaml:"socks-port"`
	RedirPort          int          `yaml:"redir-port"`
	TProxyPort         int          `yaml:"tproxy-port"`
	MixedPort          int          `yaml:"mixed-port"`
	Authentication     []string     `yaml:"authentication"`
	AllowLan           bool         `yaml:"allow-lan"`
	BindAddress        string       `yaml:"bind-address"`
	Mode               T.TunnelMode `yaml:"mode"`
	LogLevel           log.LogLevel `yaml:"log-level"`
	IPv6               bool         `yaml:"ipv6"`
	ExternalController string       `yaml:"external-controller"`
	ExternalUI         string       `yaml:"external-ui"`
	Secret             string       `yaml:"secret"`
	Interface          string       `yaml:"interface-name"`
	RoutingMark        int          `yaml:"routing-mark"`
	Tunnels            []Tunnel     `yaml:"tunnels"`

	ProxyProvider map[string]map[string]any `yaml:"proxy-providers"`
	Hosts         map[string]string         `yaml:"hosts"`
	Inbounds      []C.Inbound               `yaml:"inbounds"`
	DNS           RawDNS                    `yaml:"dns"`
	Experimental  Experimental              `yaml:"experimental"`
	Profile       Profile                   `yaml:"profile"`
	Proxy         []map[string]any          `yaml:"proxies"`
	ProxyGroup    []map[string]any          `yaml:"proxy-groups"`
	Rule          []string                  `yaml:"rules"`
}

// Parse config
func Parse(buf []byte) (*Config, error) {
	rawCfg, err := UnmarshalRawConfig(buf)
	if err != nil {
		return nil, err
	}

	return ParseRawConfig(rawCfg)
}

func UnmarshalRawConfig(buf []byte) (*RawConfig, error) {
	// config with default value
	rawCfg := &RawConfig{
		AllowLan:       false,
		BindAddress:    "*",
		Mode:           T.Rule,
		Authentication: []string{},
		LogLevel:       log.INFO,
		Hosts:          map[string]string{},
		Rule:           []string{},
		Proxy:          []map[string]any{},
		ProxyGroup:     []map[string]any{},
		DNS: RawDNS{
			Enable:      false,
			UseHosts:    true,
			FakeIPRange: "198.18.0.1/16",
			FallbackFilter: RawFallbackFilter{
				GeoIP:     true,
				GeoIPCode: "CN",
				IPCIDR:    []string{},
			},
			DefaultNameserver: []string{
				"114.114.114.114",
				"8.8.8.8",
			},
		},
		Profile: Profile{
			StoreSelected: true,
		},
	}

	if err := yaml.Unmarshal(buf, rawCfg); err != nil {
		return nil, err
	}

	return rawCfg, nil
}

func ParseRawConfig(rawCfg *RawConfig) (*Config, error) {
	config := &Config{}

	config.Experimental = &rawCfg.Experimental
	config.Profile = &rawCfg.Profile

	general, err := parseGeneral(rawCfg)
	if err != nil {
		return nil, err
	}
	config.General = general

	config.Inbounds = rawCfg.Inbounds

	proxies, providers, err := parseProxies(rawCfg)
	if err != nil {
		return nil, err
	}
	config.Proxies = proxies
	config.Providers = providers

	rules, err := parseRules(rawCfg, proxies)
	if err != nil {
		return nil, err
	}
	config.Rules = rules

	hosts, err := parseHosts(rawCfg)
	if err != nil {
		return nil, err
	}
	config.Hosts = hosts

	dnsCfg, err := parseDNS(rawCfg, hosts)
	if err != nil {
		return nil, err
	}
	config.DNS = dnsCfg

	config.Users = parseAuthentication(rawCfg.Authentication)

	config.Tunnels = rawCfg.Tunnels
	// verify tunnels
	for _, t := range config.Tunnels {
		if _, ok := config.Proxies[t.Proxy]; !ok {
			return nil, fmt.Errorf("tunnel proxy %s not found", t.Proxy)
		}
	}

	return config, nil
}

func parseGeneral(cfg *RawConfig) (*General, error) {
	externalUI := cfg.ExternalUI

	// checkout externalUI exist
	if externalUI != "" {
		externalUI = C.Path.Resolve(externalUI)

		if _, err := os.Stat(externalUI); os.IsNotExist(err) {
			return nil, fmt.Errorf("external-ui: %s not exist", externalUI)
		}
	}

	return &General{
		LegacyInbound: LegacyInbound{
			Port:        cfg.Port,
			SocksPort:   cfg.SocksPort,
			RedirPort:   cfg.RedirPort,
			TProxyPort:  cfg.TProxyPort,
			MixedPort:   cfg.MixedPort,
			AllowLan:    cfg.AllowLan,
			BindAddress: cfg.BindAddress,
		},
		Controller: Controller{
			ExternalController: cfg.ExternalController,
			ExternalUI:         cfg.ExternalUI,
			Secret:             cfg.Secret,
		},
		Mode:        cfg.Mode,
		LogLevel:    cfg.LogLevel,
		IPv6:        cfg.IPv6,
		Interface:   cfg.Interface,
		RoutingMark: cfg.RoutingMark,
	}, nil
}

func parseProxies(cfg *RawConfig) (proxies map[string]C.Proxy, providersMap map[string]providerTypes.ProxyProvider, err error) {
	proxies = make(map[string]C.Proxy)
	providersMap = make(map[string]providerTypes.ProxyProvider)
	proxyList := []string{}
	proxiesConfig := cfg.Proxy
	groupsConfig := cfg.ProxyGroup
	providersConfig := cfg.ProxyProvider

	proxies["DIRECT"] = adapter.NewProxy(outbound.NewDirect())
	proxies["REJECT"] = adapter.NewProxy(outbound.NewReject())
	proxyList = append(proxyList, "DIRECT", "REJECT")

	// parse proxy
	for idx, mapping := range proxiesConfig {
		proxy, err := adapter.ParseProxy(mapping)
		if err != nil {
			return nil, nil, fmt.Errorf("proxy %d: %w", idx, err)
		}

		if _, exist := proxies[proxy.Name()]; exist {
			return nil, nil, fmt.Errorf("proxy %s is the duplicate name", proxy.Name())
		}
		proxies[proxy.Name()] = proxy
		proxyList = append(proxyList, proxy.Name())
	}

	// keep the original order of ProxyGroups in config file
	for idx, mapping := range groupsConfig {
		groupName, existName := mapping["name"].(string)
		if !existName {
			return nil, nil, fmt.Errorf("proxy group %d: missing name", idx)
		}
		proxyList = append(proxyList, groupName)
	}

	// check if any loop exists and sort the ProxyGroups
	if err := proxyGroupsDagSort(groupsConfig); err != nil {
		return nil, nil, err
	}

	// parse and initial providers
	for name, mapping := range providersConfig {
		if name == provider.ReservedName {
			return nil, nil, fmt.Errorf("can not defined a provider called `%s`", provider.ReservedName)
		}

		pd, err := provider.ParseProxyProvider(name, mapping)
		if err != nil {
			return nil, nil, fmt.Errorf("parse proxy provider %s error: %w", name, err)
		}

		providersMap[name] = pd
	}

	for _, provider := range providersMap {
		log.Infoln("Start initial provider %s", provider.Name())
		if err := provider.Initial(); err != nil {
			return nil, nil, fmt.Errorf("initial proxy provider %s error: %w", provider.Name(), err)
		}
	}

	// parse proxy group
	for idx, mapping := range groupsConfig {
		group, err := outboundgroup.ParseProxyGroup(mapping, proxies, providersMap)
		if err != nil {
			return nil, nil, fmt.Errorf("proxy group[%d]: %w", idx, err)
		}

		groupName := group.Name()
		if _, exist := proxies[groupName]; exist {
			return nil, nil, fmt.Errorf("proxy group %s: the duplicate name", groupName)
		}

		proxies[groupName] = adapter.NewProxy(group)
	}

	// initial compatible provider
	for _, pd := range providersMap {
		if pd.VehicleType() != providerTypes.Compatible {
			continue
		}

		log.Infoln("Start initial compatible provider %s", pd.Name())
		if err := pd.Initial(); err != nil {
			return nil, nil, err
		}
	}

	ps := []C.Proxy{}
	for _, v := range proxyList {
		ps = append(ps, proxies[v])
	}
	hc := provider.NewHealthCheck(ps, "", 0, true)
	pd, _ := provider.NewCompatibleProvider(provider.ReservedName, ps, hc)
	providersMap[provider.ReservedName] = pd

	global := outboundgroup.NewSelector(
		&outboundgroup.GroupCommonOption{
			Name: "GLOBAL",
		},
		[]providerTypes.ProxyProvider{pd},
	)
	proxies["GLOBAL"] = adapter.NewProxy(global)
	return proxies, providersMap, nil
}

func parseRules(cfg *RawConfig, proxies map[string]C.Proxy) ([]C.Rule, error) {
	rules := []C.Rule{}
	rulesConfig := cfg.Rule

	// parse rules
	for idx, line := range rulesConfig {
		rule := trimArr(strings.Split(line, ","))
		var (
			payload string
			target  string
			params  = []string{}
		)

		switch l := len(rule); {
		case l == 2:
			target = rule[1]
		case l == 3:
			payload = rule[1]
			target = rule[2]
		case l >= 4:
			payload = rule[1]
			target = rule[2]
			params = rule[3:]
		default:
			return nil, fmt.Errorf("rules[%d] [%s] error: format invalid", idx, line)
		}

		if _, ok := proxies[target]; !ok {
			return nil, fmt.Errorf("rules[%d] [%s] error: proxy [%s] not found", idx, line, target)
		}

		rule = trimArr(rule)
		params = trimArr(params)

		parsed, parseErr := R.ParseRule(rule[0], payload, target, params)
		if parseErr != nil {
			return nil, fmt.Errorf("rules[%d] [%s] error: %s", idx, line, parseErr.Error())
		}

		rules = append(rules, parsed)
	}

	return rules, nil
}

func parseHosts(cfg *RawConfig) (*trie.DomainTrie, error) {
	tree := trie.New()

	// add default hosts
	if err := tree.Insert("localhost", net.IP{127, 0, 0, 1}); err != nil {
		log.Errorln("insert localhost to host error: %s", err.Error())
	}

	if len(cfg.Hosts) != 0 {
		for domain, ipStr := range cfg.Hosts {
			ip := net.ParseIP(ipStr)
			if ip == nil {
				return nil, fmt.Errorf("%s is not a valid IP", ipStr)
			}
			tree.Insert(domain, ip)
		}
	}

	return tree, nil
}

func hostWithDefaultPort(host string, defPort string) (string, error) {
	if !strings.Contains(host, ":") {
		host += ":"
	}

	hostname, port, err := net.SplitHostPort(host)
	if err != nil {
		return "", err
	}

	if port == "" {
		port = defPort
	}

	return net.JoinHostPort(hostname, port), nil
}

func parseNameServer(servers []string) ([]dns.NameServer, error) {
	nameservers := []dns.NameServer{}

	for idx, server := range servers {
		// parse without scheme .e.g 8.8.8.8:53
		if !strings.Contains(server, "://") {
			server = "udp://" + server
		}
		u, err := url.Parse(server)
		if err != nil {
			return nil, fmt.Errorf("DNS NameServer[%d] format error: %s", idx, err.Error())
		}

		// parse with specific interface
		// .e.g 10.0.0.1#en0
		interfaceName := u.Fragment

		var addr, dnsNetType string
		switch u.Scheme {
		case "udp":
			addr, err = hostWithDefaultPort(u.Host, "53")
			dnsNetType = "" // UDP
		case "tcp":
			addr, err = hostWithDefaultPort(u.Host, "53")
			dnsNetType = "tcp" // TCP
		case "tls":
			addr, err = hostWithDefaultPort(u.Host, "853")
			dnsNetType = "tcp-tls" // DNS over TLS
		case "https":
			clearURL := url.URL{Scheme: "https", Host: u.Host, Path: u.Path, User: u.User}
			addr = clearURL.String()
			dnsNetType = "https" // DNS over HTTPS
		case "dhcp":
			addr = u.Host
			dnsNetType = "dhcp" // UDP from DHCP
		default:
			return nil, fmt.Errorf("DNS NameServer[%d] unsupport scheme: %s", idx, u.Scheme)
		}

		if err != nil {
			return nil, fmt.Errorf("DNS NameServer[%d] format error: %s", idx, err.Error())
		}

		nameservers = append(
			nameservers,
			dns.NameServer{
				Net:       dnsNetType,
				Addr:      addr,
				Interface: interfaceName,
			},
		)
	}
	return nameservers, nil
}

func parseNameServerPolicy(nsPolicy map[string]string) (map[string]dns.NameServer, error) {
	policy := map[string]dns.NameServer{}

	for domain, server := range nsPolicy {
		nameservers, err := parseNameServer([]string{server})
		if err != nil {
			return nil, err
		}
		if _, valid := trie.ValidAndSplitDomain(domain); !valid {
			return nil, fmt.Errorf("DNS ResoverRule invalid domain: %s", domain)
		}
		policy[domain] = nameservers[0]
	}

	return policy, nil
}

func parseFallbackIPCIDR(ips []string) ([]*net.IPNet, error) {
	ipNets := []*net.IPNet{}

	for idx, ip := range ips {
		_, ipnet, err := net.ParseCIDR(ip)
		if err != nil {
			return nil, fmt.Errorf("DNS FallbackIP[%d] format error: %s", idx, err.Error())
		}
		ipNets = append(ipNets, ipnet)
	}

	return ipNets, nil
}

func parseDNS(rawCfg *RawConfig, hosts *trie.DomainTrie) (*DNS, error) {
	cfg := rawCfg.DNS
	if cfg.Enable && len(cfg.NameServer) == 0 {
		return nil, fmt.Errorf("if DNS configuration is turned on, NameServer cannot be empty")
	}

	dnsCfg := &DNS{
		Enable:       cfg.Enable,
		Listen:       cfg.Listen,
		IPv6:         lo.FromPtrOr(cfg.IPv6, rawCfg.IPv6),
		EnhancedMode: cfg.EnhancedMode,
		FallbackFilter: FallbackFilter{
			IPCIDR: []*net.IPNet{},
		},
	}
	var err error
	if dnsCfg.NameServer, err = parseNameServer(cfg.NameServer); err != nil {
		return nil, err
	}

	if dnsCfg.Fallback, err = parseNameServer(cfg.Fallback); err != nil {
		return nil, err
	}

	if dnsCfg.NameServerPolicy, err = parseNameServerPolicy(cfg.NameServerPolicy); err != nil {
		return nil, err
	}

	if len(cfg.DefaultNameserver) == 0 {
		return nil, errors.New("default nameserver should have at least one nameserver")
	}
	if dnsCfg.DefaultNameserver, err = parseNameServer(cfg.DefaultNameserver); err != nil {
		return nil, err
	}
	// check default nameserver is pure ip addr
	for _, ns := range dnsCfg.DefaultNameserver {
		host, _, err := net.SplitHostPort(ns.Addr)
		if err != nil || net.ParseIP(host) == nil {
			return nil, errors.New("default nameserver should be pure IP")
		}
	}

	if cfg.EnhancedMode == C.DNSFakeIP {
		_, ipnet, err := net.ParseCIDR(cfg.FakeIPRange)
		if err != nil {
			return nil, err
		}

		var host *trie.DomainTrie
		// fake ip skip host filter
		if len(cfg.FakeIPFilter) != 0 {
			host = trie.New()
			for _, domain := range cfg.FakeIPFilter {
				host.Insert(domain, true)
			}
		}

		pool, err := fakeip.New(fakeip.Options{
			IPNet:       ipnet,
			Size:        1000,
			Host:        host,
			Persistence: rawCfg.Profile.StoreFakeIP,
		})
		if err != nil {
			return nil, err
		}

		dnsCfg.FakeIPRange = pool
	}

	dnsCfg.FallbackFilter.GeoIP = cfg.FallbackFilter.GeoIP
	dnsCfg.FallbackFilter.GeoIPCode = cfg.FallbackFilter.GeoIPCode
	if fallbackip, err := parseFallbackIPCIDR(cfg.FallbackFilter.IPCIDR); err == nil {
		dnsCfg.FallbackFilter.IPCIDR = fallbackip
	}
	dnsCfg.FallbackFilter.Domain = cfg.FallbackFilter.Domain

	if cfg.UseHosts {
		dnsCfg.Hosts = hosts
	}

	if len(cfg.SearchDomains) != 0 {
		for _, domain := range cfg.SearchDomains {
			if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
				return nil, errors.New("search domains should not start or end with '.'")
			}
			if strings.Contains(domain, ":") {
				return nil, errors.New("search domains are for ipv4 only and should not contain ports")
			}
		}
		dnsCfg.SearchDomains = cfg.SearchDomains
	}

	return dnsCfg, nil
}

func parseAuthentication(rawRecords []string) []auth.AuthUser {
	users := []auth.AuthUser{}
	for _, line := range rawRecords {
		if user, pass, found := strings.Cut(line, ":"); found {
			users = append(users, auth.AuthUser{User: user, Pass: pass})
		}
	}
	return users
}

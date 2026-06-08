package router

import "testing"

func TestRoute(t *testing.T) {
	cases := []struct {
		upstream string
		iface    string
		want     string
	}{
		{"alipay", "callback", "/alipay/callback"},
		{"wechat", "callback", "/wechat/callback"},
		{"alipay", "create", "/alipay/create"},
	}
	for _, c := range cases {
		if got := Route(c.upstream, c.iface); got != c.want {
			t.Errorf("Route(%q, %q) = %q, want %q", c.upstream, c.iface, got, c.want)
		}
	}
}

func TestNotify(t *testing.T) {
	cases := []struct {
		name     string
		endpoint string
		upstream string
		want     string
	}{
		{"no trailing slash", "https://gh.ink", "alipay", "https://gh.ink/alipay/callback"},
		{"trailing slash", "https://gh.ink/", "alipay", "https://gh.ink/alipay/callback"},
		{"path prefix", "https://gh.ink/pay", "wechat", "https://gh.ink/pay/wechat/callback"},
		{"path prefix trailing slash", "https://gh.ink/pay/", "wechat", "https://gh.ink/pay/wechat/callback"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Notify(c.endpoint, c.upstream); got != c.want {
				t.Errorf("Notify(%q, %q) = %q, want %q", c.endpoint, c.upstream, got, c.want)
			}
		})
	}
}

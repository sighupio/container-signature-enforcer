package config

import (
	"testing"

	"github.com/ghodss/yaml"
)

func TestParsingGlobalConfig(t *testing.T) {
	var tests = []struct {
		testcase string
		config   string
	}{
		{testcase: "simple", config: `repositories:
      - name: 'localhost.*'
        namespace: "webhook"
        priority: 10
        trust:
          enabled: true
          trustServer: 'https://notary-server.notary.svc.cluster.local:4443'
          signers:
          - role: 'targets/jenkins'
            publicKey: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURDakNDQWZJQ0NRQzdDNHZqQnN0bW96QU5CZ2txaGtpRzl3MEJBUXNGQURCSE1ROHdEUVlEVlFRS0RBWlQKU1VkSVZWQXhGVEFUQmdOVkJBc01ESEJ5YjJSMVkzUXRkR1ZoYlRFZE1Cc0dBMVVFQXd3VWIzQmhMVzV2ZEdGeQplUzFqYjI1dVpXTjBiM0l3SGhjTk1qQXdOekUwTURrd01ESTBXaGNOTWpFd056RTBNRGt3TURJMFdqQkhNUTh3CkRRWURWUVFLREFaVFNVZElWVkF4RlRBVEJnTlZCQXNNREhCeWIyUjFZM1F0ZEdWaGJURWRNQnNHQTFVRUF3d1UKYjNCaExXNXZkR0Z5ZVMxamIyNXVaV04wYjNJd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFSwpBb0lCQVFDNWZsTk9rVXNTeFkwTW5uRjExMjJnYktFWk9mQ2R6cWlPUXVpeFBMVmEwc0huL2FEamVDQUk1K1VmCi9QRVdOL1JiZzJVdCtjZHNEUDFVV3RCMVJ6M1JvUDBZTnNtS3UyNHpvLzdTS2V4dXlFcFljalhQM1FtN3hKUEgKaXJRV2swcGNHYnIzMjJPWlRDK0t4Y0E1VVh5NGFpbElONUVIbGovcU9xM1Nzd3R5bG9GbGxBbkViRmRHcDRxWQpTWkFNczhoa0FLZU93REJjUEcxQW1WR0pOdGlrOWNscFlqSEdyUXBTOVd3OVgvUTVPNU8vK0gxSmF4ZnNCMElNCmZOdGxmTlhkTGs0STFmeGtjcTAvWlZoQ1Vmd2ZJT3NJNVdBaklVK3ZLQmx0QytENE82bUJtUmxJRDlUclZtcTEKb21KM0tRNjNZUFpHMVYxNTRnM2NhTU9KakpVNUFnTUJBQUV3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUNJcwp6VWsrWFk0NzJHSHQxWjl5VWdzOGkyN3pHQ0hUTUp3b2V3Y0RpL2FwQ0pNcFZHT3gvMEVsR1cxY2xySVZSbjhOCkN5a3NPaFlXbnBqVUVVRGYyZHY1SkRHSGpBK0ExTFNUUEVRYXhCTXEvOEhkekN1WFdsN2xrTDdXWW9KWWQvOWkKTFRIOUpBVkNtckh6VklLeWd3d1ZSVHIwZVhRbGJ1ZnpFd01TU0FUWFJmMTFwekorazZyVE1icmNIT2pJb3FreQpIVVZCOHJsb3RUMUgxdFBUOVVzcUhoR0N3eUdad2MwSkNSSXZwemJsdUc4ZUFCL1gxWXdmblQzbG9iQzczT2VXCnd0ZkdTd25EN3IwS2E0YWdoQVMraWtRNDdtdklIWFVOTzU3WUt3dXJkUllrMjZxQzZqRTZRM3haU1J3MC92SkIKbm1SaUl6SmJUWUQvYVluT3N6WT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			c := NewGlobalConfig()
			err := yaml.Unmarshal([]byte(tt.config), c.GetConfig())
			if err != nil {
				t.Errorf("got error while parsing config: %s", err.Error())
				return
			}
		})
	}

}

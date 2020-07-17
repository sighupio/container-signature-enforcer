package config

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/sighupio/opa-notary-connector/reference"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestParsingConfig(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		testcase               string
		config                 string
		image                  string
		expectedRepositoryName string
		expectedReposLen       int
	}{
		{
			testcase:               "single repository",
			image:                  "localhost/test:prod",
			expectedRepositoryName: "localhost.*",
			expectedReposLen:       1,
			config: `repositories:
      - name: 'localhost.*'
        namespace: "webhook"
        priority: 10
        trust:
          enabled: true
          trustServer: "https://notary-server.notary.svc.cluster.local:4443"
          signers:
          - role: "targets/jenkins"
            publicKey: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURDakNDQWZJQ0NRQzdDNHZqQnN0bW96QU5CZ2txaGtpRzl3MEJBUXNGQURCSE1ROHdEUVlEVlFRS0RBWlQKU1VkSVZWQXhGVEFUQmdOVkJBc01ESEJ5YjJSMVkzUXRkR1ZoYlRFZE1Cc0dBMVVFQXd3VWIzQmhMVzV2ZEdGeQplUzFqYjI1dVpXTjBiM0l3SGhjTk1qQXdOekUwTURrd01ESTBXaGNOTWpFd056RTBNRGt3TURJMFdqQkhNUTh3CkRRWURWUVFLREFaVFNVZElWVkF4RlRBVEJnTlZCQXNNREhCeWIyUjFZM1F0ZEdWaGJURWRNQnNHQTFVRUF3d1UKYjNCaExXNXZkR0Z5ZVMxamIyNXVaV04wYjNJd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFSwpBb0lCQVFDNWZsTk9rVXNTeFkwTW5uRjExMjJnYktFWk9mQ2R6cWlPUXVpeFBMVmEwc0huL2FEamVDQUk1K1VmCi9QRVdOL1JiZzJVdCtjZHNEUDFVV3RCMVJ6M1JvUDBZTnNtS3UyNHpvLzdTS2V4dXlFcFljalhQM1FtN3hKUEgKaXJRV2swcGNHYnIzMjJPWlRDK0t4Y0E1VVh5NGFpbElONUVIbGovcU9xM1Nzd3R5bG9GbGxBbkViRmRHcDRxWQpTWkFNczhoa0FLZU93REJjUEcxQW1WR0pOdGlrOWNscFlqSEdyUXBTOVd3OVgvUTVPNU8vK0gxSmF4ZnNCMElNCmZOdGxmTlhkTGs0STFmeGtjcTAvWlZoQ1Vmd2ZJT3NJNVdBaklVK3ZLQmx0QytENE82bUJtUmxJRDlUclZtcTEKb21KM0tRNjNZUFpHMVYxNTRnM2NhTU9KakpVNUFnTUJBQUV3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUNJcwp6VWsrWFk0NzJHSHQxWjl5VWdzOGkyN3pHQ0hUTUp3b2V3Y0RpL2FwQ0pNcFZHT3gvMEVsR1cxY2xySVZSbjhOCkN5a3NPaFlXbnBqVUVVRGYyZHY1SkRHSGpBK0ExTFNUUEVRYXhCTXEvOEhkekN1WFdsN2xrTDdXWW9KWWQvOWkKTFRIOUpBVkNtckh6VklLeWd3d1ZSVHIwZVhRbGJ1ZnpFd01TU0FUWFJmMTFwekorazZyVE1icmNIT2pJb3FreQpIVVZCOHJsb3RUMUgxdFBUOVVzcUhoR0N3eUdad2MwSkNSSXZwemJsdUc4ZUFCL1gxWXdmblQzbG9iQzczT2VXCnd0ZkdTd25EN3IwS2E0YWdoQVMraWtRNDdtdklIWFVOTzU3WUt3dXJkUllrMjZxQzZqRTZRM3haU1J3MC92SkIKbm1SaUl6SmJUWUQvYVluT3N6WT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="`,
		},
		{
			testcase:               "two repositories matching",
			image:                  "localhost/test:prod",
			expectedRepositoryName: "localhost/test.*",
			expectedReposLen:       2,
			config: `repositories:
      - name: 'localhost.*'
        namespace: "webhook"
        priority: 10
        trust:
          enabled: true
          trustServer: "https://notary-server.notary.svc.cluster.local:4443"
          signers:
          - role: "targets/jenkins"
            publicKey: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURDakNDQWZJQ0NRQzdDNHZqQnN0bW96QU5CZ2txaGtpRzl3MEJBUXNGQURCSE1ROHdEUVlEVlFRS0RBWlQKU1VkSVZWQXhGVEFUQmdOVkJBc01ESEJ5YjJSMVkzUXRkR1ZoYlRFZE1Cc0dBMVVFQXd3VWIzQmhMVzV2ZEdGeQplUzFqYjI1dVpXTjBiM0l3SGhjTk1qQXdOekUwTURrd01ESTBXaGNOTWpFd056RTBNRGt3TURJMFdqQkhNUTh3CkRRWURWUVFLREFaVFNVZElWVkF4RlRBVEJnTlZCQXNNREhCeWIyUjFZM1F0ZEdWaGJURWRNQnNHQTFVRUF3d1UKYjNCaExXNXZkR0Z5ZVMxamIyNXVaV04wYjNJd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFSwpBb0lCQVFDNWZsTk9rVXNTeFkwTW5uRjExMjJnYktFWk9mQ2R6cWlPUXVpeFBMVmEwc0huL2FEamVDQUk1K1VmCi9QRVdOL1JiZzJVdCtjZHNEUDFVV3RCMVJ6M1JvUDBZTnNtS3UyNHpvLzdTS2V4dXlFcFljalhQM1FtN3hKUEgKaXJRV2swcGNHYnIzMjJPWlRDK0t4Y0E1VVh5NGFpbElONUVIbGovcU9xM1Nzd3R5bG9GbGxBbkViRmRHcDRxWQpTWkFNczhoa0FLZU93REJjUEcxQW1WR0pOdGlrOWNscFlqSEdyUXBTOVd3OVgvUTVPNU8vK0gxSmF4ZnNCMElNCmZOdGxmTlhkTGs0STFmeGtjcTAvWlZoQ1Vmd2ZJT3NJNVdBaklVK3ZLQmx0QytENE82bUJtUmxJRDlUclZtcTEKb21KM0tRNjNZUFpHMVYxNTRnM2NhTU9KakpVNUFnTUJBQUV3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUNJcwp6VWsrWFk0NzJHSHQxWjl5VWdzOGkyN3pHQ0hUTUp3b2V3Y0RpL2FwQ0pNcFZHT3gvMEVsR1cxY2xySVZSbjhOCkN5a3NPaFlXbnBqVUVVRGYyZHY1SkRHSGpBK0ExTFNUUEVRYXhCTXEvOEhkekN1WFdsN2xrTDdXWW9KWWQvOWkKTFRIOUpBVkNtckh6VklLeWd3d1ZSVHIwZVhRbGJ1ZnpFd01TU0FUWFJmMTFwekorazZyVE1icmNIT2pJb3FreQpIVVZCOHJsb3RUMUgxdFBUOVVzcUhoR0N3eUdad2MwSkNSSXZwemJsdUc4ZUFCL1gxWXdmblQzbG9iQzczT2VXCnd0ZkdTd25EN3IwS2E0YWdoQVMraWtRNDdtdklIWFVOTzU3WUt3dXJkUllrMjZxQzZqRTZRM3haU1J3MC92SkIKbm1SaUl6SmJUWUQvYVluT3N6WT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
      - name: 'localhost/test.*'
        namespace: "webhook"
        priority: 11
        trust:
          enabled: true
          trustServer: "https://notary-server.notary.svc.cluster.local:4443"
          signers:
          - role: "targets/jenkins"
            publicKey: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURDakNDQWZJQ0NRQzdDNHZqQnN0bW96QU5CZ2txaGtpRzl3MEJBUXNGQURCSE1ROHdEUVlEVlFRS0RBWlQKU1VkSVZWQXhGVEFUQmdOVkJBc01ESEJ5YjJSMVkzUXRkR1ZoYlRFZE1Cc0dBMVVFQXd3VWIzQmhMVzV2ZEdGeQplUzFqYjI1dVpXTjBiM0l3SGhjTk1qQXdOekUwTURrd01ESTBXaGNOTWpFd056RTBNRGt3TURJMFdqQkhNUTh3CkRRWURWUVFLREFaVFNVZElWVkF4RlRBVEJnTlZCQXNNREhCeWIyUjFZM1F0ZEdWaGJURWRNQnNHQTFVRUF3d1UKYjNCaExXNXZkR0Z5ZVMxamIyNXVaV04wYjNJd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFSwpBb0lCQVFDNWZsTk9rVXNTeFkwTW5uRjExMjJnYktFWk9mQ2R6cWlPUXVpeFBMVmEwc0huL2FEamVDQUk1K1VmCi9QRVdOL1JiZzJVdCtjZHNEUDFVV3RCMVJ6M1JvUDBZTnNtS3UyNHpvLzdTS2V4dXlFcFljalhQM1FtN3hKUEgKaXJRV2swcGNHYnIzMjJPWlRDK0t4Y0E1VVh5NGFpbElONUVIbGovcU9xM1Nzd3R5bG9GbGxBbkViRmRHcDRxWQpTWkFNczhoa0FLZU93REJjUEcxQW1WR0pOdGlrOWNscFlqSEdyUXBTOVd3OVgvUTVPNU8vK0gxSmF4ZnNCMElNCmZOdGxmTlhkTGs0STFmeGtjcTAvWlZoQ1Vmd2ZJT3NJNVdBaklVK3ZLQmx0QytENE82bUJtUmxJRDlUclZtcTEKb21KM0tRNjNZUFpHMVYxNTRnM2NhTU9KakpVNUFnTUJBQUV3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUNJcwp6VWsrWFk0NzJHSHQxWjl5VWdzOGkyN3pHQ0hUTUp3b2V3Y0RpL2FwQ0pNcFZHT3gvMEVsR1cxY2xySVZSbjhOCkN5a3NPaFlXbnBqVUVVRGYyZHY1SkRHSGpBK0ExTFNUUEVRYXhCTXEvOEhkekN1WFdsN2xrTDdXWW9KWWQvOWkKTFRIOUpBVkNtckh6VklLeWd3d1ZSVHIwZVhRbGJ1ZnpFd01TU0FUWFJmMTFwekorazZyVE1icmNIT2pJb3FreQpIVVZCOHJsb3RUMUgxdFBUOVVzcUhoR0N3eUdad2MwSkNSSXZwemJsdUc4ZUFCL1gxWXdmblQzbG9iQzczT2VXCnd0ZkdTd25EN3IwS2E0YWdoQVMraWtRNDdtdklIWFVOTzU3WUt3dXJkUllrMjZxQzZqRTZRM3haU1J3MC92SkIKbm1SaUl6SmJUWUQvYVluT3N6WT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.testcase, func(t *testing.T) {
			t.Parallel()
			c := &Config{}
			err := yaml.Unmarshal([]byte(tt.config), c)
			assert.NoError(t, err, "error unmarshalling")
			err = c.Validate(log)
			assert.NoError(t, err, "error validating")
			ref, err := reference.NewReference(tt.image, logrus.NewEntry(logrus.StandardLogger()))
			assert.NoError(t, err, "error parsing image")
			repos, err := c.GetMatchingRepositoriesPerImage(ref, log)
			assert.NoError(t, err, "error getting matching repositories per image", ref)
			assert.Len(t, repos, tt.expectedReposLen, "wrong repos len")
			if tt.expectedReposLen >= 1 {
				assert.Equal(t, tt.expectedRepositoryName, repos[0].Name, "wrong repo returned")
			}

		})
	}

}

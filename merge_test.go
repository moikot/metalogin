package main

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func Test_Merge_Successfully_Creates(t *testing.T) {
	srcStr := `
clusters:
- cluster:
    certificate-authority-data: cad
    server: https://host:6443
  name: kubernetes
users:
- user:
    client-certificate-data: ccd
    client-key-data: ckd
  name: kubernetes-admin
`
	expStr := `apiVersion: ""
kind: ""
current-context: kubernetes-admin-host@kubernetes-host
clusters:
- cluster:
    certificate-authority-data: cad
    server: https://host:6443
  name: kubernetes-host
contexts:
- context:
    cluster: kubernetes-host
    user: kubernetes-admin-host
  name: kubernetes-admin-host@kubernetes-host
users:
- user:
    client-certificate-data: ccd
    client-key-data: ckd
  name: kubernetes-admin-host
`

	src, _ := Unmarshal([]byte(srcStr))
	dst, _ := Unmarshal([]byte(""))

	err := Merge(src, dst)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	dstBytes, _ := yaml.Marshal(dst)

	if string(dstBytes) != expStr {
		t.Errorf("expected: \n%v \n got: \n%v", expStr, string(dstBytes))
	}
}

func Test_Merge_Successfully_Updates(t *testing.T) {
	srcStr := `
clusters:
- cluster:
    certificate-authority-data: cad_2
    server: https://host:6777
  name: kubernetes
users:
- user:
    client-certificate-data: ccd_2
    client-key-data: ckd_2
  name: kubernetes-admin
`
	dstStr := `apiVersion: v1
kind: Config
current-context: kubernetes-admin-host@kubernetes-host
clusters:
- cluster:
    certificate-authority-data: cad
    server: https://host:6443
  name: kubernetes-host
contexts:
- context:
    cluster: kubernetes-host
    user: kubernetes-admin-host
  name: kubernetes-admin-host@kubernetes-host
users:
- user:
    client-certificate-data: ccd
    client-key-data: ckd
  name: kubernetes-admin-host
`
	expStr := `apiVersion: v1
kind: Config
current-context: kubernetes-admin-host@kubernetes-host
clusters:
- cluster:
    certificate-authority-data: cad_2
    server: https://host:6777
  name: kubernetes-host
contexts:
- context:
    cluster: kubernetes-host
    user: kubernetes-admin-host
  name: kubernetes-admin-host@kubernetes-host
users:
- user:
    client-certificate-data: ccd_2
    client-key-data: ckd_2
  name: kubernetes-admin-host
`

	src, _ := Unmarshal([]byte(srcStr))
	dst, _ := Unmarshal([]byte(dstStr))

	err := Merge(src, dst)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	dstBytes, _ := yaml.Marshal(dst)

	if string(dstBytes) != expStr {
		t.Errorf("expected: \n%v \n got: \n%v", expStr, string(dstBytes))
	}
}

func Test_Merge_Fails_For_Reason(t *testing.T) {
	var tests = []struct {
		src   string // source config
		error string // expected error
	}{
		{
			src: `
clusters:
- cluster:
    certificate-authority-data: cad
    server: https://host:6443
  name: foo
`,
			error: "cluster `kubernetes` is not defined in source configuration",
		},
		{
			src: `
clusters:
- cluster:
    certificate-authority-data: cad
    server: "f://o o"
  name: kubernetes
`,
			error: "server name `f://o o` is not a valid URL: parse f://o o: invalid character \" \" in host name",
		},
		{
			src: `
clusters:
- cluster:
    server: https://host:6777
  name: kubernetes
`,
			error: "user `kubernetes-admin` is not defined in source configuration",
		},
	}

	for _, test := range tests {
		src, _ := Unmarshal([]byte(test.src))
		dst, _ := Unmarshal([]byte(""))

		err := Merge(src, dst)
		if err == nil {
			t.Errorf("expected error `%v`, got `nil`", test.error)
			return
		}

		if err.Error() != test.error {
			t.Errorf("expected error `%v`, got `%v`", test.error, err.Error())
		}
	}
}

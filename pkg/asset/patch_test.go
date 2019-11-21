package asset

import (
	"testing"
)

func TestPatch(t *testing.T) {
	t.Parallel()
	type testCase struct {
		Name         string
		As           Assets
		P            []byte
		ExpectError  bool
		ExpectOutput string
	}
	cases := []testCase{
		{
			Name:         "test basic add to list",
			As:           genAsset("testadd.yaml", []byte(testAddManifest)),
			P:            testSimpleAddPatch,
			ExpectError:  false,
			ExpectOutput: resultSimpleAdd,
		},
		{
			Name:         "test simple replace value",
			As:           genAsset("testreplace.yaml", []byte(testReplaceManifest)),
			P:            testSimpleReplacePatch,
			ExpectError:  false,
			ExpectOutput: resultSimpleReplace,
		},
		{
			Name:         "test list replace value",
			As:           genAsset("testreplace.yaml", []byte(testAddManifest)),
			P:            testListReplacePatch,
			ExpectError:  false,
			ExpectOutput: resultListReplace,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ps, _ := ReadPatch(tc.P)
			err := tc.As.Patch([]PatchJSON6902{ps})
			ExpectError(t, tc.ExpectError, err)
			if err == nil {
				StringEqual(t, tc.ExpectOutput, string(tc.As[0].Data))
			}
		})
	}
}

func genAsset(name string, data []byte) []Asset {
	testAsset := Asset{
		Name: name,
		Data: data,
	}
	return []Asset{testAsset}
}

const testAddManifest = `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kube-proxy
spec:
  template:
    spec:
      containers:
      - command:
        - ./hyperkube
        - kube-proxy
        - --cluster-cidr=10.2.0.0/16
        - --hostname-override=$(NODE_NAME)
        - --kubeconfig=/etc/kubernetes/kubeconfig
        - --proxy-mode=iptables
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 1
    type: RollingUpdate
`
const testReplaceManifest = `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kube-proxy
spec:
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 1
    type: RollingUpdate
`

var testSimpleAddPatch = []byte(`
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kube-proxy
patch: |
  - op: add
    path: /spec/template/spec/containers/0/command/-
    value: --test
`)

var testSimpleReplacePatch = []byte(`
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kube-proxy
patch: |
  - op: replace
    path: /spec/updateStrategy/rollingUpdate/maxUnavailable
    value: 2
`)

var testListReplacePatch = []byte(`
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kube-proxy
patch: |
  - op: replace
    path: /spec/template/spec/containers/0/command/5
    value: --proxy-mode=ipvs
`)

const resultSimpleAdd = `apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kube-proxy
spec:
  template:
    spec:
      containers:
      - command:
        - ./hyperkube
        - kube-proxy
        - --cluster-cidr=10.2.0.0/16
        - --hostname-override=$(NODE_NAME)
        - --kubeconfig=/etc/kubernetes/kubeconfig
        - --proxy-mode=iptables
        - --test
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 1
    type: RollingUpdate
`

const resultSimpleReplace = `apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kube-proxy
spec:
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 2
    type: RollingUpdate
`

const resultListReplace = `apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kube-proxy
spec:
  template:
    spec:
      containers:
      - command:
        - ./hyperkube
        - kube-proxy
        - --cluster-cidr=10.2.0.0/16
        - --hostname-override=$(NODE_NAME)
        - --kubeconfig=/etc/kubernetes/kubeconfig
        - --proxy-mode=ipvs
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 1
    type: RollingUpdate
`

type testingDotT interface {
	Errorf(format string, args ...interface{})
}

func ExpectError(t testingDotT, expectError bool, err error) {
	if err != nil && !expectError {
		t.Errorf("Did not expect error: %v", err)
	}
	if err == nil && expectError {
		t.Errorf("Expected error but got none")
	}
}

func StringEqual(t testingDotT, expected, result string) {
	if expected != result {
		t.Errorf("Strings did not match!")
		t.Errorf("Expected: %q", expected)
		t.Errorf("But got: %q", result)
	}
}

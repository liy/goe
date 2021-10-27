package tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func getType(o interface{}) (res string) {
    t := reflect.TypeOf(o)
    for t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    return res + t.Name()
}

func ToMatchSnapshot(t *testing.T, o interface{}) {
	actualBytes, err := json.Marshal(o)
	if err != nil {
		t.Fatal(err)
	}

	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("cannot obtain caller's name")
	}
	chunks := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	testFuncName := chunks[len(chunks)-1]

	postfix := getType(o)
	if postfix == "" {
		postfix = "_snapshost.json"
	} else {
		postfix = "_" + postfix + "_snapshost.json"
	}
	p := "./" + testFuncName + postfix
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err = os.WriteFile(p, actualBytes, 0644)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		f, _ := os.Open(p)
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		require.JSONEq(t, string(bs), string(actualBytes))
	}
}
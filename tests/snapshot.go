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

func Save(o interface{}, filename string) error {
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
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
	filename := "./" + testFuncName + postfix
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := Save(o, filename)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		f, _ := os.Open(filename)
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		require.JSONEq(t, string(bs), string(actualBytes))
	}
}
package main

import "testing"

func True(t *testing.T, value bool) {
	if !value {
		t.Error("expect true, but actual is false")
	}
}

func TrueErr(t *testing.T, value bool, errMsg string) {
	if !value {
		t.Errorf("expect true, but actual is false, error: %v", errMsg)
	}
}

func False(t *testing.T, value bool) {
	if value {
		t.Error("expect false, but actual is true")
	}
}

func FalseErr(t *testing.T, value bool, errMsg string) {
	if value {
		t.Errorf("expect false, but actual is true, error: %v", errMsg)
	} else {
		t.Logf("error: %v", errMsg)
	}
}

// Equal 参数为act-expect-act-expect-...结构，其中expect为期望值，act为实际值
func Equal(t *testing.T, objects ...interface{}) {
	if len(objects)%2 != 0 {
		t.Error("参数个数必须为偶数")
	}

	for i := 0; i < len(objects); i += 2 {
		if objects[i] != objects[i+1] {
			t.Errorf("期望：%v \n          实际：%v", objects[i+1], objects[i])
		}
	}
}

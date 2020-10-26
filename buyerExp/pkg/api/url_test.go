package api

import (
	"testing"
)

func TestCheckKeySuccess(t *testing.T){
	key := "3vT43"
	err := checkKey(key)
	if err != nil{
		t.Errorf("key %s failed check", key)
	}
}

func TestCheckEmptyKeySuccess(t *testing.T){
	key := ""
	err := checkKey(key)
	if err != nil{
		t.Errorf("key %s failed check", key)
	}
}

func TestCheckKeyRunes(t *testing.T){
	key := string(KeyRunes)
	err := checkKey(key)
	if err != nil{
		t.Errorf("key %s failed check", key)
	}
}


func TestCheckKeyFailed(t *testing.T){
	key := "РусскийКлюч"
	err := checkKey(key)
	if err == nil{
		t.Errorf("key %s did't fail check", key)
	}
}

func TestCheckMixedKeyFailed(t *testing.T){
	key := "RussianКлюч"
	err := checkKey(key)
	if err == nil{
		t.Errorf("key %s did't fail check", key)
	}
}


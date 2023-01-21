package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetUrls(t *testing.T) {
	html := `<a class="button" href="/profile.html">`

	got := GetUrls(strings.NewReader(html))
	want := []string{"/profile.html"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetUrls should be %v, got %v", want, got)
	}
}

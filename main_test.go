package main

import "testing"

func TestToFixedUrl(t *testing.T) {
	fixedUrl := toFixedUrl("/aboutus.html", "http://ajinkya.com/")
	if fixedUrl != "http://ajinkya.com/aboutus.html" {
		t.Error("toFixedUrl did not get expected href")
	}

	mailToUrl := toFixedUrl("mailto:ajinkya@gmail.com", "http://ajinkya.com/")
	if mailToUrl != "http://ajinkya.com/" {
		t.Error("expected baseUrl instead of mailto link")
	}

	telephoneUrl := toFixedUrl("tel://9820098200", "http://ajinkya.com/")
	if telephoneUrl != "http://ajinkya.com/" {
		t.Error("expected baseUrl instead of telephone link")
	}
}

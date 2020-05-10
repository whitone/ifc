// Copyright 2020 Stefano Cotta Ramusino. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ifc

import (
	"strings"
	"testing"
	"time"
)

var p Person

func TestEncode(t *testing.T) {
	_, err := Encode("", "", 0, "--", "Andria")
	if err == nil {
		t.Fatalf("no error generated: expected '%s'", errSex)
	} else {
		if !strings.Contains(err.Error(), errSex) {
			t.Fatalf("unexpected error: got '%s', expected '%s'", err, errSex)
		}
	}
}

func TestEncode2(t *testing.T) {
	_, err := Encode("", "", 'f', "--", "Andria")
	if err == nil {
		t.Fatalf("no error generated: expected '%s'", `cannot parse "--" as "2006"`)
	} else {
		if !strings.Contains(err.Error(), `cannot parse "--" as "2006"`) {
			t.Fatalf("unexpected error: got '%s', expected '%s'", err, `cannot parse "--" as "2006"`)
		}
	}
}

func TestEncode3(t *testing.T) {
	_, err := Encode("", "", 'f', "1999-12-13", "Andria")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPerson_Ifc(t *testing.T) {
	p = Person{}
	_, err := p.Ifc()
	if err == nil {
		t.Fatalf("no error generated: expected '%s'", errSex)
	} else {
		if !strings.Contains(err.Error(), errSex) {
			t.Fatalf("unexpected error: got '%s', expected '%s'", err, errSex)
		}
	}
}

func TestPerson_Ifc2(t *testing.T) {
	p.Sex = 'm'
	_, err := p.Ifc()
	if err == nil {
		t.Fatalf("no error generated: expected '%s'", `cannot parse "" as "2006"`)
	} else {
		if !strings.Contains(err.Error(), `cannot parse "" as "2006"`) {
			t.Fatalf("unexpected error: got '%s', expected '%s'", err, `cannot parse "" as "2006"`)
		}
	}
}

func TestPerson_Ifc3(t *testing.T) {
	p.BirthDate = time.Now().AddDate(-30, 0, 0).Format("2006-01-02")
	_, err := p.Ifc()
	if err == nil {
		t.Fatalf("no error generated: expected '%s'", errPlace)
	} else {
		if !strings.Contains(err.Error(), errPlace) {
			t.Fatalf("unexpected error: got '%s', expected '%s'", err, errPlace)
		}
	}
}

func TestPerson_Ifc4(t *testing.T) {
	p.BirthPlace = "Lucca "
	_, err := p.Ifc()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPerson_Ifc5(t *testing.T) {
	p.Surname = "Doe"
	_, err := p.Ifc()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPerson_Ifc6(t *testing.T) {
	p.Name = "Jane"
	_, err := p.Ifc()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPerson_Ifc7(t *testing.T) {
	p.Sex = 'F'
	_, err := p.Ifc()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecode(t *testing.T) {
	person, err := Decode("RSSMRA99T13H501A")
	if err != nil {
		t.Fatal(err)
	}
	if person.Sex != strings.ToUpper(string(male))[0] {
		t.Fatalf(errSex+": got '%c', expected '%c'", person.Sex, strings.ToUpper(string(male))[0])
	}
	if person.BirthDate != "1999-12-13" {
		t.Fatalf("wrong birth date: got '%s', expected '%s'", person.BirthDate, "1999-12-13")
	}
	if person.BirthPlace != "Roma" {
		t.Fatalf("wrong birth place: got '%s', expected '%s'", person.BirthPlace, "Roma")
	}
}

func TestPerson_CheckIfc(t *testing.T) {
	err := p.CheckIfc("DOEJHN99T13Z121S")
	if err != nil {
		t.Fatal(err)
	}
	if p.Sex != strings.ToUpper(string(male))[0] {
		t.Fatalf(errSex+": got '%c', expected '%c'", p.Sex, strings.ToUpper(string(male))[0])
	}
	if p.BirthDate != "1999-12-13" {
		t.Fatalf("wrong birth date: got '%s', expected '%s'", p.BirthDate, "1999-12-13")
	}
	if p.BirthPlace != "Malta" {
		t.Fatalf("wrong birth place: got '%s', expected '%s'", p.BirthPlace, "Malta")
	}
}

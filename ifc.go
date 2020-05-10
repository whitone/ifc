// Copyright 2020 Stefano Cotta Ramusino. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ifc is a library to handle Italian fiscal code.
package ifc

//go:generate go run generator/updateMunicipalities.go && go fmt municipalities.go

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	// A small library that turns strings in to slugs
	"github.com/Machiel/slugify"
)

const (
	length = 16

	male       = 'm'
	female     = 'f'
	femaleCode = 40
	vowels     = "aeiou"
	omocodia   = "lmnpqrstuv"
	months     = "abcdehlmprst"
	consonants = "bcdfghjklmnpqrstvwxyz"

	errSex       = "invalid sex"
	errCode      = "invalid code"
	errMonth     = "invalid month"
	errLength    = "wrong length for code"
	errCin       = "wrong control interal number"
	errPlace     = "birth place not mapped to any code"
	errPlaceCode = "birth place code not mapped to any municipality or country"

	re = `(?i)^([a-z]{3})([a-z]{3})([a-z\d]{2})([` + months + `])([a-z\d]{2})([a-z][a-z\d]{3})([a-z])$`
)

var odds = []int{1, 0, 5, 7, 9, 13, 15, 17, 19, 21, 2, 4, 18, 20, 11, 3, 6, 8, 12, 14, 16, 10, 22, 25, 24, 23}

// Person represents personal data used in Italian fiscal code.
type Person struct {
	// Surname is the last name of the person.
	Surname string
	// Name is the first name of the person.
	Name string
	// Sex can be 'M' or 'F'.
	Sex byte
	// BirthDate is the birth date in the format "YYYY-MM-DD".
	BirthDate string
	// BirthPlace is the Italian municipality or a country (in Italian language) if the birth place is outside of Italy.
	BirthPlace string
}

// CheckIfc verifies if the code is a valid Italian fiscal code and decode it to person element.
func (p *Person) CheckIfc(code string) error {
	decodedPerson, err := Decode(code)
	if err != nil {
		return err
	}
	p.Name = decodedPerson.Name
	p.Surname = decodedPerson.Surname
	p.Sex = decodedPerson.Sex
	p.BirthDate = decodedPerson.BirthDate
	p.BirthPlace = decodedPerson.BirthPlace
	return nil
}

// Decode checks an Italian fiscal code.
func Decode(code string) (p *Person, err error) {
	cleanCode := strings.ReplaceAll(slugify.Slugify(code), "-", "")
	p = &Person{}
	var codeRegExp = regexp.MustCompile(re)
	codeFields := codeRegExp.FindStringSubmatch(cleanCode)
	if codeFields == nil {
		return nil, errors.New(errCode + " '" + code + "'")
	}
	p.Surname = strings.ToUpper(codeFields[1])
	p.Name = strings.ToUpper(codeFields[2])
	birthDateYear := codeFields[3]
	birthDateMonthCode := codeFields[4]
	birthDateDayCode := codeFields[5]
	birthPlaceCode := codeFields[6]
	cin := codeFields[7]
	birthDateMonth := strings.Index(months, birthDateMonthCode) + 1
	if birthDateMonth < 1 {
		return nil, errors.New(errMonth + " '" + birthDateMonthCode + "'")
	}
	for index, char := range []rune(omocodia) {
		birthDateYear = strings.ReplaceAll(birthDateYear, string(char), strconv.Itoa(index))
		birthDateDayCode = strings.ReplaceAll(birthDateDayCode, string(char), strconv.Itoa(index))
		birthPlaceCode = string(birthPlaceCode[0]) +
			strings.ReplaceAll(birthPlaceCode[1:], string(char), strconv.Itoa(index))
	}
	birthDateDay, err := strconv.Atoi(birthDateDayCode)
	if err != nil {
		return nil, err
	}
	p.Sex = 'M'
	if birthDateDay > femaleCode {
		birthDateDay -= femaleCode
		p.Sex = 'F'
	}
	birthDateCode := birthDateYear + fmt.Sprintf("%02d%02d", birthDateMonth, birthDateDay)
	birthDate, err := time.Parse("060102", birthDateCode)
	if err != nil {
		return nil, err
	}
	p.BirthDate = birthDate.Format("2006-01-02")
	if val, ok := municipalitiesByCode[birthPlaceCode]; ok {
		p.BirthPlace = val
	}
	if p.BirthPlace == "" {
		err = errors.New(errPlaceCode + " '" + birthPlaceCode + "'")
	}
	cinCheck, err := encodeCin(cleanCode[:len(cleanCode)-1])
	if err != nil {
		return nil, err
	}
	if cin != cinCheck {
		return nil, errors.New(errCin + " '" + cin + "', expected '" + cinCheck + "'")
	}
	return p, nil
}

// Ifc creates the Italian fiscal code of the person.
func (p *Person) Ifc() (code string, err error) {
	return Encode(p.Surname, p.Name, p.Sex, p.BirthDate, p.BirthPlace)
}

// Encode generates an Italian fiscal code.
func Encode(surname string, name string, sex byte, birthDate string, birthPlace string) (code string, err error) {
	birthDateCode, err := encodeBirthDate(birthDate, sex)
	if err != nil {
		return
	}
	birthPlaceCode, err := encodeBirthPlace(birthPlace)
	if err != nil {
		return
	}
	code = encodeSurname(surname)
	code += encodeName(name)
	code += birthDateCode
	code += birthPlaceCode
	cin, err := encodeCin(code)
	if err != nil {
		return
	}
	code += cin
	_, err = Decode(code)
	if err != nil {
		return
	}
	code = strings.ToUpper(code)
	return
}

// encodeSurname encodes surname to the code used in Italian fiscal code.
func encodeSurname(surname string) string {
	surnameSlug := slugify.Slugify(surname)
	surnameConsonants := subset(surnameSlug, consonants)
	surnameVowels := subset(surnameSlug, vowels)
	return indicativeChars(surnameConsonants, surnameVowels)
}

// subset extracts from a string a restrict set described in category string.
func subset(s string, category string) (subset string) {
	for _, char := range s {
		if strings.Contains(category, string(char)) {
			subset += string(char)
		}
	}
	return
}

// indicativeChars mixes consonants and vowels according to Italian fiscal code algorithm.
func indicativeChars(consonants string, vowels string) (chars string) {
	chars = consonants
	if len(chars) >= 3 {
		return chars[:3]
	}
	chars += vowels
	if len(chars) >= 3 {
		return chars[:3]
	}
	chars += "xxx"
	return chars[:3]
}

// encodeName encodes name according to Italian fiscal code algorithm.
func encodeName(name string) string {
	nameSlug := slugify.Slugify(name)
	nameConsonants := subset(nameSlug, consonants)
	if len(nameConsonants) > 3 {
		nameConsonants = nameConsonants[:1] + nameConsonants[2:]
	}
	nameVowels := subset(nameSlug, vowels)
	return indicativeChars(nameConsonants, nameVowels)
}

// encodeBirthDate encodes birth date according to Italian fiscal code algorithm.
func encodeBirthDate(birthDateString string, sex byte) (string, error) {
	cleanSex := strings.ToLower(string(sex))[0]
	if cleanSex != female && cleanSex != male {
		return "", errors.New(errSex + " '" + string(sex) + "'")
	}
	birthDate, err := time.Parse("2006-01-02", birthDateString)
	if err != nil {
		return "", err
	}
	birthDateCode := string([]rune(strconv.Itoa(birthDate.Year()))[2:])
	birthDateCode += string([]rune(months)[birthDate.Month()-1])
	birthDayCodeNumber := birthDate.Day()
	if cleanSex == female {
		birthDayCodeNumber += femaleCode
	}
	return birthDateCode + fmt.Sprintf("%02d", birthDayCodeNumber), nil
}

// encodeBirthPlace encodes birth place according to Italian fiscal code algorithm.
func encodeBirthPlace(birthPlace string) (birthPlaceCode string, err error) {
	birthPlaceSlug := slugify.Slugify(birthPlace)
	if val, ok := municipalitiesByName[birthPlaceSlug]; ok {
		birthPlaceCode = val
	}
	if birthPlaceCode == "" {
		err = errors.New(errPlace + " '" + birthPlace + "'")
	}
	return
}

// ordinalValue get numeric value of an ASCII character.
func ordinalValue(char byte) (num int) {
	if num = int(char) - int('a'); num < 0 {
		num += int('a') - int('0')
	}
	return
}

// encodeCin calculates the Control Internal Number, the check character of the Italian fiscal code.
func encodeCin(code string) (string, error) {
	if len(code) != length-1 {
		return "", errors.New(errLength + " '" + code + "'")
	}
	cinCodeNumber := odds[ordinalValue(code[length-2])]
	for i := 0; i <= length-3; i += 2 {
		cinCodeNumber += odds[ordinalValue(code[i])] + ordinalValue(code[i+1])
	}
	return string(rune(cinCodeNumber%len(odds)) + rune('a')), nil
}

[![godoc reference][godoc badge]][godoc url]
[![build status][travis badge]][travis url]
[![go report card][goreportcard badge]][goreportcard url]
[![total alerts][lgtm badge]][lgtm url]
[![license][license badge]][license url]

# ifc

> The `ifc` package allows to handle [Italian fiscal code][ifc wikipedia].

## Example

Just import the library to generate the Italian fiscal code from the data of a person.

```go
import "github.com/whitone/ifc"

code, err := ifc.Encode("Doe", "John", 'M', "1999-12-13", "Malta")
if err != nil {
  log.Fatal(err)
}
```

The library also check if a code is valid and return all available informations of the person.

```go
person, err := ifc.Decode("DOEJHN99T13Z121S")
if err != nil {
  log.Fatal(err)
}
log.Println("sex: %c, birthdate: %s, birthplace: %s", person.Sex, person.BirthDate, person.BirthPlace)
```

Another way to use the library is to set all required informations of person and then generate the code.

```go
person := &ifc.Person{}
person.Surname = "Doe"
person.Name = "Jane"
person.Sex = 'F'
person.BirthDate = time.Now().AddDate(-30, 0, 0).Format("2006-01-02")
person.BirthPlace = "Palermo"
code, err := person.Ifc()
if err != nil {
  log.Fatal(err)
}
```

The opposite is to use the check method to fulfill the person with all its informations, be aware that any
previous data will be overwritten.

```go
person := &ifc.Person{}
err := person.CheckIfc("DOEJNA90H48G273P")
if err != nil {
  log.Fatal(err)
}
log.Println("sex: %c, birthdate: %s, birthplace: %s", person.Sex, person.BirthDate, person.BirthPlace)
```

## Useful references

- [Italian fiscal code reference documentation][ifc refdoc]

## License

`ifc` is licensed under the [BSD-3-Clause License][license url].

[godoc badge]: https://img.shields.io/badge/godoc-reference-5272B4
[godoc url]: https://pkg.go.dev/github.com/whitone/ifc
[travis badge]: https://api.travis-ci.com/whitone/ifc.svg?branch=master
[travis url]: https://travis-ci.com/github/whitone/ifc
[goreportcard badge]: https://goreportcard.com/badge/github.com/whitone/ifc
[goreportcard url]: https://goreportcard.com/report/github.com/whitone/ifc
[lgtm badge]: https://img.shields.io/lgtm/alerts/g/whitone/ifc.svg
[lgtm url]: https://lgtm.com/projects/g/whitone/ifc/alerts
[license badge]: https://img.shields.io/github/license/whitone/ifc.svg
[license url]: ./LICENSE
[ifc wikipedia]: https://en.wikipedia.org/wiki/Italian_fiscal_code
[ifc refdoc]: https://www.agenziaentrate.gov.it/portale/Schede/Istanze/Richiesta+TS_CF/Informazioni+codificazione+pf/
[![godoc reference](https://img.shields.io/badge/godoc-reference-5272B4)](https://pkg.go.dev/github.com/whitone/ifc)
[![build status](https://api.travis-ci.com/whitone/ifc.svg?branch=master)](https://travis-ci.com/github/whitone/ifc)
[![go report card](https://goreportcard.com/badge/github.com/whitone/ifc)](https://goreportcard.com/report/github.com/whitone/ifc)
[![license](https://img.shields.io/github/license/whitone/ifc.svg)](./LICENSE)

# ifc

> The `ifc` package allows to handle [Italian fiscal code].

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

- [Italian fiscal code reference documentation][IFC refdoc]

## License

`ifc` is licensed under the [BSD-3-Clause License](./LICENSE).

[Italian fiscal code]: https://en.wikipedia.org/wiki/Italian_fiscal_code
[IFC refdoc]: https://www.agenziaentrate.gov.it/portale/Schede/Istanze/Richiesta+TS_CF/Informazioni+codificazione+pf/
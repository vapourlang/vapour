![](https://vapour.run/img/vapour.png)

<div align="center">
  <p>
    A typed superset of the <a href="https://www.r-project.org/">R programming language</a>
  </p>
  <a href="https://vapour.run">Docs</a> | <a href="https://vapour.run/get-started">Get Started</a> | <a href="https://vapour.run/install">Install</a>
</div>

> [!WARNING]  
> Vapour is in (very) early alpha!

```r
type person: object {
   age: int,
   name: char 
}

func create(name: char): person {
  return person(name = name)
}

@generic
func (p: any) set_age(...: any): any

@default
func(p: any) set_age(age: int): null {
  stop("not implemented")
}

func(p: person) set_age(age: int): person {
  p$age = age
  return p
}

let john: person = create("John") |>
  set_age(36)
```

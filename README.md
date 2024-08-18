![](https://vapour.run/img/vapour.png)

<div style="text-align:center">
  <p>
    A typed superset of the [R programming language](https://www.r-project.org/)
  </p>
  <a href="https://vapour.run">docs</a> | <a href="https://github.com/vapourlang/vapour/releases">releases</a>
</div>

> [!WARNING]  
> This is a work in progress!

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

func(p: default) set_age(age: int): null {
  stop("not implemented")
}

func(p: person) set_age(age: int): person {
  p$age = age
  return p
}

let john: person = create("John") |>
  set_age(36)
```

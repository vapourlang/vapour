![](https://vapour.run/img/vapour.png)

A typed superset of the [R programming language](https://www.r-project.org/),
see the [documentation](https://vapour.run) for more information.

> [!WARNING]  
> This is a work in progress!

```r
type person: object {
  age: int,
  name: char 
}

func create(name: char): person {
  stopifnot(!missing(name))
  return person(name = name)
}

func(p: person) set_age(age: int = 42): null {
  p$age = age
}

let john: person = create("john")

set_age(john, 36)
```

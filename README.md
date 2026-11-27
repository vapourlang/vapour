![](https://vapour.run/img/vapour.png)

<div align="center">
  <p>
    Typed <a href="https://www.r-project.org/">R</a>
  </p>
  <a href="https://vapour.run">Docs</a> | <a href="https://vapour.run/get-started">Get Started</a> | <a href="https://vapour.run/install">Install</a>
</div>

> [!WARNING]  
> Vapour is in (very) early alpha!

```r
type user: object {
   age: int,
   name: char 
}

func create(name: char): user {
  return user(name = name)
}

@generic
func (p: any) set_age(age: int, ...: any): any

@default
func(p: any) set_age(age: int): null {
  stop("not implemented")
}

func(p: user) set_age(age: int): user {
  p$age = age
  return p
}

let john: user = create("John") |>
  set_age(36)
```

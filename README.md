# ReqParser
multiparser for fiber request parser interface

# Idea
main idea was make the request more clean, rather than:
### before
```go
var request RequestStruct
err := c.BodyParser(&request)
if err != nil {
    return err
}
```
### after
```go
request, err := reqparser.New[RequestStruct](c).Parse(reqparser.BodyParser())
if err != nil {
    return err
}
```

### or
```go
id, err := reqparser.New[string](c).Parse(reqparser.Params("id", true))
if err != nil {
    return err
}
```

# License
[GNU](https://github.com/itzngga/reqparser/blob/master/LICENSE)

# Contribute
Pull Request are pleased to
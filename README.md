# resourcename

> parse google-type resource names to convert from DTOs to domain models.

Google's [AIP](https://google.aip.dev) defines standards on how APIs should be written.
[AIP-122](https://google.aip.dev/122) mandates that each resource must be have a resource name, by which the resource is identified.
That name must be of the format `resources/{resource_id}`

## Usage

When defining a domain resource, you can use golang field tags to specify a resource name segment (rns) for this field.
For now, only string fields are supported.

```go
type Resource struct {
    Id string `rns:"resource_id"`
}
```

You will also need to define the pattern of your resource's resource name.
The variable segments of the resource name is what we will map the tagged fields to.
Make sure, that you wrap variable segments in curly braces and that the name matches what you set in the tag.

```go
pattern := resourcename.FromPattern("/resources/{resource_id}")
```

We can use our pattern to parse a string resource name and assign it's variable segments to our resource struct.
Make sure you pass your resource by reference.

```go
resource := Resource{}
err := pattern.MatchInto("/resources/abcdefg", &resource)
```

We also support any number of parent resources.
Usually you don't have grandparent ids in your model, but this is just to prove that it would be supported.
But our field tags would look like:

```go
type Child struct {
    Id            string `rns:"child_id"`
    ParentId      string `rns:"parent_id"`
    
    // this may not make sense if parent has the reference to grandparent
    // but works if need be
    GrandparentId string `rns:"grandparent_id"` 
}

```

```go
child := Child{}

rname := "grandparents/abcd/parents/efgh/children/ijkl"

pattern := resourcenames.FromPattern(
    "grandparents/{grandparent_id}/parents/{parent_id}/children/{child_id}")

err := pattern.MatchInto(rname, &child)
if err != nil {
    return err
}
```
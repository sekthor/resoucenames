# resourcename

Google's [AIP](https://google.aip.dev) defines standards on how APIs should be written.
[AIP-122](https://google.aip.dev/122) mandates that each resource must be have a resource name, by which the resource is identified.
That name must be of the format `resources/{resource_id}`.
In the domain layer however, your domain resources will typically work with `id` only, rather than the full resource name.
Parsing the id(s) from a resource name can be tedious and repetitive.
With this package, they can be injected into a resource from a pattern and a resource name using golang struct tags.

## Problem Example

This is how the AIP defines api resources (example from AIP).
The primary identifier is not an id, but a full resource name (`/publishers/{publisher}/books/{book}`), including parent resource name and all.

```protobuf
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  // The resource name of the book.
  // Format: publishers/{publisher}/books/{book}
  string name = 1 [(google.api.field_behavior) = IDENTIFIER];

  // Other fields...
}
```

Also, all requests identify the resource by it's full resource name:

```protobuf
message GetBookRequest {
  // The name of the book to retrieve.
  // Format: publishers/{publisher}/books/{book}
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference) = {
      type: "library.googleapis.com/Book"
    }];
}
```

The resource may also have a parent (or even grandparent) resource.
It's resource name can become very "nested".

```
grandparents/{grandparent_id}/parents/{parent_id}/children/{child_id}
```

But generally, our domain models just work with the ids.
We don't pass full resource names as primary keys to the database, for example.
This package help translating api messages to domain layer objects by unmarshaling resource names into tagged domain structs.

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
err := pattern.Unmarshal("/resources/abcdefg", &resource)
```

Of course we can also do the reverse

```go
resource := Resource {
  Id: "abcdefg"
}

pattern := resourcename.FromPattern("/resources/{resource_id}")
resourceName, err := pattern.Marshal(&resource)
if err != nil {
  // ...
}
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

We can unmarshal a child resource just like any other resource.

```go
child := Child{}

rname := "grandparents/abcd/parents/efgh/children/ijkl"

pattern := resourcenames.FromPattern(
    "grandparents/{grandparent_id}/parents/{parent_id}/children/{child_id}")

err := pattern.Unmarshal(rname, &child)
if err != nil {
    return err
}
```
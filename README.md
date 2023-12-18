Zombie Survival Social Network
==============================

This is a simplistic implementation of the ZSSN case.

## Architecture

From an architectual point of view, this is a mixture of [benbjohnson/wtf][wtfdial] and [Vertical Slice Architecture][vertical-slice-architecture].

The core domain is kept pure and located in the root of the project, primarily in `zssn.go`.
Alongside the domain are defined the various features of the application.
Features are defined as vertical slices.

> So what is a "Vertical Slice Architecture"? In this style, my architecture is built around distinct requests, encapsulating and grouping all concerns from front-end to back. You take a normal "n-tier" or hexagonal/whatever architecture and remove the gates and barriers across those layers, and couple along the axis of change:
>
> When adding or changing a feature in an application, I'm typically touching many different "layers" in an application. I'm changing the user interface, adding fields to models, modifying validation, and so on. Instead of coupling across a layer, we couple vertically along a slice. **Minimize coupling between slices, and maximize coupling in a slice.**
- [Jimmy Bogard, "Vertical Slice Architecture"][vertical-slice-architecture]

I've made a slight compromise on this approach, by defining a simplistic repository pattern, which is used for the repeated task of mapping to/from persistent storage.
Typically, Vertical Slice Architecture is implemented with OR/M libraries in mind for this - which I'm not a huge fan of.
Eventually, I'd like to experiment with changing this to something like [sqlc][sqlc], but that's left out for now.

When, inevitably, ZSSN becomes the Facebook of the post-apocalypse world and more sub-domains are identified, features may simply start to be grouped into packages reflecting those sub-domains.

[wtfdial]: https://github.com/benbjohnson/wtf/
[vertical-slice-architecture]: https://www.jimmybogard.com/vertical-slice-architecture/
[sqlc]: https://sqlc.dev/

## Building and Running

The most simple approach is just running the `zssnd` binary.

```
go run ./cmd/zssnd
```

Configuration is provided in a flag-first style, with each flag being conventionally set-able through environment variables as well.
See the `.env.template` file for an example, and consider something like [`direnv`][direnv] for automatically loading the environment.

Alternatively, a container recipe is provided:

```
docker build -t zssn/zssnd -f ./cmd/zssnd/Dockerfile
```

To use Postgres for persistent storage, the database schema is found at `postgres/schema.sql`.

Alternatively, a container recipe is provided, which takes care of the initial schema creation:

```
docker build -t zssn/postgres ./postgres
```

[direnv]: https://direnv.net/

## Usage

Once `zssnd` is up and running, HTTP requests are served as configured.
For convenience, [hurl][hurl] files are provided for the different features - think of it as a plaintext alternative to PostMan etc.

For more information about the various endpoints, see the provided OpenAPI 3 specification at `cmd/zssnd/openapi.yaml`.

## TODO

* Fix casing in HTTP responses
* Clean up `./cmd/zssnd/zssnd.go`
* Add integration tests for `./postgres`
* Add deployment to Google Cloud Run

[hurl]: https://hurl.dev/

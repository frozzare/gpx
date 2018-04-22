# gpx

Execute go package binaries.

## Installation

```
go get -u github.com/frozzare/gpx
```

## Usage

Example usage of [hello](https://github.com/golang/example/tree/master/hello) example with gpx:

```
gpx github.com/golang/example/hello
```

To delete binary after execution:

```
gpx -r github.com/golang/example/hello
```

To log what happens:

```
gpx -v github.com/golang/example/hello
```

To dry run (and log everything):

```
gpx -n github.com/golang/example/hello
```

Print help output:

```
gpx --help
```

## License

MIT © [Fredrik Forsmo](https://github.com/frozzare)
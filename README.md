# zzh

Simple SSH wrapper

```sh
go install github.com/tsukinoko-kun/zzh@latest
```

```sh
brew install tsukinoko-kun/tap/zzh
```

## Usage

Add a new configuration

```sh
zzh set --user <user> --password <password> --host <host> --port <port>
```

Login to a host

```sh
zzh <user>@<host>
```

or

```sh
zzh <user>@<host>:<port>
```

# Installing Go

This guide covers the installation of Go 1.21.3 on a Linux x86\_64 machine. For
other setups, please refer to the [official installation instructions](https://go.dev/doc/install).

You can install Go to /usr/local by running the following commands:

```console
$ curl -LO https://go.dev/dl/go1.21.3.linux-amd64.tar.gz
$ sudo tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
$ rm -f go1.21.3.linux-amd64.tar.gz
```

The full Go distribution will be present in /usr/local/go. You might
need to add /usr/local/go/bin to your PATH variable in order to
run the `go` binary:

```console
$ echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile && source ~/.profile
```

Once this is done, you should be able to confirm that Go works by running:

```console
$ go version
go version go1.21.3 linux/amd64
```

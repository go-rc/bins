# Bins

Bins is a small [Up](https://github.com/apex/up) application for building Golang binaries on-demand, allowing people to install open-source Go programs without installing Go, and without requiring the author to provide cross-compiled binaries.

Here's an example using [node-prune](https://github.com/tj/node-prune):

```sh
$ curl URL/github.com/tj/node-prune/cmd/node-prune > /usr/local/bin/node-prune
$ chmod +x /usr/local/bin/node-prune
```

The initial build may take 10-15s, however subsequent requests are cached by CloudFront for a configurable period (for example one day).

This project is just conceptual for now, but it would be nice to provide it as a free service at some point when money is not a concern (CDN data transfer for large assets is expensive). Security concerns would have to be addressed as well.

---

[![GoDoc](https://godoc.org/github.com/apex/bins?status.svg)](https://godoc.org/github.com/apex/bins)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)

<a href="https://apex.sh"><img src="http://tjholowaychuk.com:6000/svg/sponsor"></a>

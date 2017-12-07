# Bins

Bins is a small [Up](https://github.com/apex/up) app for serving Golang binaries cross-compiled on-demand. Consumers of your program do not need Go installed, nor do you need to cross-compile binaries for each release and platform.


## Example

Here's an example using [node-prune](https://github.com/tj/node-prune):

```sh
$ curl URL/github.com/tj/node-prune/cmd/node-prune > /usr/local/bin/node-prune
$ chmod +x /usr/local/bin/node-prune
```

The initial build may take 10-15s, however subsequent requests are cached by CloudFront for a configurable period (for example one day).

## About

This project is just conceptual for now, but it would be nice to provide it as a free service at some point when money is not a concern (CDN data transfer for large assets is expensive). Security concerns would have to be addressed as well.

In general I feel that we have the necessary information and tooling to make this kind of a thing a reality (and perfectly secure), while taking the burden away from the program author.

## Deploy

```sh
$ up
```

---

![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)

<a href="https://apex.sh"><img src="http://tjholowaychuk.com:6000/svg/sponsor"></a>

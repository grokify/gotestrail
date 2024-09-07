# Go TestRail

[![Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

Go Client for [TestRail API](https://support.testrail.com/hc/en-us/categories/7076541806228-API-Manual).

## Usage

### Direct

```
import "github.com/grokify/gotestrail"

func main() {
    client, err := gotestrail.NewClient("https://mydomain.testrail.io/", "myusername", "mypassword")
}
```

### GoAuth Credentials File

#### .goauth.json

Create a file, e.g. `.goauth.json` to contain your credentials with [`GoAuth`](https://github.com/grokify/goauth), e.g:

```
{
    "credentials": {
        "TESTRAIL": {
            "type": "basic",
            "service": "testrail",
            "basic": {
                "serverURL": "https://<mydomain>.testrail.io/",
                "username": "<myusername>",
                "password": "<mypassword>"
            }
        }
    }
}
```

#### Code

```
import "github.com/grokify/gotestrail"

func main() {
    // ... get `goauth.Credentials`
    client, err := gotestrail.NewClientFromGoauthCredentials(creds) // `creds` is a `goauth.Credentials{}`
}
```

## Related Modules

1. [`github.com/educlos/testrail`](https://github.com/educlos/testrail)
1. [`github.com/qba73/tr`](https://github.com/qba73/tr)

 [used-by-svg]: https://sourcegraph.com/github.com/grokify/gotestrail/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/grokify/gotestrail?badge
 [build-status-svg]: https://github.com/grokify/gotestrail/workflows/test/badge.svg
 [build-status-url]: https://github.com/grokify/gotestrail/actions/workflows/test.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/gotestrail
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/gotestrail
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/gotestrail
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/gotestrail
 [loc-svg]: https://tokei.rs/b1/github/grokify/gotestrail
 [repo-url]: https://github.com/grokify/gotestrail
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/gotestrail/blob/master/LICENSE

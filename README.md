# Keruu

Aggregating RSS/Atom feeds to a single HTML page.

## Installation

First, you need to install [Go](https://golang.org/dl/) version 1.12 or higher.
After that, you can use `go get` to install Keruu:

    $ go get github.com/Lepovirta/keruu

The executable should now be in path `$GOPATH/bin/keruu` (or `~/go/bin/keruu`).

## Usage

Keruu accepts the following CLI flags:

* `-config`: Path to the configuration file (default: read from STDIN)
* `-output`: Path to the HTML output file (default: write to STDOUT)
* `-help`: Displays how to use Keruu

## Configuration

Keruu is configured using YAML. Here's all the configurations accepted by Keruu:

* `feeds`: a list of RSS/Atom feeds to aggregate. At least one feed must be provided.
* `fetch`: a map containing the following configurations
  * `httpTimeout`: duration for how long to wait for a single feed fetch
* `aggregation`:
  * `title`: title to use in the HTML output
  * `description`: description to use in the HTML output
  * `maxPosts`: maximum number of posts to include in the HTML output
  * `css`: custom CSS for the HTML output

Everything except the list of feeds is optional.

## License

GNU General Public License v3.0

See LICENSE file for more information.

# Keruu

Aggregating RSS/Atom feeds to a single HTML page.

## Installation

First, you need to install [Go](https://golang.org/dl/) version 1.12 or higher.
After that, you can use `go get` to install Keruu:

    $ go get gitlab.com/lepovirta/keruu

The executable should now be in path `$GOPATH/bin/keruu` (or `~/go/bin/keruu`).

## Usage

Keruu accepts the following CLI flags:

* `-config`: Path to the configuration file (default: read from STDIN)
* `-output`: Path to the HTML output file (default: write to STDOUT)
* `-help`: Displays how to use Keruu

## Configuration

Keruu is configured using YAML. Here's all the configurations accepted by Keruu:

* `feeds`: A list of RSS/Atom feeds to aggregate. At least one feed must be provided.
  * `name`: Name of the feed (optional)
  * `url`: URL for the feed
* `fetch`: A map containing the following configurations
  * `httpTimeout`: Duration for how long to wait for a single feed fetch (optional)
* `aggregation`:
  * `title`: Title to use in the HTML output (optional)
  * `description`: Description to use in the HTML output (optional)
  * `maxPosts`: Maximum number of posts to include in the HTML output (optional)
  * `css`: Custom CSS for the HTML output (optional)
* `links`: A list of links to generate per feed item. (optional)
  * `name`: A name to display for the link
  * `url`: An URL pattern to use for the link.
    In the pattern, `$TITLE` will be replaced with the feed item title,
    and `$URL` will be replaced by the feed item link.

Everything except the list of feeds is optional.

## Example

See [keruu-jkpl](https://gitlab.com/lepovirta/keruu-jkpl) for an example of Keruu in action.

## License

GNU General Public License v3.0

See LICENSE file for more information.

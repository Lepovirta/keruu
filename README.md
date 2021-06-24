# Keruu

Aggregating RSS/Atom feeds to a single HTML page.

## Installation

First, you need to install [Go](https://golang.org/dl/) version 1.12 or higher.
After that, you can use `go get` to install Keruu:

    $ go get gitlab.com/lepovirta/keruu

The executable should now be in path `$GOPATH/bin/keruu` (or `~/go/bin/keruu`).

Alternatively, if you have [Docker](https://docker.com/) installed, you can run the Docker image from the following location:

```
registry.gitlab.com/lepovirta/keruu
```

## Usage

Keruu accepts the following CLI flags:

* `-config`: Path to the configuration file (default: read from STDIN)
* `-output`: Path to the HTML output file (default: write to STDOUT)
* `-help`: Displays how to use Keruu

## Configuration

Keruu is configured using YAML. Here's all the configurations accepted by Keruu:

* `feeds`: A list of RSS/Atom feeds to aggregate. At least one feed must be provided.
  * `name` (optional): Name of the feed
  * `url`: URL for the feed
  * `exclude` (optional): A list of regular expression patterns to match against the feed post titles.
    If a post title matches any of the expressions, the post is excluded from the HTML output.
  * `include` (optional): A list of regular expression patterns to match against the feed post titles.
    Only posts that match the expressions are included in the HTML output unless they match the expressions in the `exclude` list.
* `fetch` (optional): A map containing the following configurations
  * `httpTimeout` (optional): Duration for how long to wait for a single feed fetch
* `aggregation` (optional):
  * `title` (optional): Title to use in the HTML output
  * `description` (optional): Description to use in the HTML output
  * `maxPosts` (optional): Maximum number of posts to include in the HTML output
  * `css` (optional): Custom CSS for the HTML output
* `links` (optional): A list of links to generate per feed item.
  * `name`: A name to display for the link
  * `url`: An URL pattern to use for the link.
    In the pattern, `$TITLE` will be replaced with the feed item title,
    and `$URL` will be replaced by the feed item link.

## Example

See [keruu-example](https://gitlab.com/lepovirta/keruu-example) for an example of Keruu in action.

## Docker

There's two Docker images available for Keruu.

### Minimal

Image tag: `registry.gitlab.com/lepovirta/keruu`

The minimal image is optimized for the use in the command-line and scripting.
Besides the Keruu tool, it only includes the bare minimum system dependencies.
Running the container runs the Keruu tool directly and the container parameters are passed to the tool.

Example:

```
docker run -v $(pwd):/workspace:z \
  registry.gitlab.com/lepovirta/keruu \
  -config /workspace/config.yaml \
  -output /workspace/output.html
```

### CI

Image tag: `registry.gitlab.com/lepovirta/keruu-ci`

The CI image is optimized for the use in CI pipelines that use containers as pipeline steps such as Gitlab CI.
In addition to the Keruu tool, it includes the basic Linux utilities from Ubuntu, curl, Git, and a few other useful tools.
By default, the container runs Bash, which is often used for executing CI step scripts.

Example use in Gitlab CI with Gitlab Pages:

```
pages:
  stage: publish
  image: registry.gitlab.com/lepovirta/keruu:ci
  before_script:
  - mkdir -p public
  script:
  - keruu -config config.yaml -output public/index.html
  artifacts:
    paths:
    - public
  rules:
  - if: $CI_COMMIT_BRANCH == "master"
```

## License

GNU General Public License v3.0

See LICENSE file for more information.

# Contentful Hugo Extractor

<img src="https://d33wubrfki0l68.cloudfront.net/21d38ec2ccdfaacf6adc0b9921add9d18406493a/e1bcd/assets/images/logos/contentful-dark.svg" width="180" /> &nbsp; <img src="https://gohugo.io/img/hugo-logo.png" width="150" />

This tool extracts all content from your Contentful space and makes it easily consumable by Hugo. You can run it locally or as part of a CI server like Travis.

## Install

### Go Install Method

Assuming Go (1.10 +) is installed as well as [dep](https://golang.github.io/dep/)
```
go get -u github.com/icyitscold/contentful-hugo
cd $GOPATH/src/github.com/icyitscold/contentful-hugo
dep ensure
go install
```

## Usage

```
contentful-hugo [Flags]

Flags:
 --space-id string         Id of the contentful space from which to extract content. If not present will default to an environment variable named 'CONTENTFUL_API_SPACE'
 --api-key string          API Key used to authenticate with contentful If not present will default to an environment variable named 'CONTENTFUL_API_SPACE'
 --config-file string      Path to the config TOML file to load. Defauls to ./extractor-config.tml
 ```

The tool requires two parameters to work, a contentful space id and API key. These can be provided as command line flags or as environment variables

```
export CONTENTFUL_API_KEY=YOUR-ACCESS-KEY-HERE
export CONTENTFUL_API_SPACE=YOUR-ID-HERE
contentful-hugo
```

```
contentful-hugo --space-id [YOUR-ID-HERE] --api-key [YOUR-ACCESS-KEY-HERE] --config-file ./export-conf.toml

```

## Expected output

Contentful Hugo Extractor stores all content under the /content directory. For each content type, it makes a subdirectory. For each item, it creates a markdown file with the all properties in TOML format.

Special cases:
 - Items of type Homepage are stored as /content/_index
   - Note that there can only be one such item
 - Fields named mainContent are used for the main content of the markdown file
 - File names are based on the ID of the item to make it easily referencable from related items (for the machine, not humans)

## Configuration
Use the `--config-file` command line flag to provide the location of a TOML configuration to laod or ensure that there is a `extract-config.toml` file in the work directory of contentful-hugo

### Configure YAML output

While the default output is in TOML format, it is also possible to output content in YAML format. Use the following key in your config file:

```
encoding = "yaml"
```

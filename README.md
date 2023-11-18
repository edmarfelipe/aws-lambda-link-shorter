# Link Shorter

AWS Lambda function to create short links

[![build](https://github.com/edmarfelipe/aws-lambda/actions/workflows/deploy.yml/badge.svg?branch=main)](https://github.com/edmarfelipe/aws-lambda/actions/workflows/deploy.yml)
[![codecov](https://codecov.io/gh/edmarfelipe/aws-lambda/graph/badge.svg?token=vBFD7NlU5Q)](https://codecov.io/gh/edmarfelipe/aws-lambda)
[![Go Report Card](https://goreportcard.com/badge/github.com/edmarfelipe/aws-lambda)](https://goreportcard.com/report/github.com/edmarfelipe/aws-lambda)

## Endpoints

### Create short link

```js
POST /shorten
```

Supported attributes:

| Attribute                | Type     | Required  |
|:-------------------------|:---------|:----------|
| `title`                  | string   |    Yes    |
| `link`                   | string   |    Yes    |

Example:

```js
curl --request POST \
  --url {{baseURL}}/shorten \
  --header 'Content-Type: application/json' \
  --data '{
	"title": "My Link",
	"link": "http://www.example.com",
}'
```

### Redirect to original link

```js
GET /{{shortLink}}
```

Example:

```js
curl --request GET \
  --url {{baseURL}}/{{shortLink}}
```
<h1>🔗 b1a link shortener</h1>

**`Golang & PostgreSQL`**

This http server can shorten a long url and save it. It can also show you the links you last visited and shared with you. When sharing links, you can add a subject and message for the recipient.
<h2>Get started</h2>

To use all features, it is recommended to register

<hr>
<h3>/register</h3>

/register takes json in POST request. 

Headers:

```Content-Type: application/json```

Body:
```json
{
    "email": "users@gmail.com",
    "password": "1"
}
```

If response status code is 200, it returns nothing, otherwise an error(not json)
<hr>
<h3>/login</h3>

Now that you have registered, you need to log into your account to receive an authentication token. It accepts the same options as /register.

Headers:

```Content-Type: application/json```

Body:
```json
{
    "email": "users@gmail.com",
    "password": "1"
}
```

If response status code is 200, it returns json with message and token, otherwise an error(not json)
Example of resposne json:
```json
{
	"message": "Signed-in successful",
	"token": "token"
}
```
<hr>
<h3>/user</h3>

/user provides information about user, such as email, last visited 5 urls and last shared with him 5 urls.

Headers:
```
Authorization: token
```

Response body:
```json
{
	"email": "example@gmail.com",
	"history": [
		"I",
		"SXtP8",
		"r",
		"4La",
		"EF"
	],
	"shared": [
		"EF",
		"I",
		"d",
		"JdMgG",
		"TKSi"
	]
}
```
<hr>
<h3>/shorten</h3>

/shorten shortens long url.

Headers:
```
Content-Type: application/json
Authorization: token (optional)
```

Body:
```json
{
	"url":"https://goole.com",
	"sharedWith":"g@mail.ru",
	"topic":"Google",
	"message":"Here it is google website",
	"prefered":"google"
}
```

In body only url is required, others are optional to add. Prefered- prefered shorten url.
As a response it return json if response status code is 200, otherwise an error(not json)
For example:
```json
{
	"email": "example@gmail.com",
	"shorturl": "qk8v"
}
```
<hr>
<h3>/{shorten url}</h3>

Returns long url from short

Headers:
```
Authorization: token
```

Response body:
```json
{
	"email": "example@gmail.com",
	"url": "https://goole.com"
}
```
<hr>
<h3>/url</h3>

/url returns information about shorten url 

Query:
```
url: {shortenUrl}
```

Response body:
```json
{
	"url": {shortenUrl},
	"topic": "",
	"message": "",
	"clicked": "5"
}
```

<hr>

TODO:
- private urls
- public long urls should not be repeated
- website 

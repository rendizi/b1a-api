Requests: 


<h3>/register</h3>

Body:
```json
{
    "email": "example@gmail.com",
    "password": "12345678"
}

```
Headers:
```
Content-Type: application/json
```
Method: POST

Returns error or nothing if response status code is 200 

<h3>/login</h3>

Body:
```json
{
	"email":"example@gmail.com",
	"password":"12345678"
}
```

Headers:

```
Content-Type: application/json
```

Returns an error or json 
```json
{
	"message": "Signed-in successful",
	"token": <some auth token>
}
```

<h3>/user</h3>

Headers:


```
Authorization: <your auth token>
```

Returns an error or json if status code is 200 
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

<h3>/shorten</h3>

Body:

```json
{
	"url":"https://goole.com",
	"sharedWith":"someone@gmail.com", //optional
	"topic":"Google", //optional
	"message":"Here is google website" //optioanl
}
```
Headers (optional):

TODO: make a swagger instead and finish documentating

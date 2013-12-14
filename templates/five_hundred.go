package templates

const SERVER_ERROR_TEMPLATE = `
<html>
<head>
<style>
body {
	font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
	padding: 10px;
	color: #333;
}
pre {
	padding: 10px;
	background-color: white;
	border: 1px solid #333;
	border-radius: 4px;
}
</style>
</head>
<body>
<h1>Wowsers!</h1>
<h2>This request will self-destruct.</h2>
<p>You wouldn't see this if debug were false -- we'd render your own 500.html instead.</p>
<pre>
{{.}}
</pre>
</body>
</html>
`

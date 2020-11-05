# templatex

Package templatex implements additional go template utilities.

For now, it contains the following utilities:

* [Namespaces](#Namespaces)

## Namespaces

Namespaces implements a hierarchical template parser.

It takes a directory containing templates, parses it and creates a hierarchy of
templates where templates in child directories contain all templates parsed
along the path to that child template.

This registers namespace paths to parsed child template directories in
Namespaces by which they can later be addressed.

### Example

Given the following directory structure:
```
/home
/home/index.html
/index.html
/header.html
```

And the content of files:

`/index.html`
```html
<!DOCTYPE html>
<head>
	<title>{{ .Title }}</title>
</head>
<body>
	{{ template "header" . }}
	{{ block "content" . }}<p>Default root content.</p>{{ end }}
</body>
</html>
```

`/header.html`
```html
{{ define "header" }}<header>This is the header template.</header>{{ end }}
```

`/home/index.html`
```html
{{ define "content" }}<p>This is the Home template.</p>{{ end }}
```

Running following code:

```Go
ns, err := ParseRoot(".", "index", ".html")
if err != nil {
	log.Fatal(err)
}
if err := tt.ExecuteNamespace(buos.Stdout, "/", struct{Title string}{"Hello!"}); err != nil {
	t.Fatal(err)
}
```

Yields:
```html
<!DOCTYPE html>
<head>
        <title>Hello!</title>
</head>
<body>
        <header>This is the header template.</header>
        <p>This is the Home template.</p>
</body>
</html>
```


## License

MIT. See included LICENSE file.
## Statikit

This is an older project I was working on.

The goal is to provide a way to render static .gohtml files with dynamic data.

## Example Usage

For example, if the user has the following file tree under a directory named `project`:

* _statikit/
    * schema/
        * animals/
            * index.toml
        * index.toml
    * config.toml
* animals/
    * index.gohtml
* templates/
    * head.html
* index.gohtml

With the following data for each file:

*_statikit/config.toml*
```toml
Ignore = ["templates"]
```

*_statikit/schema/index.toml*
```toml
[Data]
Name = "Golang"
[FileSub]
Head = "templates/head.html"
```

*_statikit/schema/animals/index.toml*
```toml
[Data]
Name = "Gopher"
[FileSub]
Head = "templates/head.html"
```

*templates/head.html*
```html
<head><title>Golang rocks!</title></head>
```

*index.gohtml*
```html
{{.FileSub.Head}}
<p>A fun language to use for programming is {{.Data.Name}}.</p>
```

*animals/index.gohtml*
```html
{{.FileSub.Head}}
<p>A cute animal is the {{.Data.Name}}.</p>
```

And the program is run as follows:

`statikit render project`

Then the following files are produced:

* animals/
    * index.html
* index.html

With contents:

*index.html*
```html
<head><title>Golang rocks!</title></head>
<p>A fun language to use for programming is Golang.</p>
```

*animals/index.html*
```html
<head><title>Golang rocks!</title></head>
<p>A cute animal is the Gopher.</p>
```

## Intent

The project was intended to provide a way to create a blog, and to add new blogposts all one would have to do is add a new blogpost TOML file and re-render.

## Notable features
* The project uses [afero](https://github.com/spf13/afero) to be able to generalize the file system. This could provide the ability to render directly into cloud storage such as a github pages repo, or an azure file blob storage (unfinished)
* The renderer uses goroutines to parallelize the rendering process and make it more efficient.  (finished)
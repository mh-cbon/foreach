# foreach

bin util to iterate over a list of separated things and call for an external command `foreach` thing found.

# usage

```sh
echo "this is a thing" | foreach - as word '\s+' echo "{{.word}}"
echo "many wonders\nint this\nworld!" | foreach - as line '\n' echo "{{.index}}: {{.line}}"
```

# install

```sh
go get github.com/mh-cbon/foreach
```

# todo

- implement

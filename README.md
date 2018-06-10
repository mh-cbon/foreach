# foreach

bin util to iterate over a list of separated things and call for an external command `foreach` thing found.

# usage

```sh
$ foreach -h
Usage of foreach:
  -funcs string
      the set of functions to import for the template processing, one of sprig|gtf (default "sprig")
  -i	regexp match case insensitive (default true)
  -kind string
      the kind of html processing, one of text|html (default "html")
  -s	enable strict variable naming
```

# Example

```sh
$ echo "this is a thing" | foreach - as word '\s+' echo "{{.word}}"
this
is
a
thing

$ echo -e "many wonders\nint this\nworld!" | foreach - as line '\n' echo "{{.index}}: {{.line}}"
0: many wonders
1: int this
2: world!

$ echo "this is a thing" | foreach  - as word '\s+' echo -n "{{.word |upper}} "
THIS IS A THING

$ echo "this is a thing" | foreach -kind text  - as word '\s+' echo -n "{{.wordf}}-"
<no value>-<no value>-<no value>-<no value>-

$ echo "this is a thing" | foreach -kind html - as word '\s+' echo -n "{{.wordf}}-"
----

$ echo "this is a thing" | foreach - as word '\s+' echo -n "{{.wordf}}-"
----

$ ls test/src/ | foreach - as file '\s+' cp -v test/src/{{.file}} test/dst/{{.file}}
'test/src/1' -> 'test/dst/1'
'test/src/2' -> 'test/dst/2'

$ echo "this is a thing" | foreach -s - as word '\s+' echo -n "{{.wordf}}-"
panic: template: :1:2: executing "" at <.wordf>: map has no entry for key "wordf"

goroutine 1 [running]:
html/template.Must(0xc4200a0b40, 0x5c3720, 0xc4200bc120, 0x571820)
	/home/mh-cbon/.gvm/gos/go1.10/src/html/template/template.go:372 +0x54
main.mustExecTemplate(0x7ffca90e821d, 0xb, 0x7ffca90e820c, 0x4, 0xc4200c80e0, 0x4, 0x0, 0x1, 0xc4200c80f0, 0x2)
	/home/mh-cbon/gow/src/github.com/mh-cbon/foreach/main.go:132 +0x2bf
main.main()
	/home/mh-cbon/gow/src/github.com/mh-cbon/foreach/main.go:92 +0x524

  $ echo "this is a thing" | foreach  - as word '\s+' echo -n "{{.word |upper}-"
  panic: template: :1: unexpected "}" in operand

  goroutine 1 [running]:
  html/template.Must(0x0, 0x6c6900, 0xc420096a10, 0x0)
  	/home/mh-cbon/.gvm/gos/go1.10/src/html/template/template.go:372 +0x54
  main.mustExecTemplate(0x6930f1, 0x4, 0x7ffdb88de20a, 0x10, 0x7ffdb88de1f9, 0x4, 0xc4200aca78, 0x4, 0x0, 0xc420098ea0, ...)
  	/home/mh-cbon/gow/src/github.com/mh-cbon/foreach/main.go:148 +0x565
  main.main()
  	/home/mh-cbon/gow/src/github.com/mh-cbon/foreach/main.go:110 +0x6bb
  exit status 2

```

# install

```sh
go get github.com/mh-cbon/foreach
```


# Tpl

__Language agnostic template processor__

Command line tool to process templates in a language agnostic way.

Its main use case is to be used as lightweight configuration management, able to keep local 
configuration files up-to-date using data from:

- environment variables
- user defined variables
- other files' content

Default settings for placeholder formats and indentation help to have valid input and output files in most cases:

`$ cat testdata/yaml/template.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-nginx
spec:
  selector:
    matchLabels:
      run: my-nginx
  replicas: 2
  template:
    #!-- tpl:{file:input.yaml} --#
    spec:
      containers:
      - name: my-nginx
        image: nginx
        ports:
        - containerPort: 80
```

`$ cat testdata/yaml/input.yaml`

```yaml
metadata:
  labels:
    run: my-nginx
```

`$ tpl -s testdata/yaml/template.yaml -d output.yaml`

`$ cat output.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-nginx
spec:
  selector:
    matchLabels:
      run: my-nginx
  replicas: 2
  template:
    metadata:
      labels:
        run: my-nginx
    spec:
      containers:
      - name: my-nginx
        image: nginx
        ports:
        - containerPort: 80
```

**Placeholders**

Tpl allows a full placeholders configuration at run time using command line arguments.

Moreover, `tpl` has safe defaults as it is file type aware:
if not supplied, start and end delimiters are choosen based on the file type,
i.e. `<!--` and `-->` for xml, `#!--` and `--#` for yaml and so on...

Its default placeholders are a proper comment block, this way a template file will be always a valid file:

```xml
 <foo><!-- comment --> <!-- tpl:{var:property.key} --><foo>
```

```yaml
foo:
  test: true
  #- tpl:{env:TEST} .#
  ping: pong
```

If the file type is unknown, the default delimiters are `${` and `}`.

Placeholder _must_ be in the format:

`FORMAT  : [start delimiter][space]tpl:{[replacement source]}[space][end delimiter]`

`EXAMPLE :  <!--                   tpl:{ var:property.key   }        -->           `

but everything is configurable at runtime.

Using `-S` and `-E` you can set start and end placeholders' delimiters:

```sh
tpl -D \
    -s testdata/complete/template.txt \
    -d /tmp/tpl/test/out.txt \
    -S '#!--' -E '--#' \
    --var foo=hola --var config.hello=ciao

```

Using `--placeholder-separator` you can set a separator other than `:` for the block `var:property.key`.

Separator is also configurable in the template as the first non alphanumeric character in placeholder, i.e.:

```xml
<foo><!-- tpl:{|file|c:\a\path|c:\default} --></foo>
```


**Values from env vars**

Variables can be loaded from environment:

```xml
<!-- tpl:{env:THE_VAR} -->
```

**Values from user defined vars**

Variables can be declared in a properties file supplied in the option `--varfile` or
passed as arguments from command line using `--var my.property=bar`.

The `varfile` is expected in properties format:

```properties
my.property=foo
```

Input file:

```xml
<info>
  <build-hash><!-- tpl:{var:build.hash} --></build-hash>
</info>
```

Run tpl:

```sh
tpl -s template.html \
    -d out.html \
    --var build.hash=da32da4e

```

Output:

```xml
<info>
  <build-hash>da32da4e</build-hash>
</info>
```

**Values from files**

Values can be loaded from full file contents (path can be absolute or relative to the template file):

```xml
<!-- tpl:{file:/path/to/file} -->
```

**Values from file resolved at runtime**

```
#!-- tpl:{varfile:in-[ENV].txt} --#
```

```
bin/tpl -D -s testdata/varfile/template.txt -d build/varfile.txt --var ENV=PROD
```


**Indentation**

By default, replacement contents are indented, safe for formats where indentation is important, ie yaml.

Using `--no-indent` this behaviour is blocked.

**Write to standard output**

.............

**Test drive**

See in action:

```
tpl -D \
    -s testdata/complete/template.txt \
    -d /tmp/tpl/test/out.txt \
    -S '#!--' -E '--#' \
    --var foo=hola --var config.hello=ciao

tpl -s testdata/yaml/template.yaml -d output.yaml
```

## License

Apache 2.0 - see LICENSE file.

Copyright 2020 tpl contributors

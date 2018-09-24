# dp-zebedee-utils/content

Simple command line tool for generating the Zebedee-CMS directory structure and populating it with default content. Simply
provide an absolute path to a directory for the generator to create the directory structure. 

### Build the binary

```
cd content
go build -o builder
```

### Run it
```
./builder [ARGS]
```

#### Args 

| Arg        | Description                                                                   |
| ---------- |-------------------------------------------------------------------------------|
| -h / -help | Display the help menu.                                                        |
| -r         | The absolute path of the directory to generate the zebedee file structure in. |


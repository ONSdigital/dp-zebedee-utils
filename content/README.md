# dp-zebedee-utils/content

Command line tool for generating the Zebedee-CMS directory structure and populating it with default content.

Simply provide an absolute path to a directory for the generator to create the directory structure. In addition a 
personalised `/generated/run-cms.sh` is generated which can be used to run Zebedee in publishing/CMS mode using the dev
local config.

### Build the binary

```
cd content
go build -o builder
```

### Run it
```
./builder [FLAGS]
```

#### Flags 

| Flag       | Description                                                                   |
| ---------- |-------------------------------------------------------------------------------|
| -h / -help | Display the help menu.                                                        |
| -r         | The absolute path of the directory to generate the zebedee file structure in. |
| -cmd       | If `true` a CMD service account will be generated, the default is false.      |


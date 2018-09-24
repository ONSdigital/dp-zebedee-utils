# dp-zebedee-utils/content

Command line tool for generating the Zebedee-CMS directory structure and populating it with default content. Simply 
provide an absolute path to a directory for the generator to create the directory structure.

In addition a personalised `/generated/run-cms.sh` is generated which can be used to run Zebedee in publishing/CMS mode 
using the dev local config. Once generated copy file into the root of your Zebedee project.

_NOTE_: If you use the generated `run-cms.sh` you may be required to make it an executable.
```
sudo chmod +x run-cms.sh
```

### Prerequisites
- Go 1.10.2
- [Govendor][1] 

### Getting started
```
go get github.com/ONSdigital/dp-zebedee-utils
cd content
go build -o builder
```

### Run it
```
./builder -r=[YOUR_PATH]
```

#### Flags 

| Flag       | Description                                                                   |
| ---------- |-------------------------------------------------------------------------------|
| -h / -help | Display the help menu.                                                        |
| -r         | The absolute path of the directory to generate the zebedee file structure in. |
| -cmd       | If `true` a CMD service account will be generated, the default is false.      |

[1]: https://github.com/kardianos/govendor
# dp-zebedee-utils/content

Command line tool for generating the Zebedee-CMS directory structure and populating it with default content. Simply 
provide an absolute path to a directory for the generator to create the directory structure. In addition a personalised
 `/generated/run-cms.sh` is generated which can be used to run Zebedee in publishing/CMS mode using the dev local config.


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

| Flag       | Description                                                                   |
| ---------- |-------------------------------------------------------------------------------|
| -h / -help | Display the help menu.                                                        |
| -r         | The absolute path of the directory to generate the zebedee file structure in. |
| -cmd       | If `true` a CMD service account will be generated, the default is false.      |

Once the script has run successfully you will have a the Zebedee folder structure under the dir you provided for `-r`.
If you wish to use the generated `run-cms.sh` to run Zebedee CMS simply copy it to the root of your Zebedee project and 
run:
```
./run-cms.sh
``` 
_NOTE_: You may be required to make it an executable before you can run it.
```
sudo chmod +x run-cms.sh
```



[1]: https://github.com/kardianos/govendor
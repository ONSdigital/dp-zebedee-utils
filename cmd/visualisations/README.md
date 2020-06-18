# Visualisations Google Analytics Fix

Script for commenting any google analytics related code within visualisations on the ONS website.
The script copies the content to fix from master, into a new collection of the given name. 

 _Note:_ content can only be moved if that content is not already in another collection. 

## Setting up from scratch

1. SSH into the environment

```
dp ssh develop publishing_mount 1
```

2. Run a golang container with a volume that maps the content directory on the publishing box:
```
sudo docker run -i -t --name data-fix \
   --userns=host \
   -v <CONTENT_DIR>:<VOLUME_NAME.:rw \
   golang /bin/bash
```

```
sudo -i
sudo docker run -i -t --name data-fix --userns=host -v /var/florence/zebedee:/zebedee:rw golang /bin/bash
```

3. Update the container
```
apt-get update && apt-get install vim
```

4. Clone this repo:
```
cd /
git clone -b master https://github.com/ONSdigital/dp-zebedee-utils.git
```

5. Move into the directory of the binary
```
cd dp-zebedee-utils/cmd/visualisations/
```

6. Run the app
```
go build -o visualisations
export HUMAN_LOG="true"
./visualisations -zeb_root="/zebedee" -collection="visualisationsGA"
```

7. Cleanup the docker container:
```
docker rm data-fix
```

## Running the script.

### Config

| Flag       | Description                                    |
|------------|:-----------------------------------------------|
| zeb_root   | The zebedee root directory                     |
| collection | The name of the collection to use for the move |

### Example

Compile:
```
go build -o visualisations
```
Run:
```
./visualisations -zeb_root="/zebedee" \
            -collection="visualisationsGA" 
```
# Docker container tool

Docker container tool is an utility can help developer while creating and testing Docker images.

Features:
- [x] list all images managed by registry
- [x] load saved images from a folder
- [x] clean dangling images
- [x] remove ALL stopped and running container

### Usage example
To get command Usage
```sh
$ dct
```
To get images managed by Docker private registry with default port
```sh
$ dct listi registryName
```
To get images managed by Docker private registry with port
```sh
$ dct listi registryName:port
```
To load saved Docker images from current folder
```sh
$ dct loadi ./
```

### Language
golang

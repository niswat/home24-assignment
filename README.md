# Web Crawler

## Overview
This repository is an implementation of a web crawler which takes a website URL as an input and provides general information
about the contents of the page:
- HTML Version
- Page Title
- Headings count by level
- Amount of internal and external links
- Amount of inaccessible links
- If a page contains a login form

## Pre-requisites
- Internet Browser with lastest version installed , recommended Chrome/Firefox
- Docker latest version installted, recommended v20.10.7
- Make sure port 8082 is not in use as our aplpication will be listening on this port

## Quickstart

### Cloning the git repository
```
git clone git@github.com:niswat/home24-assignment.git
```
### Building the Docker Image
```
docker build -t home24assign .
```

### Running the application 
```
docker run -d -p 8082:8082 home24assign
```
### List all containers
```
 docker ps -a
```
### Accessing the application

Access the browser and enter the command : 
```
localhost:8082
```
### Crawling an URL

Add a url in the box say `https://www.example.com` and click submit button.

Wait for sometime for the browser to display the expected output.

If the application fails , access the logs to understand the failure and after fixing it , re-trigger the application using `docker run -d -p 8082:8082 home24assign`

### Accessing the logs

```
docker logs <id of running or stopped container>
```


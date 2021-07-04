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
- Internet Browser with lastest version installed , recommended Chrome
- Docker latest version installted, recommended v20.10.7
- Make sure port 8082 is not in use as our application will be listening on this port
- If using any VPN, make sure does it not blocks access to any website.

## Quickstart

### Cloning the git repository
```
git clone git@github.com:niswat/home24-assignment.git
```
### Building the Docker Image
Before running the application, the first step is to create a docker image so as to run it inside the container.
```
docker build -t home24assign .
```

### Running the application 

To start the application, execute the command
```
docker run -d --name web-scraper -p 8082:8082 home24assign
```

### Accessing the application

Launch the browser and in a new tab enter the command : 
```
localhost:8082
```
### Crawling an URL

Add a url in the box say `https://www.google.com` and click `check` button.

Wait for sometime for the browser to display the output.

### List all containers
To list all the running/exited containers.
```
 docker ps -a
```
### Stop and Start the Containers

To stop the running application, execute the command `docker stop <container id> `.  
To restart the application, execute the command `docker start <conatiner id>`

### Removing the application
To permanently delete the application, execute the command
```
docker rm -f <contianer id>
```
### Accessing the logs
```
docker logs <id of running or stopped container>
```


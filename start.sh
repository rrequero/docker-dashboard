docker build -t docker-dashboard . && docker run -ti -v $(pwd):/go/src/app -v /var/run/docker.sock:/var/run/docker.sock docker-dashboard

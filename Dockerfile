FROM golang:1.16-alpine

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
COPY *.go ./

# Copy the entire project, including the Kademlia code, into the image
COPY kademlia/*.go ./kademlia/
COPY main.go .
RUN go build -o /kadlab

# Install prerequisites
RUN apk --no-cache add curl


EXPOSE 8080
CMD ["/kadlab"]
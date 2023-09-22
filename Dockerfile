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

RUN mkdir /app

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod .

RUN go mod download

COPY ./main/main.go ./
COPY ./kademlia ./

ADD docker /app

# Build the Go app
RUN go build -o /main.go

CMD ["/main.go"]
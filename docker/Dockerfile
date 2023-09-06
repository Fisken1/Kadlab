FROM golang:alpine

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

# Copy the entire project, including the Kademlia code, into the image
COPY ../kademlia .

RUN go build -o my-kademlia-app

CMD ["./my-kademlia-app"]
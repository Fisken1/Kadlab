# Kadlab
This is the repo for labs in D7024E

# Docker
copy over files via scp:
scp -P 27070 -r Kadlab olijoh-9@130.240.207.20:/home/olijoh-9 (email for password if you need this)

# How to build the app
sudo docker build -t kadlab .

# How to run a single node
sudo docker run -it -p 8080:1550 kadlab

# How to run 50 nodes
sudo docker compose up

# How to take down 50 nodes
sudo docker compose down

# Get all running nodes 
sudo docker ps

# Run CLI with node
sudo docker attach (name of node)

# Enter node to be able to use curl (http)
sudo docker exec -it (name of node) sh

# Curl commands
curl -X POST "(ip of node or localhost):5000/objects/(What to store)" Example: curl -X POST "localhost:5000/objects/hej"
curl -X GET "(ip of node or localhost):5000/objects/(Hash)" Example: curl -X GET "localhost:5000/objects/jkljsdfhasdfjllkjsdfsdfsdfd7983543jsdfs"



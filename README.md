# Wow Such Program
docker run -d --name mysql-container -e MYSQL_ROOT_PASSWORD=1234 -p 3306:3306 mysql:8.0
docker exec -it mysql-container mysql -u root -p
docker-compose up
docker build . -t app-multistage
This program is very simple, it connects to a MySQL database based on the following env vars:
* MYSQL_HOST
* MYSQL_USER
* MYSQL_PASS
* MYSQL_PORT

And exposes itself on port 9090:
* On `/healthcheck` it returns an OK message, 
* On GET it returns all recorded rows.
* On POST it creates a new row.
* On PATCH it updates the creation date of the row with the same ID as the one specified in query parameter `id`
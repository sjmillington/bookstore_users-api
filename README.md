# bookstore_users-api
Users API for the bookstores example in GO


## Logging into mysql 

`mysql -h localhost -P3006 -u root -p`

## Setup env variables

e.g for zsh: `echo 'export mysql_users_password=password' >> ~/.zshenv `

required for db connection:

- `mysql_users_username`
- `mysql_users_password`
- `mysql_users_host`
- `mysql_users_schema`

## Posting

`curl localhost:8080/users -X POST -d @data/user.json`
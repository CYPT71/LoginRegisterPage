# User 

type 
  - webauthn user 
  - roles number corresponding of bytes associte to Roles
  - Email 
  - password backup to webauthn
  - Credentials webauthn credentials  in database and credentials are charge in code 

# Role 

type 
  - name of the role
  - byte number 

# api request 

**POST** __register/start/:username__ begin registration and send Credentials to user create a new session form user

**POST** __register/end/:username__ finish the registration and return User Credentials and save in database

**POST** __login/start/:username__ begin the login and send Credentials to user create a new session form the user if not exist

**POST** __login/end/:username__ finish the login and return User Credentials and update in database

# Run the project

go version go1.19.2
docker version 1.1.4
docker compose version v2.12.0
```sh
docker compose up -d
go mod tidy
go run .
```
New terminal or same

```sh
cd front
npm i 
npm run serve
```

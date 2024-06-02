# GreenBasket API
 
GreenBasket is an E-commerce application Backend Service specifically for farmers to sell their crops directly to customers. 
Products sold are fresh vegetables, fruits, eggs, meat, fish and various kitchen spices. 
This backend service is built using a monolithic architecture, all service components are integrated in one cohesive codebase, making it easy to develop and maintain and deployment is faster and simpler because it involves only one artifact or build that needs to be managed.

This API is open source, so you can use it for the backend needs of your E-commerce project that you have customized.

# API Installation


## Required
To run the API on your local machine, make sure these items are installed on your computer:
* Ensure you have Golang installed. Refer to the official Golang website for installation instructions: https://go.dev/doc/install.
* Ensure you have Docker installed. Refer to the official Docker website for installation instructions: https://docs.docker.com/desktop/.
* Ensure you have MongoDB installed. Refer to the official MongoDB website for installation instructions: https://www.mongodb.com/docs/manual/installation/.
Or You can use MongoDB Atlas instead.

The other items that You need to setup are:
* Create App Password using Your Google Account. You can follow the step on this blog: https://knowledge.workspace.google.com/kb/how-to-create-app-passwords-000009237
* Create Oauth2 credential use Google API Client Libraries or Google OAuth 2.0 endpoints to implement OAuth 2.0 authorization to access Google APIs.
You can follow the step on this docs: https://developers.google.com/identity/protocols/oauth2/web-server
* Create Oauth2 credential use Facebook Developer.
You can follow the step on this docs: https://baserow.io/user-docs/configure-facebook-for-oauth-2-sso
* Create Your Midtrans account for the payment gateway that using by this project. Note: use the sandbox version. You see the documentation [here](https://docs.midtrans.com/docs/midtrans-account)

## Steps
To run the API on your local machine, make sure You follow these steps:

Clone the Repository
```bash
git clone https://github.com/IndraSty/GreenBasket.git

cd GreenBasket
```

Initialize Go module and download dependencies:
```bash
go mod init GreenBasket
go mod download
```

Start up a docker container running MongoDB and Redis. A `docker-compose` file is provided to make this easier.Make sure you are in the same root folder as `docker-compose` is located.
```bash
docker compose up
```

Create an `.env` file and fill it according to the `.env.example` 
file which has been completed according to the required items explained previously.

Import data from the `/data` folder to your database. We recommend using a GUI like MongoDB Compas to make this easier.

Run the API on Your local machine
```bash
go run cmd/server/main.go
```

To build the source code running
```bash
go build
```

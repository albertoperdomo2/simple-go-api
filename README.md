# Golang API
## A simple Golang API for learning purposes. 

With this project, what I am trying to learn is how API's works and how to efficiently design (a very basic) one. It is an extremely simple API with just a few endpoints, which I'll take as a base point to develop more complex projects. 

At the moment, this API does not have any type of DB interaction, but next step in the process would be to add PostgreSQL as backend service for this API. So, it is initialized empty. 

It should be noted that this API is **intended to work with image links** in the future. 

The **functional requisites** that this API meets are:
```
* GET /images return a JSON listing all the images present. 
* GET /images/random return a JSON with the information of a random-picked image. 
* GET /admin requires basic authentication. 
* GET /images/{id} return a JSON with the information of image {id}. 
* POST /images allows image information in JSON format to be inserted, and return an error if type of content is not application/json. 
```

To run this sever in localhost:8080, you will need to declare the *ADMIN_PASSWORD* environment variable, so e.g.:
```bash 
ADMIN_PASSWORD=secret go run api.go
```

You can test out the API using the following example command lines: 
POST new content: 
```bash
curl localhost:8080/images -X POST -d '{"fileName": "exampleFileName", "author": "exampleAuthor", "size": 12.3}' -H "Content-type: application/json"
```
AUTHENTICATE: 
```bash 
curl localhost:8080/admin -u admin:secret
```

As mentioned previously, this is a first step in what I would consider a more complex project and it is intended to help me along the learning process, so mistakes and errors along the way are expected. 


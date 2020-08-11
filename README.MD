##### Usage
###### Server
Cart service listening to port 8000.

To run locally: `make && bin/shopping-cart -server` 

To run in scalability mode (**Docker needed**): `make docker compose SCALE=5` just set SCALE to the number of services instances you whish.

To run unit and functional tests: `make test`

###### Client
This is a console client connecting to http://localhost:8000. This can be overridden via command line argument "server".

To run locally: `make && bin/shopping-cart -client` 

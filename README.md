# Service Preview Demo

This is a simple demo application with four services to demonstrate the Ambassador
Service Preview feature.  The application is a simple aircraft inventory, displaying
a list of available aircraft and details on each.  The asking price of each aircraft is
in US Dollars but may be converted to other currencies such as Euros or Norwegian Kroner.

There are four services:
- an Application service (appservice), which renders the web pages for the application;
- an Image service (imageservice), which, given a registration number, returns an image of the aircraft;
- an Inventory service (inventoryservice), which returns a list of all available aircraft;
- and a Specs service (specsservice), which, given an aircraft model, returns that aircraft's specifications

Each service is a Kubernetes Service with a Deployment, and routing to each Service is
done through an Ambassador Mapping Resource.  Communication between Services is via JSON.




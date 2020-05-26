# Service Preview Demo

This is a simple demo application with four services to demonstrate the Ambassador
Service Preview feature.  The application is a simple aircraft inventory, displaying
a list of available aircraft and details on each.  The asking price of each aircraft is
in US Dollars but may be converted to other currencies such as Euros or Norwegian Kroner.

There are four services:
- an Application service (appservice, written in Go), which renders the web pages for the application;
- an Image service (imageservice, written in Go), which, given a registration number, returns an image of the aircraft;
- an Inventory service (inventoryservice, written in Python), which returns a list of all available aircraft;
- and a Specs service (specsservice, written in Python), which, given an aircraft model, returns that aircraft's specifications

Each service is a Kubernetes Service with a Deployment, and routing to each Service is
done through an Ambassador Mapping Resource.  Communication between Services is via JSON.

## Application Flow

The appservice is responsible for communicating with the other services and rendering the pages using the golang
template system.  The pages are:

- ```home.html```, the landing page for the site;
- ```results.html```, the inventory list with a thumbnail of each aircraft and some basic information;
- ```detail.html```, a detail page on a single aircraft, with a photo and specifications for that aircraft model.

Rendering the Home page simply lays out "Welcome to the Aircraft Bazaar", and requests a static image that is
served from the appservice itself using the http.FileServer module.  A link at the bottom navigates to the
```results.html``` page.

## Results Page

The ```results.html``` page, served up by the appservice, requires two other services to render its content: the
```inventoryservice``` and the ```imageservice```.

appservice [requests inventory from] ==> ```inventoryservice``` [returns list of aircraft]


## How to run the Service Preview demo


`edgectl install`                         # Install the Ambassador Edge Stack
`edgectl license <license key here>`      # Apply the license, required for edgectl intercept.

`make all`                                # Create the Docker images for all services and push to Docker Hub
`make deploy`                             # Apply the service YAML files and the traffic agent RBAC

`kubectl get pods`                        # See all the running pods
`kubectl describe pod appservice...`      # The appservice should run with the traffic-agent sidecar
`kubectl describe deployment appservice...`

`edgectl daemon`

`edgectl connect`
`edgectl status`
`edgectl intercept available`

Now, launch your local services so that when they are intercepted by AES and sent to your localhost, they
will provide the services that are currently running in the cloud.
`edgectl intercept add appservice -t localhost:8080`
`edgectl intercept add specsservice -t localhost:8081`
`edgectl intercept add inventoryservice -t localhost:8082`


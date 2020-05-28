# Service Preview Demo

This is a simple demo application with four services to demonstrate the Ambassador
Service Preview feature, which is documented [here](https://www.getambassador.io/docs/latest/topics/using/edgectl/)

The application is a simple aircraft inventory, displaying
a list of available aircraft and details on each.  The asking price of each aircraft is
in US Dollars but may be converted to other currencies such as Euros or Norwegian Kroner.

There are four services:
- an Application service (`appservice`, written in Go), which renders the web pages for the application;
- an Image service (`imageservice`, written in Go), which, given a registration number, returns an image of the aircraft;
- an Inventory service (`inventoryservice`, written in Python), which returns a list of all available aircraft;
- and a Specs service (`specsservice`, written in Python), which, given an aircraft model, returns that aircraft's specifications

Each service is a Kubernetes Service with a Deployment, and routing to each Service is
done through an Ambassador Mapping Resource.  Communication between Services is via JSON.

The Image, Inventory and Specs services provide a database-style lookup of images, individual aircraft information,
and specifications for each model of aircraft in the inventory.  To make the demo source as simple as possible,
this information is stored directly in each service in a `data` directory.  In the case of the Image service, the
individual aircraft images are stored as `.jpg` files where each file's name matches that aircraft's registration.
In the Inventory and Specs services, the information is stored in individual `json` files, `data/inventory.json` and
`data/specs.json` respectively. These files are read in at startup time and the values returned as `json` when the service
handles a request.

## Application Flow

The appservice is responsible for communicating with the other services and rendering the pages using the golang
template system.  The pages are:

- `home.html`, the landing page for the site;
- `results.html`, the inventory list with a thumbnail of each aircraft and some basic information;
- `detail.html`, a detail page on a single aircraft, with a photo and specifications for that aircraft model.

Rendering the Home page simply lays out "Welcome to the Aircraft Bazaar", and requests a static image that is
served from the `appservice` itself using the `http.FileServer` module.  A link at the bottom navigates to the
`results.html` page.

## Inventory Page

The `inventory.html` page, served up by the `appservice`, requires two other services to render its content: the
`inventoryservice` and the `imageservice`.

## Detail Page

The `detail.html` page, served up by the `appservice`, requires two other services to render its content: the
`specsservice` and the `imageservice`.

## How to run the Service Preview demo

### Customizing to your environment

The toplevel `Makefile` (`service-preview-demo/Makefile`) has a number of functions:

- `make build-images` builds the images for each of the services

- `make push-images` pushes these images to your desired Docker repository.

- `make deploy` deploys the services, and applies the `traffic-agent-rbac.yaml` to your cluster.

- `make all` builds and pushes all the images to the repository but does not deploy the services.

You may also build and push individual service images:
- `make app` builds and pushes the Application Docker image (code, `html` and `css` files)
- `make inventory` builds and pushes the Inventory Docker image (code and related `json` data)
- `make specs` builds and pushes the Specs Docker image (code and related `json` data)
- `make images` builds and pushes the Image service Docker image (code and `jpeg` files)

These are used when you have modified one of the services and want to run those changes in the cluster.
You'll need to `kubectl delete pod <pod-id>` for the pods running those specific services to restart each service
with the code and data changes.  Of course, with Service Preview, this will be unneccessary most of the time since
you will be running the services locally and can modify and re-run any service to see your changes much more
quickly than pushing a large image to a repository and restarting a pod.

Two environment variables specify your Docker registry and project name, and can be customized to your needs:

- `$DEV_REGISTRY` defaults to `docker.io/brucehorn` but is overridden if defined in the environment.

- `$PROJECT_NAME` defaults to `service_preview`, and can also be overridden if otherwise defined.

### Installing AES

Install the Ambassador Edge Stack:

`edgectl install`

Apply the license, required for `edgectl intercept`.

`edgectl license <license key here>`

### Building and Deploying

Create the Docker images for all services, and push the images to Docker Hub.  You will want to 

`make all`

Apply the service YAML files and the traffic agent RBAC

`make deploy`

See all the running pods

`kubectl get pods`

The appservice should run with the traffic-agent sidecar

`kubectl describe pod <appservice pod UID here>`

`kubectl describe deployment <appservice pod UID here>`

Start the edgectl daemon, needed for the client side to connect to the cluster.

`edgectl daemon`

Connect to the cluster and check its status.  The proxy should be ON.

`edgectl connect`
`edgectl status`

List the available intercepts.  These are the deployed services.

`edgectl intercept available`

Now, launch your local services so that when they are intercepted by AES and sent to your localhost, they
will provide the services that are currently running in the cloud.

`edgectl intercept add appservice -t localhost:8080`

`edgectl intercept add specsservice -t localhost:8081`

`edgectl intercept add inventoryservice -t localhost:8082`

## Some easy modifications to see the Service Preview working


# WIG - Server

## Table of Contents
1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
   - [Prerequisites](#prerequisites)
   - [Initial Setup](#initial-setup)
3. [API Documentation](#api-documentation)
   - [User Controller](#user-controller)
   - [Borrower Controller](#borrower-controller)
   - [Location Controller](#location-controller)
   - [Ownership Controller](#ownership-controller)
   - [Scanner Controller](#scanner-controller)

---

## Introduction <a name="introduction"></a>

WIG is a self-hosted inventory management service to remember "What I Got".<br>
It is built for the everyday person in need of organization.<br>
WIG-Server is the backend database manager that handles all functionality behind the application.

## Getting Started <a name="getting-started"></a>
The WIG-Server is only needed if the plan is to Self-Host.
A free tier is available for all with access to the [WIG-Android](https://github.com/Narrow-Code/WIG-Android) application.

### Prerequisites <a name="prerequisites"></a>
To run the WIG-Server, Docker is required as all components are containerized. Please ensure Docker is installed on your system before proceeding with the installation. If Docker is not yet installed, you can install it by following [Docker Installation](https://docs.docker.com/engine/install/)

### Initial Setup <a name="initial-setup"></a>
- To setup the WIG-Server clone the repository to the desired directory in your server.
- In the wIG-Server directory, rename .envTemplate to .env
- Make any personal changes to .env
- Any changes made to .env must be changed in docker-compose.yml
- From the WIG-Server directory, run the following command:

```bash
docker stack deploy —compose-file docker-compose.yml wig
```
- After running this command, WIG-Server should be containerized and running on Port 30001
- Follow instructions on [WIG-Android](https://github.com/Narrow-Code/WIG-Android) to setup Self Hosting for the Android Application.

## API Documentation <a name="api-documentation"></a>
Bellow is the API documentation for understanding.

### User Controller <a name="user-controller"></a>
The user controller API calls interact with the User table.

| URI | Supported Methods | Description |
|-----|-------------------|-------------|
| user/signup | POST | Creates a user |
| user/{username}/salt | GET | Returns the Users salt |
| user/login | POST | Logs the user into the application creating a token |
| user/validate | GET | Validates that the Users token is still active |

### Borrower Controller <a name="borrower-controller"></a>
The borrower controller API calls interact with the Borrower table.

| URI | Supported Methods | Description |
|-----|-------------------|-------------|
| app/borrower | POST, GET, DELETE | POST = Create new Borrower<br>GET = Get all Borrowers<br>DELETE = Delete borrower |
| app/borrower/{borrowerUID}/checkout | PUT | Checks out all included Ownerships from the body to the Borrower in the URI |
| app/borrower/check-in | POST | Checks in all included Ownerships from the body |
| app/borrower/checked-out | GET | Returns all checked-out Ownerships and their coorosponding Borrowers|

### Location Controller <a name="location-controller"></a>
The location controller API calls interact with the Location table.

| URI | Supported Methods | Description |
|-----|-------------------|-------------|
| app/location | POST, PUT, DELETE | POST = Create new Location<br>PUT = Edit location fields<br>DELETE = Delete location |
| app/location/{locationUID}/set-parent | PUT | Edits the parent location of the location | 
| app/location/{locationUID} | GET | Returns all Ownerships associated with the Location in the URI |
| app/location/search | POST | Searches inventory and returns matches for specified Location |
| app/inventory | GET | Returns all Locations and Ownerships associated with the user |

### Ownership Controller <a name="ownership-controller"></a>
The ownership controller API calls interact with the Ownership table.

| URI | Supported Methods | Description |
|-----|-------------------|-------------|
| app/ownership | POST, PUT, DELETE | POST = Create new Ownership<br>PUT = Edit Ownership fields<br>DELETE = Delete specified Ownership |
| app/ownership/{ownershipUID}/quantity/{type} | PUT | Changes the quantity based on the types, as follows:<br><br>increment = Increases quantity by 1<br>decrement = Decreases quantity by 1<br>set = Sets the quantity to the desired amount |
| app/ownership/{ownershipUID}/set-parent | PUT | Sets the Location associated with the Ownership |
| app/ownership/search | POST | Searches inventory and returns matches for specified Ownership |

### Scanner Controller <a name="scanner-controller"></a>
The scanner controller API calls interact with sent barcode and QR codes.

| URI | Supported Methods | Description |
|-----|-------------------|-------------|
| app/scan/{barcode} | POST | Searches for existing barcode in the database, if one exists, it returns the Ownership associated. If Ownership does not exist, searches for existing item in the database and creates Ownership based on the item returned |
| app/scan/check | GET | Checks a QR code to see if it is associated with a location or ownership or if it is an unused QR. |
| app/scan/location | GET | Returns the Location associated with the barcode or QR code |
| app/scan/ownership | GET | Returns the Ownership assocaited with the barcode or QR code |

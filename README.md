#  MyAppBooth: Lightweight Web Application Search and Purchase Platform

### *Introduction*

This page documents the software design details for the `MyAppBooth` web application which mainly supports user authentication, App item search and payment checkout features

---

### *Developer Environment*

* Frontend
  * Framework: React
  * Programming Language: JSX (converted to client browser compatible JS version via Babel under the hood)
  * Backend: node.js
  
* Backend
  * Developer Environment: Google Compute Engine (GCE) virtual machine
  * Programming Language: Go
  * CORS Policy: enabled as ok

---

### *Frontend Software Architecture*

![Frontend Design](images/frontend_architecture.PNG)

Figure 1. Frontend software architecture

---
### *Frontend Software Components*

#### 1. App

This is the main application component
* Layout
  * Header
    * Display user menu if user is logged in
  * Content
    * Display LoginForm if user is not logged in
    * Display HomePage if user is logged in
 
#### 2. LoginForm

It supports user login function
  * Form
    * FormItem: username
    * FormItem: password
    * FormItem: login button
  * SignupForm
  
#### 3. SignupForm

It supports user signup function
  * Form
    * Form Item: username
    * Form Item: password
    * Form Item: signup button

#### 4. HomePage

It supports tab switching functions
* TabPane
  * BrowseApps: search App item by title or description and support App item payment checkout function
  * PostApps: upload App item by current user
  * MyApps: display uploaded App item from current user and support App item deleting function

#### 5. BrowseApps

It supports App item search and payment checkout functions
* Form
* ListItem
  * Card: display App item information
  * Button: payment checkout App item
  * Image: display App item
  
#### 6. PostApps

It supports App item uploading function
* Form
  * FormItem: App item description
  * FormItem: App item title
  * FormItem: App item pictue
  * FormItem: App item price
  * FormItem: submit App item uploading button
  
#### 7. MyApps

It displays App uploaded from current user and supports App deletion function
* ListItem
  * Card: display App item information
  * Button: delete App item
  * Image: display App item

---

### *Backend Software Architecture*


![Backend Design](images/backend_architecture_flow.PNG)

 Figure 2. Backend software architecture


---
### *Backend Software Components (Handler)*

The handler responses for handling specific client request redirected from HTTP router. There is a JWT middleware layer added on top of the handler to verify user
authentication before handling the client request. The user authenication is JWT token based. A generated token will be sent to client once the first time user login is verified. The client request carries the token afterward for verifying the user information. The CORS request is enabled to unblocks cross domain query.

#### 1. UploadHandler

* Parse the HTTP request into App item object
* Call the `SaveApp` from the service layer for futher App uploading actions 

#### 2. SignupHandler
* Parse the HTTP request into user object
* Call the `AddUser` from the service layer for further user registration actions

#### 3. LoginHandler
* Parse the HTTP request into user object
* Call the `CheckUser` from service layer for further user authentication actions

#### 4. SearchHandler
* Parse the HTTP request into title, description and username query string 
* Call the `SearchApp` from service layer for further App searching actions

#### 5. DeleteHandler
* Parse the HTTP request into App item id query string 
* Call the `DeleteApp` from service layer for further App deletion actions

#### 6. CheckoutHandler
* Parse the HTTP request into App item id query string 
* Call the `Checkout` from service layer for further App payment checkout actions

---

### *Backend Software Components (Service)*

#### 1. SaveApp

* Call `CreateAppWithPrice` from backend client layer to create App product id and price id via Stripe API
* Call `SaveToGCS` from backend client layer to save media file for the APP item and create the media file link
* Call `SaveToES` from baclend client layer to save App item metadata

#### 2. AddUser
* Verify the user is not existed
* Call `SaveToES` from baclend client layer to save user information metadata


#### 3. CheckUser
* Call `ReadFromES` from baclend client layer to retrieve user information
* Verify the user password

#### 4. SearchApp
* Provide default search method
* Provide search by username only method
* Provide search by title only method
* Provide search by description only method
* Call `ReadFromES` from baclend client layer to search App items


#### 5. DeleteApp
* It allows to delete only the App items uploaded by the current user 
* Call `DeleteFromES` from baclend client layer to delete App item

#### 6. Checkout
* Call `SearchApp` from current layer to retrieve the App item price id 
* Call `CreateCheckoutSession` from backend client layer to retreive the payment checkout link

---

### *Backend Software Components (Backend Client)*

#### 1. ReadFromES
* Use ElasticSearch backend client object to read data based on input query string and index

#### 2. SaveToES
* Use ElasticSearch backend client object to save data based on input query string and index

#### 3. DeleteFromES
* Use ElasticSearch backend client object to delete data based on input query string and index

#### 4. SaveToGCS
* Use Google Cloud Storage backend object client to save media data

#### 5. CreateAppWithPrice
* Call Stripe API to create App item product id and price id

#### 6. CreateCheckoutSession
* Call Stripe API to retrieve App item payment checkout link based on App item price id


---
### *Database*

#### 1. ElasticSearch (ES)
  * NoSQL database  
  * Store user information and App item metadata (app basic information, media file url, product id and price id for payment checkout link query)
  * Create inverted index for App item title and description
  * Support fast keyword search from user input
  
 The following figures show the schema for App item and user information in ElasticSearch (ES):
 
 ![Backend Design](images/app_item_schema.PNG)
 
 Figure 3. App item schema
 
 
 ![Backend Design](images/user_information_schema.PNG)
 
 Figure 4. User information schema
 
  
#### 2. Google Cloud Storage (GCS)
  * Blob storage for non-structured data 
  * Store media files for App item
  * Link of each media file is stored as metadata in ElasticSearch (ES)
  
---
 
### *Deployment*


* Frontend
  * Amazon Web Service (AWS) based
  * Deploy frontend build package to AWS Simplify (PaaS)  

* Backend
  * Google Cloud Platform (GCP) based
  * Docker containerization technology
  * Create Dockerfile to include deployment configuration
  * Deploy the backend code package to Google App Engine [GAE] (PaaS)

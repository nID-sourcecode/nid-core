# Wallet

The purpose of the wallet is to give citizens insight and control over personal data. The user should register (in de future with DigiD), whereafter they can give or deny consent for sharing data. Within the wallet exists an overview
of organisations whom consent was given to, so have the ability to retrieve their data. Furthermore, the wallet provides
insight of which data is actually provided to which organisations.

The wallet is implemented as application and can be installed on mobile devices. Secrets are being stored in the devices
secure storage after the registration process, which is comparable with mobile banking applications.

## Backend

The wallet backend is consists of two services, RPC and GQL.

#### RPC

The RPC service is responsible for registering new users and signing in existing users. Within this service are auth
functionalities which can be used for signing in and registering a device. Lastly there is an endpoint for creating consent,
which is being called after the user accepts the request from a client to share their data.

#### GQL

The GQL service is responsible for exposing CRUD functionalities to the application about the citizens data.

## Frontend

The frontend provides the user with an interface where various tasks can be done, in the form of a mobile application.
It provides multiple functionalities for the user to supply insight and control over the given consents. When an auth request
is sent to the user they are provided with a QR code. This code can be scanned with the wallet app, after the user is authenticated].
The user can give consent to the client by accepting this request, where optional data can be toggled off, or denying the request.

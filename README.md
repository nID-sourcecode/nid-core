# nID Framework
## nID Framework
The nID Framework is a reference architecture and collection of (technical) open source
components enabling identification, authentication, authorization, access, data protection
and security in data sharing and disclosure across organizations. It offers also citizens control
over their own data by disclosing personal attributes to third parties via an authorization
mechanism. Data consumers ( affiliated parties) can view and use data within the scope of
their mandate (consent). Parties participating in the nID Framework can make mutual
agreements about data sharing and the consultation of data takes place by means of
consents granted.
The nID Framework can be implemented in federations (distributed). These federationscan
also exchange data with each other. Switching from one framework to another should be
frictionless as long as the same open standards and protocols are used.
The nID Framework was executed under project name nID System (Dutch: nID Stelsel). This
way different modules and documentation refer to nID.
## Philosophy
nID started out of the need to break the dilemma surrounding data and privacy.
The dilemma concerns organizations that need broad use of data to improve services or
solve problems. While doing so they may compromise the privacy of individuals.
nID aims to facilitate parties in exchanging data through its framework in a secure and
privacy-friendly manner (trusted network). The aim is to ensure that the large amount of
data that is available can be used for better services and solving social problems while
enabling individuals to manage their own privacy. This type of relationship is based on trust.
By means of nID, individuals gain a position in the data eco-system with the possibility to
break the relationship to their personal data if that trust is broken. Data exchange and
revocations are always within the scope and margin with a legal base (consent, agreement,
etc.).

## How to use the nID Framework
The nID Framework consists of multiple components that together make up the entire
Framework. The development is based on several open source components that interact
with each other and on which additional functionality has been developed. The following
open source components/standaards have been used to develop the nID Framework:
- GraphQL;
- Istio;
- Kubernetes;
- Json-LD;
- Lua;
- JWT;
- OpenID Connect;

- OAuth 2.0.
The main component is nID Core, a trusted network in which trusted parties can onboard
the network and based on a consent-authorization structure are able to share data. Some
components run as multi-party shared services, others run on-premise at organizations. All
components are maintained in a single repository. This means that a developer has all the
tools and code to build and test the complete framework. It simplifies version and
dependency management and allows changes that affect multiple components to be
combined in a single feature branch and merge-request.
If you want to develop locally, or run your own nID Framework you will likely want to start all
the components first.
## Troubleshooting
If you are experiencing problems or running into other issues, please visit the
troubleshooting page.
## Features
- **Accessible:** easy to implement and low-threshold requirements
- **Transparent:** authorization consents are transparent and readable
- **Privacy by design:** improving awareness among parties
- **Federated:** enabling data eco-systems to interact with each other
- **Open source and international standards:** community driven/interoperable
- **Security by design:** &#39;Zero trust&#39; model to enforce correct actions
## Licence
[Licensed under the
EUPL](https://github.com/SecureDataSharingFramework/SDSF/blob/master/LICENCE.md)
&gt; Written with [StackEdit](https://stackedit.io/). 

# The authorization component

This document explains the functionality of the authorization component and how it is implemented within the TWI platform.

## Actors within an nID-system

Seeing as nID is an authorization platform, actors in the system are named after their primary role in the authorization system.

### Client

A client _consumes_ data.

### Audience

An audience _provides_ data. From an authorization perspective, this is called an _audience_ because it is receives an access token and is expected to verify it.

### End-user

An end-user is a subject of data, equivalent to _burger_.

### Platform

The _platform_ contains the nID components and facilitates authorzation.

## Kinds of authorization (summary)

### End-user authorization

This scenario involves an _end-user_ granting access through an openid-flow, using the wallet app. See **end-user authorization** for more details.

This scenario is currently _supported_ by the TWI platform.

### Automatic data-based consent

This scenario involves granting consent based on data itself. For example, there may be a law which grants a public health office access to a municipality's data about a certain person if that person has been assigned to said public health office. A consent-by-law script may automatically pick this up and grant consent to the PHO in question.

This scenario is currently **not supported** by the TWI platform.

### Automatic role-based consent

Some clients may automatically be granted access to certain audiences because they have a certain role within the system. Information processing services, for example, fall into this category. Such services may have full access to certain data sources because they are trusted to handle this information with care.

This scenario is currently _supported_ by the TWI platform.

### Current state of TWI authorization

| Kind of authorization        | Supported          |
| ---------------------------- | ------------------ |
| End-user consent             | :heavy_check_mark: |
| Automatic data-based consent | :x:                |
| Automatic role-based consent | :heavy_check_mark: |

## Automatic token verification

Generally, in oAuth systems, it is the responsibility of the audience to verify a token and its claims. Using nID removes this responsibility from an audience by passing its inbound requests through a filter that automatically verifies the token.

Since the query and access token are both present in the request, the filter can simply verify that the query matches the scopes that are defined in the access token.

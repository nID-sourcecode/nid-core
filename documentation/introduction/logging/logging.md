Logging binnen het platform wordt onderverdeeld in 2 onderwerpen. Technische logging is logging uit alle services die draaien binnen het platform. Anderzijds is er audit logging. Dit zijn logs over kritieke gebeurtenissen binnen het platform.

## Audit logging

Op het moment van schrijven is er een conceptversie voor audit logging opgezet. Hiervoor is een envoy filter toegevoegd aan alle informatie services en databronnen die aangesloten zitten op het platform. Dit filter logt alle inkomende requests en de response die terugkomt op deze requests. Op deze manier wordt een volledig inzicht verschaft over alle requests die door het platform heen gaan. Deze logs worden momenteel verzameld door een service en als technische logs gelogd. Dit wordt momenteel verzameld door [ELK](https://www.elastic.co/what-is/elk-stack). De definitieve vorm waarin audit logs verzameld moeten worden moet nog worden bepaald. Weave heeft alleen controle over wat er geaudit logt wordt binnen functionaliteit van services die onderdeel zijn van de nID oplossing en het netwerkverkeer.

## Technische logging

Elke service in het platform is verantwoordelijk voor zijn eigen logging. Momenteel wordt alle technische logging verzameld door de ELK stack.

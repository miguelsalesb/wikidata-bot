# wikidata-bot
This script is not yet complete and was tested using a local Wikidata installation that runs in Docker. It retrieves data from the National Library of Portugal bibliographic and authorities MarcXchange repositories, exports it to Wikidata and to a MySQL database.

Initially, the following data shoud be added to Wikidata:
P1 = date of birth (P569)
P2 = date of death (P570)
P3 afirmado em - stated in (P248)
P4 endereço eletrónico da referência - reference URL (P854)
P5 instância de (P5) - instance of (P31)
P6 identificador PTBNP - Portuguese National Library ID (P1005)
P7 data de acesso (P7) - retrieved (P813)
P8 obra destacada - notable work (P800)
P9 data de publicação - publication date (P577)
P10 país de origem - country of origin (P495)
Q1 ser humano (Q1) - human (Q5)
Q2 BNP - National Library of Portugal (Q245966)
Q3 obra escrita - written work (Q47461344)
Q4 Portugal - Portugal (Q45)

For the moment, it creates new Portuguese author's entities with:
label in Portuguese and English
description in Portuguese
Portuguese aliases
instance of human
date of birth
date of death

Title entities:
labels in Portuguese and if the work's original language is either English, French or Spanish, it also adds that information
notable work: written work
country of origin: Portugal
publication date

For each property, the following references are created:
stated in: BNP - National Library of Portugal
reference URL
retrieved

The script calculates the probability of an author already existing in Wikidata, and if there is none, it creates a new author entity. It also checks which author's occupations are registered in the author's repository and not in Wikidata and registers that information in the MySQL database, so that it can in a second moment, export the lacking occupations to Wikidata (this part is not yet developed).


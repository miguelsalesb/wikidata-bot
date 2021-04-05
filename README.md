# wikidata-bot
This script is not yet complete and was tested using a local Wikidata installation that runs in Docker. It retrieves data from the National Library of Portugal bibliographic and authorities MarcXchange repositories, exports it to Wikidata and to a MySQL database.
<br><br>
Initially, the following data shoud be added to Wikidata:<br>
P1 = date of birth (P569)<br>
P2 = date of death (P570)<br>
P3 afirmado em - stated in (P248)<br>
P4 endereço eletrónico da referência - reference URL (P854)<br>
P5 instância de (P5) - instance of (P31)<br>
P6 identificador PTBNP - Portuguese National Library ID (P1005)<br>
P7 data de acesso (P7) - retrieved (P813)<br>
P8 obra destacada - notable work (P800)<br>
P9 data de publicação - publication date (P577)<br>
P10 país de origem - country of origin (P495)<br>
Q1 ser humano (Q1) - human (Q5)<br>
Q2 BNP - National Library of Portugal (Q245966)<br>
Q3 obra escrita - written work (Q47461344)<br>
Q4 Portugal - Portugal (Q45)<br>
<br><br>
For the moment, it creates new Portuguese author's entities with:<br>
label in Portuguese and English<br>
description in Portuguese<br>
Portuguese aliases<br>
instance of human<br>
date of birth<br>
date of death<br>
<br><br>
Title entities:<br>
labels in Portuguese and if the work's original language is either English, French or Spanish, it also adds that information<br>
notable work: written work<br>
country of origin: Portugal<br>
publication date<br>
<br><br>
For each property, the following references are created:<br>
stated in: BNP - National Library of Portugal<br>
reference URL<br>
retrieved<br>
<br><br>
The script calculates the probability of an author already existing in Wikidata, and if there is none, it creates a new author entity. It also checks which author's occupations are registered in the author's repository and not in Wikidata and registers that information in the MySQL database, so that it can in a second moment, export the lacking occupations to Wikidata (this part is not yet developed).


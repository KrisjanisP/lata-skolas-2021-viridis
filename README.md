![logo](web/assets/logo.png)

Krišjānis Petručeņa / komanda "One Lone Coder" 

Atvērto ģeotelpisko datu hakatons skolēniem 2021

2021.10.07 - 2021.11.04

![example](./example.png)
# Info par hakatonu

[Atvērto ģeotelpisko datu hakatons skolēniem 2021](https://www.lata.org.lv/skolas-2021)

## Mērķis

> Veidot sabiedrības, tostarp jauniešu, izpratni par ģeotelpisko atvērto datu nozīmīgumu inovatīvu produktu un pakalpojumu izstrādē un veicināt ģeotelpisko datu atvēršanu Latvijā, kā arī izglītot jauniešus par atvērto datu tehnoloģiju izmantošanu.

## Uzdevums

> Radoša, inovatīva un praktiska brīvo un atvērto datu izmantošana, radot sabiedrībai noderīgu produktu vai pakalpojumu, iesaistot tehnoloģijas.

## Rezultāts

> Hakatona dalībnieku sasniedzamais rezultāts ir radīt inovatīvu lietotni, kas risina sabiedrībai būtisku problēmu, izmantojot atvērtos ģeotelpiskos datus.

## Lokālā projekta startēšana

Pirmkārt, vēlams lai uz datora būtu:
* Uzstādīta Go jeb golang programmēšanas valoda
* Uzstādīta Python3 programmēšanas valoda
* Uzstādīts sqlite3 CLI
* Uzstādīts npm jeb node package manager
* Pieejami vismaz papildus 8 GB brīvpiekļuves atmiņas

Datubāzes sākotnējā izveide:
```
go run scripts/database/gen-database.go
```
Datubāzes sākotnējā aizpilde:
```
go run scripts/tiles/fetch-tile-names.go
go run scripts/links/fetch-tile-links.go
```
TKS93 kartes sadalījuma ģenerēšana:
```
go run scripts/geotiff/gen-geotiff.go
```
Pierakstīšanās sistēmas pieslēgšana:
```
jāizveido .env fails, kurā jānorāda sekojošās vērtības:
AUTH0_DOMAIN
AUTH0_CLIENT_ID
AUTH0_CLIENT_SECRET
AUTH0_CALLBACK_URL
, kuras var iegūt piereģistrējoties https://auth0.com/
```
Kad visi iepriekšējie soļi izpildīti,
projekta palaišana uz porta 8080:
```
go run .
```
Lai kaut ko kodēt mājaslapas javascript:
```
cd web/src && npm install
npm run dev
```
## Vērtēšanas kritēriji

1. Atvērto ģeotelpisko datu inovatīvs pielietojums.
2. Tehniskais izpildījums
3. Idejas dzīvotspēja
4. Projekta progress
5. Projekta prezentācija

## Pieejamo datu kopas

1. Latvijas Ģeotelpiskās informācijas aģentūras datu kopas
2. OpenStreetMap datu kopa
3. LVRTC datu kopa
4. Valsts Zemes dienests datu kopas
5. Atvērtie dati Latvijā
6. Atvērtie dati Eiropas Savienībā


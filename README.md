Social media summarizer
===
[![wercker status](https://app.wercker.com/status/47e29fd42244d65e3e738ccb7c6bc3e2/m/master "wercker status")](https://app.wercker.com/project/bykey/47e29fd42244d65e3e738ccb7c6bc3e2)
Project to create a summarizer for:
* Twitter (initially)
* Reddit (possibly)

Architecture
---

There are 3 separate go routines that communicate through channels:

* data.go - getting the data from twitter through Steaming API
* process.go - processing of data (i.e sorting the results)
* main.go - Web & API endpoints using [Gin](https://github.com/gin-gonic/gin)

Right now all the data is in memory, and during redeployment the app saves it in dump_tweets, dump_links.

Deployement
---

App is hosted at DigitalOcean at [162.243.99.74](https://162.243.99.74:8080). The deployment is done automatically by Wercker after successful CI.

Possible next steps
---
- [ ] Better frontend using ReactJS
- [ ] Storing information in some database
- [ ] More advanced processing of tweets


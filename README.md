# PublicationService
CHESS Publication service. Its implementation is based on
[Zenodo](https://zenodo.org) APIs, see his [blog](https://felipecrp.com/2021/01/01/uploading-to-zenodo-through-api.html)

### APIs
```
# list docs
curl http://localhost:8355/docs

# create new DOI resource
curl http://localhost:8355/create -v -X POST -H "Content-type: application/json" -H "Authorization: bearer $t" -d '{}'

# list specific doc
curl http://localhost:8355/docs/10593866

# add new file to our doi resource
# for that see response above for bucket uuid
curl -X PUT http://localhost:8355/add/d21bc027-8218-4357-8640-07e7f44976b1/bla.md --upload-file bla.md
```

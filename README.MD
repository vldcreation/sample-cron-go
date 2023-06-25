# github.com/vldcreation/sample-cron-go
Disclaimer!!!
This is just simple program, no pattern write in this project

## How to Use
```
1. install depdencies
make install

2. setup environtment file
in this project i use yaml for management environtment
just run: 
cp internal/config/gcs_account.json.dist internal/config/gcs_account.json #setup config gcs service account
cp internal/config/pubsub_account.json.dist internal/config/pubsub_account.json #setup config pubsub service account
go to https://console.cloud.google.com/ for reference
cp config.yaml.dist config.yaml
setup all the env based on your own

3. make sure you coverage all environtment

4. Just Run Code Gunners !!!
make run

3. Run With opt params
by default use default configuration based on env dev
opts param:
 APP_ENV
 STORAGE_BUCKET
if you want to dynamicaly switching the storage, add those env while run program

make run APP_ENV=dev STORAGE_BUCKET=accell-go # dev, prod -> assume dev for using minio and prod use gcs

```
## Changelog (Based on accel quiz)
1. Setup project [Y]
2. Setup storage [Y]
3. Implement switcher [Y]
4. Implement Pub/Sub [Y]

## TODO
1. Implement Design Pattern
2. implement asym encryption
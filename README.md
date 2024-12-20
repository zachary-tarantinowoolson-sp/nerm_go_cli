# CLI Tool to make API requests and generate files

Built with GO and Cobra. This utilized the [go-keyring](https://www.github.com/zalando/go-keyring) to store and retrive secrets

## Installing

Download and install the latest version of Go: [Download and Install](https://go.dev/doc/install)

While still working on the v1 release of this, you can download the pre-release versions under "Releases". Download the Zip and extract it to a directory of your choice. Open a terminal window in the nerm_go_cli-x folder, and run `go build -buildvcs=false; go install -buildvcs=false`. 

Then, you will be able to run the NERM CLI commands from any directory. 


## Usage

In a terminal, type `nerm -h` to see the available commands. Example:

Use `nerm env` to CRUD environments in order to make use of the other commands (the nerm_config.yaml files gets created in the .nerm folder of your User directory)
User `nerm profiles get` with optional flags to pull a JSON and CSV report of Profile dat from a tenant

AFTER ID usage
to get all profiles: nerm profiles get --after_id=""
to get profiles after a certain page : nerm profiles get --after_id profile_id

### Configuration
There are default settings configured in the `nerm_config.yaml` file (in the .nerm folder of your User directory). These are:
- default_output_location : Currently set to `default_output_location` . This is where files generate by this CLI tool will be sent to.
- limit : Currently set to `100`. This is the value which feeds the `limit` query parameter for GET requests.


#### ToDo
- [ ] Updating profiles
    - [ ] using JSON from a File
    - [ ] using prompts or flags for what attributes / values to set
- [ ] Creating profiles
    - [ ] using JSON from a File
    - [ ] using prompts or flags for what attributes / values to set
- [ ] Consolidation reporting
    - Get records
    - Delete records
    - Importer
- [x] Workflow Session searching and reporting
    - pull last x days of failed workflows 
    - [ ] fix progress bar reporting numbers to not just be the get max if using -d
- [ ] Better input error checking (number of profile type, env, etc is within range)
- [ ] Job status table for mass profile change / import
- [x] Change Yaml to https://github.com/zalando/go-keyring

# CLI Tool to make API requests and generate files

Built with GO and Cobra

## Installing

Download and install the latest version of Go: [Download and Install](https://go.dev/doc/install)

While still working on the v1 release of this, you can download the pre-release versions under "Releases". Download the Zip and extract it to a directory of your choice. Open a terminal window in the nerm_go_cli-x folder, and run `go build -buildvcs=false; go install -buildvcs=false`. 

Then, you will be able to run the NERM CLI commands from any directory. 


## Usage

In a terminal, type `nerm -h` to see the available commands. Example:

Use `nerm env` to CRUD environments in order to make use of the other commands (the nerm_config.yaml files gets created in the .nerm folder of your User directory)
User `nerm profiles get` with optional flags to pull a JSON and CSV report of Profile dat from a tenant

### Configuration
There are default settings configured in the `nerm_config.yaml` file (in the .nerm folder of your User directory). These are:
- default_output_location : Currently set to `default_output_location` . This is where files generate by this CLI tool will be sent to.
- limit : Currently set to `100`. This is the value which feeds the `limit` query parameter for GET requests.


#### ToDo
- [x] Build environment manager
    - [x] Create
    - [x] Show
    - [x] Update
    - [x] Delete
    - [x] Update
    - [x] List all
- [x] Health Check
- [x] Pulling profiles into a report
    - Allow query parameters 
    - Print to a file / json to csv
- [x] Basic Profile Counts
    - Show Profile counts based on status and for each profile type
    - Table display
- [ ] Advanced Searching
    - [ ] See all saved Searches (list)
    - [ ] Create/Run Search via file (-f)
    - [ ] Create Search via prompts (-c)
        - Utility for getting all stored attributes
        - Utility for getting Profile Types
        - Creating a search both saves it to the app and stores the id / name to yaml
- [ ] Updating profiles
    - [ ] using JSON from a File
    - [ ] using prompts or flags for what attributes / values to set
- [ ] Creating profiles
    - [ ] using JSON from a File
    - [ ] using prompts or flags for what attributes / values to set
- [ ] IDP reporting
    - Get records
    - Delete records
    - Importer
- [ ] Workflow Session searching and reporting
    - pull last x days of failed workflows 
    - using a settings file to link workflow name to workflow ID (readability)
- [ ] Input error checking (number of profile type, env, etc is within range)
- [ ] Job status table for mass proifle change / impot
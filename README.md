# CLI Tool to make API requests and generate files

Built with GO and Cobra

## Configuration
There are default settings configured in the `settings.env` file. These are:
- OUTPUT_FOLDER : Currently set to `default_output_location` . This is where files generate by this CLI tool will be sent to.
- DEFAULT_LIMIT_PARAM : Currently set to `100`. This is the value which feeds the `limit` query parameter for GET requests.
- DEFAULT_GET_LIMIT : Currently set to `Float::INFINITY`. This is the value which controls when bulk GET requests will stop. At Infinity, this will GET records until all records are recieved


## ToDo
- :heavy_check_mark: Build environment manager
    - :heavy_check_mark: Create
    - :heavy_check_mark: Show
    - :heavy_check_mark: Update
    - :heavy_check_mark: Delete
    - :heavy_check_mark: Update
    - :heavy_check_mark: List all
- :heavy_check_mark: Health Check
- [ ] Pulling profiles into a report
    - Allow query parameters 
    - Table display / print to a file
- [ ] Basic Profile Counts
    - Show Profile counts based on status and for each profile type
    - Table display / print to a file
- [ ] Advanced Searching
- [ ] IDP reporting
    - Get records
    - Delete records
    - Importer
- [ ] Workflow Session searching and reporting
    - pull last x days of failed workflows 
    - using a settings file to link workflow name to workflow ID (readability)
- [ ] Input error checking (number of profile type, env, etc is within range)
- [ ] Job status table for mass proifle change / impot
- [ ] Add limit, profile type and forcebacked to pull_profiles
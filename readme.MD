# drivvn backend test

## Description

Basic API set up to the specs of https://github.com/drivvn/backend-task - I probably took a bit longer than the time given, perhaps 5 hours in total but I was enjoying the task :)

Written in Go using the (Go Fiber Framework)[https://gofiber.io/], and using an SQLite database.

Data is stored in `./data/cardata.db` and `./data/cardata_test.db`, the latter is identically set up, and used for testing.

## Building and running

Built binaries are in /bin for Windows, Mac and Linux, just run to start the API - defaults to `localhost:8000`

You can build from source with `go build`. I've also included a basic shell script which builds for Windows, Mac and Linux all at once - see `./build_all.sh`. This will place binaries into `/bin` which can be executed, e.g `./bin/drivvnapi_windows_amd64.exe`

If you provide the `-r` flag the database will be reset to contain one car and the standard four colours at startup.

If you provide the `-c` flag the database will be reset to contain no car or colour information at startup.

The above information is also available with `-help`

## Usage

See Future Improvements below for Future iterations would document the API using the Open API spec, owing to its ease of use, wide array of tooling, and good support for generating documentation and code directly from the spec.

For now here are curls to demonstrate the API functionality.

### Get All Cars

```
curl --request GET \
--url http://127.0.0.1:8000/cars
```

### Get Car By ID

```
curl --request GET \
--url http://127.0.0.1:8000/car/1
```

### Delete Car By ID
Note the request for this endpoint was to be /cars/:id but I felt /car/id made more sense as we only delete a single car. Easily changed though.

```
curl --request DELETE \
--url http://127.0.0.1:8000/car/1
```

### Add Car

```
curl --request POST \
  --url http://127.0.0.1:8000/cars \
  --header 'Content-Type: application/json' \
  --data '[
	{
		"Make": "BMW",
		"Model": "3 Series",
		"BuildDate": "2022-04-04",
		"ColourID": 1
	}
]'
```

## Testing

Run `go test` to run tests

Testing in Go using the built in testing library currently doesn't seem to make it hugely easy to run before and after hooks, so currently each test contains a start and finish line to reset the test DB back to its default state.

Potential future improvement is to switch to a more fully featured community library, or investigate whether TestMain can be used for this purpose - my initial reading suggested that it could only provide a before all and after all but not a before each and after each. 

Nevertheless I think my test code could be shortened by treating all tests as a struct slice and passing the test function in each time, as a future improvement.


## Requirements checks, issues, future improvements

### Requirements Checks

Minimum:
- The API should allow the adding, deleting and listing of cars ✓
- The age of each car submitted can not be older than four years ✓
- There should be four initial colour options - red, blue, white and black ✓
- The API should respond with appropriate HTTP response codes and messages ✓
- The API should accept and return valid JSON ✓
- A suite of suitable tests should be created for these requirements ✓

Optional:
- Endpoints to add, update, delete and list additional colours ✘ (adding these would take the same approach as car endpoints however)
- A short description of how extra data models could improve the design ✓ (see below)
- A short description on how best the API could be documented ✓ (see below)

### Extra Data Models

Given the constantly fluctuating nature of car configurations, it may be a good improvement for Make and Model to themselves be data models. In particular, a Make can have many Models so forming a one to many relationship would provide a nice safeguard for data integrity.

### API Documentation

I'd plan on adding an OpenAPI 3 yml spec file to document the API, and generating a static HTML of the documentation to then serve at the root of the project, or if root is required in the future then at /help or /docs. In my experience adding an OpenAPI doesn't take much time, is universally understood and is both human readable as well as supported by numerous tools for generating code, request collections and documentation from.

### Other Potential Improvements

Route grouping: Would make it easier to manage at a larger scale (for example) /api/v1 using [Fiber's grouping](https://docs.gofiber.io/guide/grouping/).

Database setup: Setup and resetting is pretty hardcoded at the moment although the `database.go` file is mostly intended to just provide utility functions to enable the easier testing of the API itself. In reality those functions would likely not be needed and in fact would be rather dangerous since they reset the data at the drop of a hat!

Testing: The system tests currently only check for status codes and could be much improved with error state checking, as well as response body checks.

Configuration: Much configuration is done at the top of `main.go` regarding database paths, IPs and ports. This would need extracting for use in a cloud environment and more command line flags or reading and parsing a config file if used in a cloud environment.

Authentication: It's likely that authentication would be needed on the POST and DELETE endpoints, depending on usage. My first thought would be to use the Fiber middleware for [Basic Auth](https://docs.gofiber.io/api/middleware/basicauth) or [Keyauth](https://docs.gofiber.io/api/middleware/keyauth)

Rate Limiting: Depending on the use of the API and whether this is already provided at another layer (for example by Azure Front Door or equivalent), rate limiting could likely be achieved fairly easily using [Limiter](https://docs.gofiber.io/api/middleware/limiter)

Error messaging: Very bare bones currently, and could do with being fleshed out to perhaps use error codes and a standardised response layout of `{"Error": errorCode, "Message": errorMessage}`
# VMA Simple User Rest API CRUD operations

This project aims to provide a basic Rest API to Create/Read/Update/Delete users using MongoDB and a security/trust acceptable level.

Requirements:
 - Keep the user information: Name, Age, Email, Password and Address
 - Use Go lang or Node.js
 - Use MongoDB or MySql
 - Key evaluation items: robustness, scalability, performance and security

## Development environment specs

Visual Studio Code (version 1.79.2)
Go lang (version go1.20.7 windows/amd64) library:
 - github.com/golang-jwt/jwt v3.2.2+incompatible
 - github.com/gin-gonic/gin v1.9.1
 - github.com/swaggo/swag v1.16.1
 - golang.org/x/crypto v0.9.0
 - golang.org/x/net v0.10.0
 - go.mongodb.org/mongo-driver v1.12.0
Atlas MongoDB (https://cloud.mongodb.com/):
 - M0 Free cluster MongoDB 6.0
 - Replication Factor 3 Nodes
 - X.509 Certificates Authentication
 - Maximum of 500 simultaneous connections
 - Maximum of 100 databases and 500 collections
 - Able to deploy to MS Azure, AWS and GCP in a subset of regions
 - Data Transfer Limits are 10 GB in and 10 GB out per period
 - Data recovery: always takes a single daily snapshot at the same time
 - Alerts & Monitoring: Connections, Logical Size, Network and Opscounter
 - Throughput: 100 operations per second
Microsoft Windows 11 Pro x64-based PC (Version 10.0 Build 22621)
Processor: 11th Gen Intel(R) Core(TM) i7-11800H CPU @ 2.30GHz, 16 Core(s)
Total Physical Memory 32.7 GB

## 1.1 Configuration file

The main feature of this parser is the flexibility regarding column header identification. It allows changing the header keyword to capture the desired fields properly using a regular expression format. This behavior makes it possible to process a great number of csv files, without being stuck with a rigid property pattern.

For a better experience and clarification, follow a short description regarding the configuration file and its properties:
 - **id**: related to the employee/data row identificator, it is a key item and must be valid and unique.
 - **name**: the candidate name, is possible to use the first name or the surname (default is the first name).
 - **email**: a valid email to contact the candidate, it is also a key item and must be valid and unique.
 - **salary**: the payment amount, no period, type or currency restrictions (weekly, monthly, annual, USD, CAD, etc).

The regular expression language was selected because it is easy to use and has been broadly applied to identify textual patterns. You can check for references at [Regular Expression Language - Quick Reference](https://docs.microsoft.com/en-us/dotnet/standard/base-types/regular-expression-language-quick-reference).


Follow the sample configuraiton available for testins:
```json
{
    "config" : {
        "header" : {
            "id": "(?i)number|id|emp",
            "name" : "(?i)name|first",
            "email" : "(?i)mail",
            "salary" : "(?i)wage|rate|salary"
        },
        "outorder" : {
            "id": 0,
            "name": 1,
            "email": 2,
            "salary": 3
        }
    }
}
```

The `header` property contains the expressions to identify each column name. The intention is to obtain the right column ID where the desired information is declared.

The `outorder` property represents the output file column order. The current version can process only 4 elements, so, the indexes are from 0 to 3. It is just to adjust the order of the columns in the parsed file.

If a unique value violation occurs or a key property isn't available/right the process action will fail and a bad file will be generated.

## 1.2 Flags description

There are three available flags to declare, the csv file to be processed, a configuration file, and the output file. It's recommended to write the full path for those items or use in the same level from the binary file:
 - **config**: the configuration file as cited at [section 1.1](#1.1-configuration-file).
 - **input**: a csv file containing the information to be parsed and organized.
 - **output**: the output file that will have the processed information.

Example with all flags declared with the full path and file name:
```shell
rainparser.exe -config="C:\temp\rainparser\sample\config.json" -input="C:\temp\rainparser\sample\roster1.csv" -output="C:\temp\rainparser\sample\roster_parsed.csv"
```

If the files are at the same directory as the binary you can delcare like?
```shell
rainparser.exe -config="config.json" -input="roster1.csv" -output="roster_parsed.csv"
```

If no flag information provided, the program will try to use the default settings:
 - **config**: `<bin folder>\sample\config.json`.
 - **input**: `<bin folder>\sample\roster1.csv`.
 - **output**: `<bin folder>\sample\roster1_parsed.csv`.

# 2 Chosen architecture

The initial approach adopted was to design the behavior to see how attached the understanding to the requirements was.

So, the sequence diagram was created as follows:
![seqdiagram](vmausers_api_diagram-all.png)

The first archiceture chosen was related with a Model View Controller aiming to meet the requirements in the construction of a project that can be scalable and that possible changes in any of the layers are made without interference in the other layers. MVC is based on the separation of data (model), user interface (view), and business logic (controller).

Using this design approach the following segregation was applied:
 - **Model**: items related with the data manipulation, basically entities classes.
 - **View**: the user will interact through the command line interface.
 - **Controller**: the logic applied, the controller and the engine under the hood (standardizer and the parser).

With this in mind was more comfortable completing the sequence diagram. Also thinking to provide a clean and robust code, the intent of separating each action as little as possible in each method was almost entirely achieved.

# 3 Future improvements

Regarding the available time and the initial requirements a Minimal Viable Program was provided, to confirm the achievements and if the path is clear regarding the development and the delivered features.

However, there are a couple of adjustments that could bring more reliability, flexibility and usage for a parser like that:
 - the bad file with detailed errors: as there wasn't a bad file description, the approach adopted was to only copy the current input renaming to a new name, maybe beyond this copy would be great to add the flaw description.
 - multiple file processing: using the multithreaded resource to process a bunch of files simultaneously in a background task
 - create a flexible entity: the current development has its entity fixed with 4 elements (id, name, email, salary), it would be great to use an extern modeling file (XSD-like) to build the entity at runtime.
 - cyclic log file: to improve debug and problem tracking
 - support large files: the current program has a limited workspace to process files, to be able to process large files instead of reading a file from beginning to end, the proper approach could be processing each line or fixed blocks of text and buffering the results to save when the parser is complete.

Beyond these items, thinking of an online environment would be great to design APIs to attend this program and, to be able to process a large number of requests, like a streaming situation, maybe a message hub would be required and another evolution of this tool to support and perform good autoscaling to attend a demand. Also items like a Redis for caching this information and gain speed. As there is no limit to ideas, see if it's legal to store the processed information and use it with an ML to predict trends and obtain analysis regarding a topic.

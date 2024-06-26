# CLI

##  shield action 

Manage actions

###  shield action create [flags] 

Create an action

```
-f, --file string     Path to the action body file
-H, --header string   Header <key>:<value>
````

###  shield action edit [flags] 

Edit an action

```
-f, --file string   Path to the action body file
````

###  shield action list 

List all actions

###  shield action view 

View an action

##  shield auth 

Auth configs that need to be used with shield

##  shield completion [bash|zsh|fish|powershell] 

Generate shell completion scripts

##  shield config

Manage client configurations

###  shield config init 

Initialize a new client configuration

###  shield config list 

List client configuration settings

##  shield environment 

List of supported environment variables

##  shield group 

Manage groups

###  shield group create [flags] 

Create a group

```
-f, --file string     Path to the group body file
-H, --header string   Header <key>:<value>
````

###  shield group edit [flags] 

Edit a group

```
-f, --file string   Path to the group body file
````

###  shield group list 

List all groups

###  shield group view [flags] 

View a group

```
-m, --metadata   Set this flag to see metadata
````

##  shield namespace 

Manage namespaces

###  shield namespace create [flags] 

Create a namespace

```
-f, --file string   Path to the namespace body file
````

###  shield namespace edit [flags] 

Edit a namespace

```
-f, --file string   Path to the namespace body file
````

###  shield namespace list 

List all namespaces

###  shield namespace view 

View a namespace

##  shield organization 

Manage organizations

###  shield organization admadd [flags] 

add admins to an organization

```
-f, --file string   Path to the provider config
````

###  shield organization admlist 

list admins of an organization

###  shield organization admremove [flags] 

remove admins from an organization

```
-u, --user string   Id of the user to be removed
````

###  shield organization create [flags] 

Create an organization

```
-f, --file string     Path to the organization body file
-H, --header string   Header <key>:<value>
````

###  shield organization edit [flags] 

Edit an organization

```
-f, --file string   Path to the organization body file
````

###  shield organization list 

List all organizations

###  shield organization view [flags] 

View an organization

```
-m, --metadata   Set this flag to see metadata
````

##  shield policy 

Manage policies

###  shield policy create [flags] 

Create a policy

```
-f, --file string     Path to the policy body file
-H, --header string   Header <key>:<value>
````

###  shield policy edit [flags] 

Edit a policy

```
-f, --file string   Path to the policy body file
````

###  shield policy list 

List all policies

###  shield policy view 

View a policy

##  shield project 

Manage projects

###  shield project create [flags] 

Create a project

```
-f, --file string     Path to the project body file
-H, --header string   Header <key>:<value>
````

###  shield project edit [flags] 

Edit a project

```
-f, --file string   Path to the project body file
````

###  shield project list 

List all projects

###  shield project view [flags] 

View a project

```
-m, --metadata   Set this flag to see metadata
````

##  shield role 

Manage roles

###  shield role create [flags] 

Create a role

```
-f, --file string     Path to the role body file
-H, --header string   Header <key>:<value>
````

###  shield role edit [flags] 

Edit a role

```
-f, --file string   Path to the role body file
````

###  shield role list 

List all roles

###  shield role view [flags] 

View a role

```
-m, --metadata   Set this flag to see metadata
````

##  shield server

Server management

###  shield server init [flags] 

Initialize server

```
-o, --output string      Output config file path (default "./config.yaml")
-r, --resources string   URL path of resources. Full path prefixed with scheme where resources config yaml files are kept
                         e.g.:
                         local storage file "file:///tmp/resources_config"
                         GCS Bucket "gs://shield-bucket-example"
                         (default: file://{pwd}/resources_config)
                         
-u, --rule string        URL path of rules. Full path prefixed with scheme where ruleset yaml files are kept
                         e.g.:
                         local storage file "file:///tmp/rules"
                         GCS Bucket "gs://shield-bucket-example"
                         (default: file://{pwd}/rules)
````                      

###  shield server migrate [flags] 

Run DB Schema Migrations

```
-c, --config string   Config file path
````

###  shield server migration-rollback [flags] 

Run DB Schema Migrations Rollback to last state

```
-c, --config string   Config file path
````

###  shield server start [flags] 

Start server and proxy default on port 8080

```
-c, --config string   Config file path
````

##  shield user 

Manage users

###  shield user create [flags] 

Create an user

```
-f, --file string     Path to the user body file
-H, --header string   Header <key>:<value>
````

###  shield user edit [flags] 

Edit an user

```
-f, --file string   Path to the user body file
````

###  shield user list 

List all users

###  shield user view [flags] 

View an user

```
-m, --metadata   Set this flag to see metadata
````
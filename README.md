# Go Buddy CLI

Go Buddy is a command-line tool designed to interact with the Buddy CI/CD Platform. It allows users to manage and run pipelines for continuous integration (CI) and continuous deployment (CD). With this tool, you can easily deploy to staging or production environments, ensuring a smooth and automated workflow for your development and deployment processes.

## Installation

```
go install https://github.com/JacobAndrewSmith92/gobuddy@latest
```



## Available Commands 
1. `config`
2. `deploy`



### Setting Up Your Workspace With `config`
The config command is used to create and update your configuration while using this tool. It provides subcommands to get, set, and reset configuration settings. This setup is required in order to use the tool. Below are the details on how to use these subcommands.

**Need to create your buddy token first?** [Set it up here]()

#### Commands

| Subcommand | Description                |
| :-------- |  :-------------------------|
| `get` | Retrieves your current configuration |
| `set` | Create a new configuration |
|`set <key> <value>`| Set a specific configuration key |
|`reset` | Reset your configuration |

##### Examples

**`get`**


```bash
$ gobuddy config get

Current Configuration:
Token: a-valid-token
Workspace: fizzbuzz
```

**`set`** (if one does not exist)

```
$ gobuddy config set

Enter your Buddy API token: ********
Enter your Buddy workspace: foobar
```

**`set <key> <value>`**

One of:
- `token`
- `workspace`

```bash
$ gobuddy config set token some-value

Token updated to: some-value
Configuration updated successfully!
```
```bash
$ gobuddy config set workspace some-workspace

Workspace updated to: some-workplace
Configuration updated successfully!
```

**`reset`**

```bash
$ gobuddy config reset

Are you sure you want to reset the configuration? (yes/no): yes
Configuration has been reset.
```

### Deploying A Service With `deploy`
This command allows you to trigger a deployment pipeline for your project using the Buddy CI/CD platform. You can select a project, branch, and pipeline for deployment either by passing arguments and flags or through an interactive process.

#### Commands

| Subcommand | Type | Description                |Required |
| :-------- | :------ |  :-------------------------|:--------|
| `<project>` | `argument` | Pass the project name (repo) you want to deploy  |`false`|
| `-b or --branch` |`flag`| Pass this flag followed by a value if you want to specify your own git branch | `false`|
|`-p or --pipeline`|`flag`| Pass this flag followed by a value if you want to specify your own pipeline ID |`false`|

#### Interactive
If you don’t pass all arguments and flags, Go Buddy will pick up where you left off and guide you through some interactive steps:

1. **Project Selection**: Displays a list of available projects.
2. **Branch Selection**: Fetches and displays the branches for the selected project.
3. **Pipeline Selection**: Displays a list of pipelines associated with the project.

##### Examples

**Running with no arguments/flags passsed**
```bash
$ gobuddy deploy
```
**Running with the project**
```bash
$ gobuddy deploy project-foobar
```
**Running with the project and branch passed**
```bash
$ gobuddy deploy project-foobar -b fizz-buzz
```

**Running with all flags passed**
```bash
$ gobuddy deploy project-foobar -b fizz-buzz -p 12345
```

### Check Pipeline Status
Once you have ran a deployment, Go Buddy will ask you if you'd like to check the status of the deployment. You can do so by typing yes. As of today (09/18/2024), if you select no, you won't be able to check the status again. That logic will come in future improvements.

#### [Known Statuses](https://buddy.works/docs/api/pipelines/executions/get-details-and-logs)
- `SUCCESSFUL`
- `FAILED`
- `INPROGRESS`
- `ENQUEUED`
- `SKIPPED`
- `TERMINATED`
- `NOT_EXECUTED`
- `INITIAL`
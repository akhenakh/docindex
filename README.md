# docindex

An indexer for mkdocs documentation, it uses [Bleve](https://github.com/blevesearch/bleve), so [complex queries](http://blevesearch.com/docs/Query-String-Query/) can be used.


First index your data
```
Usage of ./index:
  -docPath string
        mkdocs config path (default "mkdocs.yml")
  -indexPath string
        index db path (default "index.db")
```

Then use the query tool

```
Usage of ./query:
  -indexPath string
        index db path (default "index.db")
  -limit int
        limit response count (default 10)
  -lucky
        I feel lucky, open the first found result
  -q string
        query string
```
## Example:
```
./index -docPath ../my-docs/mkdocs.yml
2020/10/13 09:04:32 main.go:145: Indexing: [index.md] Home Home [Home,Home]
2020/10/13 09:04:32 main.go:145: Indexing: [office-hours.md] Office Hours [Home,Office Hours]
2020/10/13 09:04:32 main.go:145: Indexing: [contributing/contributing.md] Contributing [Home,Contributing,MkDocs]
2020/10/13 09:04:32 main.go:145: Indexing: [contributing/tips.md] Tips and Tricks [Home,Contributing,Tips and Tricks]
...

./query -q "drone -gcp" 
63 matches, showing 1 through 10, took 11.469714ms
    1. drone/v1/cli.md (0.972849)
        Title
                Tools & Platforms CI/CD <mark>Drone</mark> CLI
        Tags
                Tools & Platforms CI/CD <mark>Drone</mark> CLI
        Body
                … If you used the <mark>Drone</mark> v0.8 CLI tool, you'll need to update to the latest version and configuration for it to work properly with <mark>Drone</mark> v1. Install the latest version of the CLI version using the follo…
    2. drone/v1/run_jobs_manually.md (0.715997)
        Tags
                Tools & Platforms CI/CD <mark>Drone</mark> Run Jobs Manually
        Body
                … webhooks, you can run <mark>Drone</mark> jobs on your computer's Docker daemon.
We still recommend using the server during your day-to-day development process. !!! info
    Ensure the <mark>Drone</mark> CLI is  , as well as D…
    3. drone/README.md (0.690244)
        Tags
                Tools & Platforms CI/CD <mark>Drone</mark> Overview
        Body
                <mark>Drone</mark> is an   built on container technologies.
    ....
```


## TODO web frontend for mkdocs
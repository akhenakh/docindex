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

## TODO web frontend for mkdocs
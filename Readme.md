# Installation

## For driver to golang
```
go get github.com/lib/pq
```

## For jsonpath library
```
go get github.com/oliveagle/jsonpath
```

## Db

```
psql -a -U postgres --password -f grants.sql
psql -a -d hotels -U hotels_manager -f schema.sql
```



# Configuration
At config.json you find credentials for connection at db, locale (because in two files use different langugages) and parsers config
for keys to retrieving each filed of hotel.
In json parser usings jsonp (except photos, because error in lib when querying for array).
In xml parser using pathes of tags (no xpath).
In csv only names of columns.

In each parser config, photos is similar configuration for extract photo array from source. It have two parameters
1) In - the key or path for element which contains array of photos
2) Key - the key or path for tag which contains photo data

In json files i see the stars value is bigger than real. And have conclusion that it must divide to 10. It is in config too.
Also in XML parser config you can see root 'In' element - name of tag which contains hotel info.


# Run
```
go run main.go <your file with extension [json|xml|csv]>
```

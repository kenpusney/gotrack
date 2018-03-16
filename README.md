# KANA

Kana is a golang implementation of https://github.com/kenpusney/simple-analytics .

### Usage:

Build and run the server (see following section), and just include following script in
your website.

```html
<script src="http://localhost:10086/ka.js?id=chooseYourOwnId"></script>
```

Replace the host and `id` param with your self id, and then every request will hit on
on this analytics server.

You can retrieve collected analytics data via report.php by using following 3 params:
```
/report.php?limit=100&skip=0
```

For each of these params:
  - `limit`: the amounts you want to retrieve, default: 20.
  - `skip`: the offsite you want to skip.

### Up and Running

Simply using following script to get an executable
```
export GOPATH=`pwd`
go get github.com/kenpusney/kana
```
the `kana` executable will be in your `./bin` path.

Move the template file `ka.js.tmpl` and config file `conf.json` to same directory,
then the server will be up and running now.

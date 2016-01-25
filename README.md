# osetup
windows, mac, linux software version management system.

````
windows : v0.1.1 Correspondence

mac : v0.1.0 Correspondence

linux : Not compatible
````

****

## install
git clone https://github.com/minaduking/osetup.git

****

## usage
You create config.json in the project directory.

#### run osetup

```
bin/osetup [project directory]
```

ex) bin/osetup ./sample/

#### create config.json ####
```
sample
├── config.json
```

```config.json
{
	"packages": [{
		"name": "java", 
		"version": "latest",
		"option": {
			"windows": {

			},
			"darwin": {
				"type": "cask",
				"tap": "caskroom/cask"
			},
			"linux": {
				"type": "cask",
				"tap": "caskroom/cask"
			}
		}
	}]
}
```

## License
MIT

****

## Contact

github
https://github.com/minaduking

twitter
https://twitter.com/minaduking

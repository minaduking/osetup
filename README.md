# osetup
windows, mac, linux software version management system.

````
windows : Not compatible

mac : v0.1 Correspondence

linux : Not compatible
````

****

# install
git clone https://github.com/minaduking/osetup.git

****

# usage
bin/osetup [project directory]

You create config.json in the project directory.

#### create config.json ####
````
sample
├── config.json
````

ex) bin/osetup ./sample/

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

****

# Contact

github
https://github.com/minaduking

twitter
https://twitter.com/minaduking

# HPE Docker Machine OneView plugin

[![Build Status](https://travis-ci.org/HewlettPackard/docker-machine-oneview.svg?branch=master)](https://travis-ci.org/HewlettPackard/docker-machine-oneview)

The Docker Machine plugin for HPE OneView automates the provisioning of physical infrastructure on-demand, from a private cloud using templates from HPE OneView to enable customers to treat infrastructure-as-code.  

## Instructions and Setup Steps

Get the instructions for setting up the plugin with [HPE OneView here](/docs/oneview.md).

## Testing your changes

### From a container
```
USE_CONTAINER=1 make test
```

### From your local system
* Install golang 1.6 or better
* Install go packages listed in .travis.yml

```
make test
```

## Building a plugin
```
USE_CONTAINER=1 make build
```

## Contributing

Want to hack on docker-machine-oneview? Please start with the [Contributing Guide](https://github.com/HewlettPackard/docker-machine-oneview/blob/master/CONTRIBUTING.md).

## License
This project is licensed under the Apache License, Version 2.0.  See LICENSE for full license text.

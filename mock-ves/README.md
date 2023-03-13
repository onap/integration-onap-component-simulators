# mock-ves
This is a simple ves mock created for ONAP integration tests.

Build image
===========
To build ves-mock image you can simply run command:

 ```
docker build . -t mock-ves
```

Run image
=========

To run ves-mock image use below command:
```
docker run -d --net=host --name mock-ves mock-ves
```

Stop container
==============
To stop the container use:
```
docker rm -f mock-ves
```


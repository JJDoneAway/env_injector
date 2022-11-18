# Problem to be solved
Most cloud providers inject all variables you need to run during the start and runtime of your application.

E.g.: You have a frontend single page web application. The application needs to know the URL of the backend services (REST) to get its content from. But the backend URL is different in each environment (DEV, TEST, QA, PROD). The cloud tells you the correct one during environment variables in your runtime container.

On the other hand a web application is not able to read environment variables during runtime any more as it just consists out of a set of html and JS files which are delivered via a ordinary web server. So it adds all variables during build time (in the CI/CD container) into the code. 

E.g.: You implement your single page web application with Vue.js. In Vue.js you make use of .env files to define your environment specific variables and use them in the code via `process.env.VUE_APP_XYZ`. If you build your application (`npm run build --mode=production`) the values out of .env.production will be written into your code and can't be changed later on any more.

This gap between build time and runtime injection of environment variables must be gapped.

# The solution
1. You create a extra env file `.env.odj`and add a mapping between your local variables and the runtime environment names:
e.g.:
```
VUE_APP_URL=ODJ_DEP_BACKEND4ADDRESSBOOK_URI
VUE_APP_SERVICE=ODJ_DEP_BACKEND4ADDRESSBOOK_SERVICE
```  

2. You use your variables during your coding in the common way:
e.g.:
```
...
axios
    .get(`${process.env.VUE_APP_URL}/status/up`)
    .then((response) => {
...
```

2. You build your application in the common way
e.g.: 
```
npm run build
```

3. You copy your distribution files, the .env.odj file and this script (`env_injector`) to your docker file
e.g.:
```
FROM nginx

COPY nginx.conf /etc/nginx/nginx.conf
COPY env_injector /usr/share/env_injector
COPY .env.odj /usr/share/.env.odj
COPY dist/ /usr/share/nginx/html
COPY .well-known/ /usr/share/nginx/html/.well-known/

CMD ["/bin/sh",  "-c",  "/usr/share/env_injector /usr/share/nginx/html/ /usr/share/.env.odj && exec nginx -g 'daemon off;'"]
```

==> Just before the web server is started `env_injector` is updating all the system environment names with its correct values 


# Run & Build
## how to call the script
```
env_injector {path to your distribution folder} {path to your .env.odj file}
```

## during development
### run the code
```
go run . {path to your distribution folder} {path to your .env.odj file}
```

### test the code
in the `test`folder you'll find a bigger test which can be run again adn again to test your expected behavior.
```
cd ./test
chmod 777 runTest.sh
./runTest.sh 
```


## build the code
Please mind that go applications must be compiled for the target environment. This is in most cases a Linux distribution (Docker container).
For this you must hand over the corresponding go env parameters:

```
env GOOS=linux GOARCH=amd64 go build github.com/JJDoneAway/env_injector
```

# Contact
If you would like to discuss or extend the concept please feel free to send me a mail hoehne@magic-inside.de
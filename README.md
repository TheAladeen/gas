AGS - AngularJS + Golang + Sqlite
===

## Intro

This project intended to build a brand new shopping cart web application base on the latest AngularJS and Google Golang, both of them are very new and interesting to me.

After spent months of my spare time on designing and coding, I actually made a very simple yet functional prototype, include user, product, order, report modules, and a blogging module. Enquire module enables user to send messages to the owner to enquire business solution.

I hope to grow this project to a fully functional shopping cart application, forks and patches are welcome.


## Design - keep it simple

AngularJS is responsible for layout, sending restful request to server, and presenting data retrieved from server. (ags/webapp)

While Golang modules are responsible for providing restful api to angularjs webapp. (ags/service)

Using Sqlite for database mainly because it's simplicity, and should be more enough for me.


## Install - Linux & MacOS

1. Nginx.   
        Nginx is used for providing static resources, while restful requests were forwarded to 8080 port. Check out the conf file for nginx at data/nginx.conf.
2. Golang.  
	Download tar.gz file and untar it to /golang/go.  in your ~/.bashrc  
	export GOROOT=/golang/go  
	export GOPATH=/golang/ext/  
        export PATH=$GOROOT/bin:$GOPATH/bin:$PATH  
3. go get github.com/featen/ags
4. Modify nginx conf for reflect the source dir above and run nginx.
5. Modify your site config at data/ags.conf.
6. go run app.go



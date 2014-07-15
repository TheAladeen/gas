AGS - AngularJS + Golang + Sqlite
===

This open source project provide you the fastest way to deploy your own frontpage and blogging web app, with a product module for promote your own business.

AngularJS is responsible for layout, sending restful request to server, and presenting data retrieved from server. (ags/webapp)

While Golang modules are responsible for providing restful api to angularjs webapp. (ags/service)

Using Sqlite for database mainly because it's simplicity, and should be more enough for me.


## Install

1. Install Golang.  
2. go get github.com/featen/ags
3. Modify your site config at data/ags.conf.
4. go run app.go

## Deploy to Openshift (free)

1. Create an openshift account and create a golang app - e.g. blog.
2. git clone your blog code from openshift, copy ags app code into blog dir.
3. rm web.go from original blog code, git add 'the rest code from ags', git commit -a; git push origin master.


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


## Install

1. Install Golang.  
2. go get github.com/featen/ags
3. Modify your site config at data/ags.conf.
4. go run app.go

## Deploy to Openshift (free)

1. Create an openshift account and create a golang app.
2. Install openshift command line tool and ssh to your app. If you have example app running, please kill it.
3. Because golang env is already in place, just repeat the 2-4 steps in Install section.



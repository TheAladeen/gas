angular.module('featen.product').factory("Products", ['$http', 'Alerts', function($http, Alerts) {


        // Get all lists.
        this.getall = function( scall, ecall) {
            var promise = $http.get("/service/product");
            var error = {
                type: "warning",
                strong: "Warning!",
                message: "Unable to retrieve all products. Try again in a few minutes."
            };
            Alerts.handle(promise, error, undefined, scall, ecall);

            return promise;
        };

        this.searchcount = function(searchtext, scall, ecall) {
       	 var promise = $http.get("/service/product/search/" + searchtext + "/count");
            var error = {
                type: "warning",
                strong: "Warning!",
                message: "No response..."
            };
            Alerts.handle(promise, error, undefined, scall, ecall);

            return promise;
       };
       
       this.search = function(searchtext, pagenumber, scall, ecall) {
       	 var promise = $http.get("/service/product/search/" + searchtext +"/page/"+pagenumber);
            var error = {
                type: "warning",
                strong: "Warning!",
                message: "No match..."
            };
            Alerts.handle(promise, error, undefined, scall, ecall);

            return promise;
       };

    this.getrand = function(scall, ecall) {
        var promise = $http.get("/service/product/");
        Alerts.handle(promise, undefined, undefined, scall, ecall);
        return promise;
    };

        this.get= function(name, scall, ecall) {
            var promise = $http.get("/service/product/" + name);
            var error = {
                type: "warning",
                strong: "Warning!",
                message: "Unable to retrieve products information. Try again in a few minutes."
            };
            Alerts.handle(promise, error, undefined, scall, ecall);

            return promise;
        };
        
        
        this.add = function(data, scall, ecall) {
            var promise = $http.post("/service/product/", data);
            var error = {
                type: "error",
                strong: "Failed!",
                message: "Add product failed."
            };
            var success = {
                type: "success",
                strong: "Success!",
                message: "Add product success."
            };
            Alerts.handle(promise, error, success, scall, ecall);
            return promise;
        };

        this.saveproduct = function(data, scall, ecall) {
            var promise = $http.put("/service/product/", data);
            var error = {
                type: "error",
                strong: "Failed!",
                message: "Update product info failed."
            };
            var success = {
                type: "success",
                strong: "Success!",
                message: "Update product info success."
            };
            Alerts.handle(promise, error, success, scall, ecall);
            return promise;
        };

        return this;
}]);


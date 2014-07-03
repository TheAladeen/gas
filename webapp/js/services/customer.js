angular.module('featen.customer').factory("Customers", ['$http', 'Alerts', function($http, Alerts) {
        
        this.searchcount = function(searchtext, scall, ecall) {
        	 var promise = $http.get("/service/customer/search/" + searchtext + "/count");
             var error = {
                 type: "warning",
                 strong: "Warning!",
                 message: "No response..."
             };
             Alerts.handle(promise, error, undefined, scall, ecall);

             return promise;
        };
        
        this.search = function(searchtext, pagenumber, scall, ecall) {
        	 var promise = $http.get("/service/customer/search/" + searchtext +"/page/"+pagenumber);
             var error = {
                 type: "warning",
                 strong: "Warning!",
                 message: "No match..."
             };
             Alerts.handle(promise, error, undefined, scall, ecall);

             return promise;
        };

        this.add = function(data, scall, ecall) {
            var promise = $http.post("/service/customer/", data);
            var error = {
                type: "error",
                strong: "Failed!",
                message: "Unable to add new customer."
            };
            var success = {
                type: "success",
                strong: "Success!",
                message: "New customer added."
            };
            Alerts.handle(promise, error, success, scall, ecall);
            return promise;
        };
        
        this.getcustomer = function(id, scall, ecall) {
          var promise = $http.get("/service/customer/"+id);
          var error = {
                type: "warning",
                strong: "Warning!",
                message: "Unable to get this customer."
            };
            Alerts.handle(promise, error, undefined, scall, ecall);

            return promise;
        };
        
        this.savecustomer = function(data, scall, ecall) {
            var promise = $http.post("/service/customer", data);
            var error = {
                type: "error",
                strong: "Failed!",
                message: "Unable to update customer info."
            };
            var success = {
                type: "success",
                strong: "Success!",
                message: "Update customer info success."
            };
            Alerts.handle(promise, error, success, scall, ecall);
            return promise;
        };
        
}]);


angular.module('featen.user').factory("User", ['$http', 'Alerts', function($http, Alerts) {
        this.data = {
			user: null,
			authenticated: false,
                    };
	
    this.get = function(call) {
        var promise = $http.get("/service/user/");
        
        Alerts.handle(promise, undefined, undefined, call);
    };
    
    this.signin = function(data, scall, ecall) {
        var promise = $http.post("/service/user/signin", data);
        var error = {
            type: "error",
            strong: "Failed!",
            message: "Could not sign in. Try again in a few minutes."
        };
        var success = {
            type: "success",
            strong: "Success!",
            message: "Sign in success."
        };
        Alerts.handle(promise, error, success, scall, ecall);

        return promise;
    };

    this.signout = function(data, scall, ecall) {
    	var promise = $http.post("/service/user/signout");
        var error = {
                type: "error",
                strong: "Failed!",
                message: "Could not signout."
            };
            var success = {
                type: "success",
                strong: "Success!",
                message: "Sign out success."
            };
            Alerts.handle(promise, error, success, scall, ecall);

        return promise;
    };
    
    return this;
}]);


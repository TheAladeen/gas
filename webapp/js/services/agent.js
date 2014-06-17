angular.module('featen.agent').factory("Agents", ['$http', 'Alerts', function($http, Alerts) {
        // Get all lists.
        this.getall = function(scall, ecall) {
            var promise = $http.get("/service/agents/");
            var error = {
                type: "warning",
                strong: "Warning!",
                message: "Cannot fetch this info, please try it later."
            };
            Alerts.handle(promise, error, undefined, scall, ecall);

            return promise;
        };


        this.get = function(id, scall, ecall) {
            var promise = $http.get("/service/agents/" + id);
            var error = {
                type: "warning",
                strong: "Warning!",
                message: "Cannot fetch this info, please try it later."
            };
            Alerts.handle(promise, error, undefined, scall, ecall);

            return promise;
        };



        this.create = function(data, scall, ecall) {
            var promise = $http.post("/service/agents/", data);
            var error = {
                type: "error",
                strong: "Failed!",
                message: "Cannot create this object, please try it later."
            };
            var success = {
                type: "success",
                strong: "Success!",
                message: "Object created."
            };
            Alerts.handle(promise, error, success, scall, ecall);

            return promise;
        };

        return this;
    }]);

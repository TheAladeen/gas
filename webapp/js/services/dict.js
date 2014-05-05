angular.module('featen.dict').factory("Dict", ['$http', 'Alerts', function($http, Alerts) {
        this.query = function(q, scall, ecall) {
            var promise = $http.get("/service/dict/" + q);
            var error = {
                type: "warning",
                strong: "Warning!",
                message: "无法查询该内容，请检查输入."
            };
            Alerts.handle(promise, error, undefined, scall, ecall);

            return promise;
        };

        return this;
    }]);

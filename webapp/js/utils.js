utils = angular.module('utils');

utils.factory("Global", [function() {
        var _this = this;
        _this._data = {
            user: window.adminuser,
            authenticated: !!window.adminuser,
        };

        return _this._data;
}]);

utils.factory("Alerts", ['$rootScope', function($rootScope) {
        this.alerts = [];
        var self = this;

        this.add = function(type, strong, message) {
            this.alerts.unshift({
                type: type,
                strong: strong,
                message: message
            });
            window.setTimeout(function() {
                $rootScope.$apply(self.remove(self.alerts.length - 1));
            }, 5000);
        };

        this.remove = function(index) {
            this.alerts.splice(index, 1);
        };

        this.handle = function(promise, error, success, scall, ecall) {
            promise
                    .success(function(data, status, headers, config) {
                        if (success !== undefined) {
                            self.add(success.type, success.strong,
                                    success.message);
                        }

                        if (scall !== undefined) {
                            scall(data, status, headers, config);
                        }
                    })
                    .error(function(data, status, headers, config) {
                        if (error !== undefined) {
                            self.add(error.type, error.strong,
                                    error.message);
                        }

                        if (ecall !== undefined) {
                            ecall(data, status, headers, config);
                        }
                    });
        };
        return this;
}]);

utils.factory("StageData", [function() {

        // for ProductAddController
        var StageData = [];
        var currIndex = 0;
        this.get = function(id) {
            var data;
            angular.forEach(StageData, function(d) {
                if (d.id === parseInt(id))
                    data = d.data;
            });
            return data;
        };
        this.add = function(adddata) {
            var i = currIndex++;
            StageData.push({id: i, data: adddata});
            return i;
        };
        this.del = function(id) {
            var oldStageData = StageData;
            StageData = [];
            angular.forEach(oldStageData, function(d) {
                if (d.id !== parseInt(id))
                    StageData.push(d);
            });
        };
        
        return this;
}]);


utils.controller("AlertController", ['$scope', "Alerts", function($scope, Alerts) {
		$scope.alerts = Alerts.alerts;
		
		$scope.remove = function(index) {
				Alerts.remove(index);
		};
}]);

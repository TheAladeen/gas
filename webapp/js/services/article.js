angular.module('featen.article').factory("Articles", ['$http', 'Alerts', function($http, Alerts) {
        this.get= function(name, scall, ecall) {
            var promise = $http.get("/service/article/" + name);
            var error = {
                type: "warning",
                strong: "Warning!",
                message: "Cannot get this blog."
            };
            Alerts.handle(promise, error, undefined, scall, ecall);

            return promise;
        };
        
        this.getPageArticles = function(page, scall, ecall) {
        	var promise = $http.get("/service/article/page/" + page);
        	var error = {
        			type: "warning",
        			strong: "Warning!",
        			message: "Can not fetch articles for current page, please try it later."
        	};
        	Alerts.handle(promise, error, undefined, scall, ecall);
        	return promise;
        };
        
        this.getTotalPageNumber = function(scall, ecall) {
        	var promise = $http.get("/service/article/totalpagenumber");
        	var error = {
        			type: "warning",
        			strong: "Warning!",
        			message: "Can not fetch total page number right now."
        	};
        	Alerts.handle(promise, error, undefined, scall, ecall);
        	return promise;
        };

        this.create = function(data, scall, ecall) {
            var promise = $http.post("/service/article/", data);
            var error = {
                type: "error",
                strong: "Failed!",
                message: "Cannot create article right now."
            };
            var success = {
                type: "success",
                strong: "Success!",
                message: "Article created."
            };
            Alerts.handle(promise, error, success, scall, ecall);

            return promise;
        };

        this.save = function(data, scall, ecall) {
            var promise = $http.put("/service/article/", data);
            var error = {
                type: "info",
                strong: "Failed!",
                message: "Cannot save your article right now."
            };
            Alerts.handle(promise, error, undefined, scall, ecall);

            return promise;
        };

        this.del = function(data, scall, ecall) {
            var promise = $http({
                method: 'DELETE',
                url: "/service/article/" + data.Id}
            );
            var error = {
                type: "error",
                strong: "Failed!",
                message: "Delete article failed."
            };
            var success = {
                type: "success",
                strong: "Success!",
                message: "Delete aritcle success."
            };
            Alerts.handle(promise, error, success, scall, ecall);

            return promise;
        };

        return this;
    }]);

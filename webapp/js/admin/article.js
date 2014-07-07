article = angular.module('article');

article.factory("Articles", ['$http', 'Alerts', function($http, Alerts) {
        this.get= function(id, scall, ecall) {
            var promise = $http.get("/service/article/" + id);
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


article.controller('ArticleEditController', ['$scope', '$window', '$document', '$routeParams', '$location', 'Global', 'StageData', 'Articles', function($scope, $window, $document, $routeParams, $location, Global, StageData, Articles) {
        $scope.data = {};
        $scope.editstate = true;
        
        $scope.preview = function(){
        	if ($scope.data.content !== undefined && $scope.data.content !== "") {
        		var htmlcontent = marked($scope.data.content);
                $('#htmlcontentdiv').html(htmlcontent);	
        	} else {
        		$('#htmlcontentdiv').html("");
        	}
            
            $scope.editstate = false;
        };

        $scope.edit = function() {
            $scope.editstate = true;
        };

        var savedDataId = $routeParams.SavedDataId;
        var uploadedUrlsId = $routeParams.UploadedUrls;
        $scope.getEditingArticle = function() {
            if (savedDataId !== undefined && savedDataId !== '') {
                var stagedata = StageData.get(savedDataId);
                if (stagedata !== undefined) {
                    $scope.data = stagedata;
                    StageData.del(savedDataId);
                } else {
                    var r = $location.path().split("/")[1];
                    $location.path("/" + r);
                }
                if (uploadedUrlsId !== undefined && uploadedUrlsId !== '') {
                    $scope.data.blogphotos = StageData.get(uploadedUrlsId).split(";");
                    $scope.data.coverphoto = $scope.data.blogphotos[0];
                    StageData.del(uploadedUrlsId);
                }
            } 
        };
        
        
        $scope.jumptoupload = function() {
            var stageDataId = StageData.add($scope.data);
            var r = $location.path().split("/")[1];
            var redirecturl = "/" + r + "/savedid/" + stageDataId;
            $location.path('/uploadfile/redirect/'+Base64.encode(redirecturl));
        };
      
        $scope.create = function() {
            Articles.create({"Info":JSON.stringify({"Title": $scope.data.title, "CoverPhoto": $scope.data.coverphoto, "Intro": $scope.data.intro, "Content": $scope.data.content, "CreateTime": Date.now()})}, function(l) {
                $location.path('/');
            });
        }; 
        }]);

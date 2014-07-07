product = angular.module('product');

product.factory("Products", ['$http', 'Alerts', function($http, Alerts) {
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

        this.get= function(id, scall, ecall) {
            var promise = $http.get("/service/product/" + id);
            var error = {
                type: "warning",
                strong: "Warning!",
                message: "Sorry, unable to retrieve products information."
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

        this.update = function(data, scall, ecall) {
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

product.controller('ProductsController', ['$scope', '$routeParams', '$location', 'Global', 'Products', function ($scope, $routeParams, $location, Global, Products) {

    $scope.searchtext = "";
    $scope.searchcount = {};
    $scope.currPage = 1;
    $scope.totalPageNumber = 1;

    $scope.search = function() {
    	$scope.currPage = 1;
    	var t = $scope.searchtext;
    	if ($scope.searchtext.length == 0)
    		t = " ";
    	Products.searchcount(t, function(sc) {
    		$scope.searchcount = sc;
    		$scope.totalPageNumber = Math.ceil(sc.Total / sc.PageLimit);
    	});
    	Products.search(t, 1, function(cs) {
    		$scope.products = cs;
            $scope.products.forEach(function(elem, index, array) {
                elem.Info = JSON.parse(elem.Info);
            });
    	});
    };

    $scope.setpage = function(n) {
    	var t = $scope.searchtext;
    	if ($scope.searchtext.length == 0)
    		t = " ";
    	Products.search(t, n, function(cs) {
        	$scope.currPage = n;
    		$scope.products = cs;
            $scope.products.forEach(function(elem, index, array) {
                elem.Info = JSON.parse(elem.Info);
            });
    	});
    };
}]);

product.controller('ProductEditController', ['$scope', '$routeParams', '$location', 'Global','StageData','Products', function ($scope, $routeParams, $location, Global,StageData, Products) {
    $scope.data = {};
    var id = $routeParams.id;
    var savedDataId = $routeParams.SavedDataId;
    var uploadedUrlsId = $routeParams.UploadedUrls;

    $scope.get = function() {
        if (savedDataId !== undefined && savedDataId !== '') {
            var stagedata = StageData.get(savedDataId);
            if (stagedata !== undefined) {
				$scope.id = stagedata.id
                $scope.data = stagedata.info;
                StageData.del(savedDataId);
            } else {
                Products.get(id , function(p) {
                    $scope.id = p.Id;
                    $scope.data = JSON.parse(p.Info);
                });
            }
            if (uploadedUrlsId !== undefined && uploadedUrlsId !== '') {
            	$scope.data.Photos = $scope.data.Photos.concat(StageData.get(uploadedUrlsId).split(";"));
                $scope.data.CoverPhoto = $scope.data.Photos[0];
                StageData.Del(uploadedUrlsId);
            }
        } else {
                Products.get(id , function(p) {
                    $scope.id = p.Id;
                    $scope.data = JSON.parse(p.Info);
                });
        }
    };


    $scope.save = function() {
        var price = parseFloat($scope.data.Price);
        var discount = parseFloat($scope.data.Discount);
        $scope.data.Price = price;
        $scope.data.Discount = discount;

        Products.update({"Id":$scope.id, "Info":JSON.stringify($scope.data)}, function(c) {
            $location.path("/products");
        });
    };

    $scope.jumptoupload = function() {
        var stageDataId = StageData.add({id:$scope.id,info:$scope.data});
        var r = $location.path().split("/");
        var redirecturl = "/" + r[1] + "/"+ r[2] + "/savedid/" + stageDataId;
        $location.path('/uploadfile/redirect/'+Base64.encode(redirecturl));
    };
}]);

product.controller('ProductAddController', ['$scope', '$routeParams', '$location', 'Global','StageData','Products', function ($scope, $routeParams, $location, Global,StageData, Products) {
    $scope.data = {};


    var savedDataId = $routeParams.SavedDataId;
    var uploadedUrlsId = $routeParams.UploadedUrls;
    $scope.get = function() {
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
                $scope.data.Photos = StageData.get(uploadedUrlsId).split(";");
                $scope.data.CoverPhoto = $scope.data.Photos[0];
                StageData.del(uploadedUrlsId);
            }
        }
    };


    $scope.save= function() {
        var price = parseFloat($scope.data.Price);
        var discount = parseFloat($scope.data.Discount);
        $scope.data.Price = price;
        $scope.data.Discount = discount;

        Products.add({"Info":JSON.stringify($scope.data)}, function(c) {
            $location.path("/products");
        });
    };

    $scope.jumptoupload = function() {
        var stageDataId = StageData.add($scope.data);
        var r = $location.path().split("/")[1];
        var redirecturl = "/" + r + "/savedid/" + stageDataId;
        $location.path('/uploadfile/redirect/'+Base64.encode(redirecturl));
    };
}]);

angular.module('featen.product').controller('ProductsController', ['$scope', '$routeParams', '$location', 'Global', 'Products', function ($scope, $routeParams, $location, Global, Products) {
    
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

angular.module('featen.product').controller('ProductEditController', ['$scope', '$routeParams', '$location', 'Global','StageData','Products', function ($scope, $routeParams, $location, Global,StageData, Products) {
    $scope.data = {};
    var name = $routeParams.Name;
    var savedDataId = $routeParams.SavedDataId;
    var uploadedUrlsId = $routeParams.UploadedUrls;

    $scope.get = function() {
        if (savedDataId !== undefined && savedDataId !== '') {
            var stagedata = StageData.get(savedDataId);
            if (stagedata !== undefined) {
                $scope.data = stagedata;
                StageData.del(savedDataId);
            } else {
                Products.get(name, function(p) {
                    $scope.data = p;
                    $scope.data.Info = JSON.parse($scope.data.Info);
                });
            }
            if (uploadedUrlsId !== undefined && uploadedUrlsId !== '') {
            	$scope.data.Photos = $scope.data.Photos.concat(StageData.get(uploadedUrlsId).split(";"));
                $scope.data.CoverPhoto = $scope.data.Photos[0];
                StageData.Del(uploadedUrlsId);
            }
        } else {
                Products.get(navname, function(p) {
                    $scope.data = p;
                    $scope.data.Info = JSON.parse($scope.data.Info);
                });
        } 
    };


    $scope.update = function() {
        var price = parseFloat($scope.data.Price);
        var discount = parseFloat($scope.data.Discount);
        $scope.data.Price = price;
        $scope.data.Discount = discount;

        Products.update($scope.data, function(c) {
            $location.path("/products");
        });
    };

    $scope.jumptoupload = function() {
        var stageDataId = StageData.add($scope.data);
        var r = $location.path().split("/");
        var redirecturl = "/" + r[1] + "/"+ r[2] + "/savedid/" + stageDataId;
        $location.path('/uploadfile/redirect/'+Base64.encode(redirecturl));
    };
}]);

angular.module('featen.product').controller('ProductAddController', ['$scope', '$routeParams', '$location', 'Global','StageData','Products', function ($scope, $routeParams, $location, Global,StageData, Products) {
    $scope.data = {};
	
    var savedDataId = $routeParams.SavedDataId;
    var uploadedUrlsId = $routeParams.UploadedUrls;
    $scope.getNewProduct = function() {
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


    $scope.add= function() {
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

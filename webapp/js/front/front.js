(function() {
var front = angular.module('front');

front.factory("Articles", ['$http', 'Alerts', function($http, Alerts) {
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
        return this;
}]);
front.factory("Products", ['$http', 'Alerts', function($http, Alerts) {
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

        return this;
}]);


front.factory("Dict", ['$http', 'Alerts', function($http, Alerts) {
        this.query = function(q, scall, ecall) {
            var promise = $http.get("/service/dict/" + q);
            var error = {
                type: "warning",
                strong: "Warning!",
                message: "Unable to query."
            };
            Alerts.handle(promise, error, undefined, scall, ecall);

            return promise;
        };

        return this;
    }]);
 

front.controller('PageController', ['$scope', '$routeParams', '$location', 'Global', 'Articles', 'Products', 'Dict', function($scope, $routeParams, $location, Global, Articles, Products, Dict) {
	$scope.currPage = 1;
	
	$scope.getPageArticles = function() {
		var p = parseInt($routeParams.Num);
		if (p != undefined && (p >= 1))
			$scope.currPage = p; 
		Articles.getTotalPageNumber(function(n) {
			$scope.totalPageNumber = parseInt(n);
		});
		Articles.getPageArticles($scope.currPage, function(ps) {
			$scope.articles = ps;
            $scope.articles.forEach(function(elem, index, array) {
                elem.Info = JSON.parse(elem.Info);
            });
		});
		Products.getrand(function(ds) {
			$scope.deals = ds;
            $scope.deals.forEach(function(elem, index, array) {
                elem.Info = JSON.parse(elem.Info);
            });
		});
	};
	
	$scope.setPage = function(page) {
		$scope.currPage = page;
		Articles.getPageArticles($scope.currPage, function(ps) {
			$scope.articles = ps;
		});
	};
        $scope.query = function() {
            var q = $scope.searchtext;
            if ($scope.searchtext.length == 0)
                return;
            Dict.query(q, function(r) {

                $scope.fanyi = r;
            });

        };
	
}]);


front.controller('ArticleViewController', ['$scope', '$routeParams', '$location', 'Global', 'Articles', 'Products', function($scope, $routeParams, $location, Global, Articles, Products) {
        $scope.global = Global;

        $scope.getblog = function() {
            var id = $routeParams.id;

            Articles.get(id, function(a) {
                $scope.article = a;
                $scope.info = JSON.parse(a.Info);
                var htmlcontent = marked($scope.info.Content);
                $('#htmlcontentdiv').html(htmlcontent);
            });
            
            Products.getrand(function(ds) {
    			$scope.deals = ds;
                $scope.deals.forEach(function(elem, index, array) {
                    elem.Info = JSON.parse(elem.Info);
                });
    		});
        };
    }]);

front.controller('DealsController', ['$scope', '$routeParams', '$location', 'Global', 'Products', function ($scope, $routeParams, $location, Global, Products) {
    $scope.global = Global;
    /*----------------------------------------------------*/
    /*  Flexslider
    /*----------------------------------------------------*/
    var loadSlider = function() {
    $('#intro-slider').flexslider({
          namespace: "flex-",
          controlsContainer: "",
          animation: 'fade',
          controlNav: false,
          directionNav: true,
          smoothHeight: true,
          slideshowSpeed: 7000,
          animationSpeed: 600,
          randomize: false,
       });
    };

    $scope.load = function() {
        var id = $routeParams.id;
        Products.get(id, function(p) {
            $scope.deal = p;
            $scope.deal.Info = JSON.parse($scope.deal.Info);
            $scope.photos = $scope.deal.Info.Photos;
            setTimeout( loadSlider, 1);

            var htmlintro = marked($scope.deal.Info.Introduction);
            $('#introduction').html(htmlintro);
            var htmlspec = marked($scope.deal.Info.Spec);
            $('#specs').html(htmlspec);
        });
    };
}]);


})();

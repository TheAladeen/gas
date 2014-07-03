angular.module('featen.article').controller('PageController', ['$scope', '$routeParams', '$location', 'Global', 'Articles', 'Products', 'Dict', function($scope, $routeParams, $location, Global, Articles, Products, Dict) {
	$scope.currPage = 1;
	
	$scope.getPageArticles = function() {
		(adsbygoogle = window.adsbygoogle || []).push({});

		var p = parseInt($routeParams.PageNum);
		if (p != undefined && (p >= 1))
			$scope.currPage = p; 
		Articles.getTotalPageNumber(function(n) {
			$scope.totalPageNumber = parseInt(n);
		});
		Articles.getPageArticles($scope.currPage, function(ps) {
			$scope.articles = ps;
            $scope.articles.forEach(function(elem, index, array) {
                elem.Info = JSON.parse(elem.Info);
                var ct = elem.Info.CreateTime;
                elem.Info.CreateTime = d3.time.format('%Y-%m-%d %I:%M%p')(new Date(ct));
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


angular.module('featen.article').controller('ArticleViewController', ['$scope', '$routeParams', '$location', 'Global', 'Articles', 'Products', function($scope, $routeParams, $location, Global, Articles, Products) {
        $scope.global = Global;

        $scope.getblog = function() {
		(adsbygoogle = window.adsbygoogle || []).push({});
            var nav= $routeParams.Nav;

            Articles.get(nav, function(a) {
                $scope.article = a;
                $scope.info = JSON.parse(a.Info);
                $scope.info.CreateTime = d3.time.format('%Y-%m-%d %I:%M%p')(new Date($scope.info.CreateTime));
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

angular.module('featen.deal').controller('DealsController', ['$scope', '$routeParams', '$location', 'Global', 'Products', function ($scope, $routeParams, $location, Global, Products) {
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
        var name = $routeParams.Name;
        Products.get(name, function(p) {
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



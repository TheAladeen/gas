angular.module('featen.article').controller('ArticlesController', ['$scope', '$routeParams', '$location', 'Global', 'Articles', function($scope, $routeParams, $location, Global, Articles) {
        $scope.global = Global;


        $scope.render = function(c){
            return marked(c); 
        };
        $scope.getall = function() {
            Articles.getall(function(ps) {
                $scope.articles = ps;
            });

        };


    }]);

angular.module('featen.article').controller('PageController', ['$scope', '$routeParams', '$location', 'Global', 'Articles', 'Deals', function($scope, $routeParams, $location, Global, Articles, Deals) {
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
		Deals.getPageDeals($scope.currPage, function(ds) {
			$scope.deals = ds;
		});
	};
	
	$scope.setPage = function(page) {
		$scope.currPage = page;
		Articles.getPageArticles($scope.currPage, function(ps) {
			$scope.articles = ps;
		});
	};
	
}]);


angular.module('featen.article').controller('ArticleViewController', ['$scope', '$routeParams', '$location', 'Global', 'Articles', 'Deals', function($scope, $routeParams, $location, Global, Articles, Deals) {
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
            
            Deals.getPageDeals(1, function(ds) {
    			$scope.deals = ds;
    		});
        };
    }]);

angular.module('featen.article').controller('ArticleEditController', ['$scope', '$window', '$document', '$routeParams', '$location', 'Global', 'StageData', 'Articles', function($scope, $window, $document, $routeParams, $location, Global, StageData, Articles) {
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
            Articles.create({"Info":JSON.stringify({"Title": $scope.data.title, "Nav": $scope.data.navname.replace(/ /g,"_"), "CoverPhoto": $scope.data.coverphoto, "Intro": $scope.data.intro, "Content": $scope.data.content, "CreateTime": Date.now()})}, function(l) {
                $location.path('/');
            });
        }; 
        }]);

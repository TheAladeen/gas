angular.module('featen.dict').controller('DictController', ['$scope', '$routeParams', '$location', 'Global', 'Dict', function($scope, $routeParams, $location, Global, Dict) {
        $scope.global = Global;

	$scope.init = function() {
		(adsbygoogle = window.adsbygoogle || []).push({});
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


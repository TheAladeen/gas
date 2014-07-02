angular.module('featen.deal').controller('DealsController', ['$scope', '$routeParams', '$location', 'Global', 'Products', function ($scope, $routeParams, $location, Global, Products) {
    $scope.global = Global;
    /*----------------------------------------------------*/
    /*	Flexslider
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

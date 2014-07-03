angular.module('featen.user').controller('SignInController', ['$scope', '$routeParams', '$route', '$rootScope', '$location', 'Global', 'User', function($scope, $routeParams, $route, $rootScope, $location, Global, User) {
        $scope.data = {email: '', password: ''};
        $scope.signin = function() {
            User.signin({"NameOrEmail": $scope.data.email, "Passwd": $scope.data.password}, function(u) {
                Global.user = u;
                Global.authenticated = true;
                $location.path('/');
            });
        };
    }]);


angular.module('featen.user').controller('UsersCtrl', ['$scope', '$routeParams', '$route', '$location', 'Global', 'User', 'Alerts', function($scope, $routeParams, $route, $location, Global, User, Alerts) {
        $scope.getCurrentUser = function() {
            User.get(function(data) {
                Global.user = data;
                Global.authenticated = true;
            });
        };

        $scope.init = function() {
            $scope.getCurrentUser();            
            $scope.global = Global;
        };

        $scope.signout = function() {
            User.signout(function() {
                Global.user = null;
                Global.authenticated = false;
                $location.path('/');
            });
        };
    }]);

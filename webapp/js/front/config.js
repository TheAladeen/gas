//Setting up route
window.app.config(['$routeProvider', '$locationProvider',
    function($routeProvider, $locationProvider) {
        $routeProvider.
                when('/blog/:id', {
                    templateUrl: '/views/viewblog.html'
                }).
                when('/page/:Num', {
                    templateUrl: '/views/page.html'
                }).
                when('/deal/:id', {
                    templateUrl: '/views/deal.html'
                }).
                when('/404', {
                    templateUrl: '/views/404.html'
                }).
                when('/', {
                    templateUrl: '/views/page.html'
                }).
                otherwise({
                    redirectTo: '/'
                });

                if (window.history && window.history.pushState) {
          $locationProvider.html5Mode(true);
                }
    }
]);


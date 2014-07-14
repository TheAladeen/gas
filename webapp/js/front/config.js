//Setting up route
window.app.config(['$routeProvider', '$locationProvider',
    function($routeProvider, $locationProvider) {
        $routeProvider.
                when('/page/:Num', {
                    templateUrl: '/views/blogs.html'
                }).
                when('/blog/:id', {
                    templateUrl: '/views/viewblog.html'
                }).
                when('/deal/:id', {
                    templateUrl: '/views/deal.html'
                }).
                when('/404', {
                    templateUrl: '/views/404.html'
                }).
                when('/', {
                    templateUrl: '/views/blogs.html'
                }).
                otherwise({
                    redirectTo: '/'
                });
/*
                if (window.history && window.history.pushState) {
          $locationProvider.html5Mode(true);
                }
               */
    }
]);


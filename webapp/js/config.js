//Setting up route
window.app.config(['$routeProvider',
    function($routeProvider) {
        $routeProvider.
                when('/blog/:Nav', {
                    templateUrl: '/views/viewblog.html'
                }).
                when('/page/:PageNum', {
                	templateUrl: '/views/page.html'
                }).
                when('/deal/:Name', {
                    templateUrl: '/views/deal.html'
                }).
                when('/', {
                    templateUrl: '/views/page.html'
                }).
                otherwise({
                    redirectTo: '/'
                });
    }
]);

//Setting HTML5 Location Mode
window.app.config(['$locationProvider',
    function($locationProvider) {
        $locationProvider.hashPrefix("!");
    }
]);

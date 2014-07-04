//Setting up route
window.app.config(['$routeProvider',
    function($routeProvider) {
        $routeProvider.
                when('/blog/:id', {
                    templateUrl: '/views/viewblog.html'
                }).
                when('/page/:PageNum', {
                	templateUrl: '/views/page.html'
                }).
                when('/deal/:id', {
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


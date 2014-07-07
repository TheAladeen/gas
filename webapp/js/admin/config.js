//Setting up route
window.app.config(['$routeProvider', '$locationProvider',
    function($routeProvider, $locationProvider) {
        $routeProvider.
                when('/signin', {
                    templateUrl: '/views/signin.html'
                }).
                when('/uploadfile/redirect/:ReUrl', {
                    templateUrl: '/views/dropboxupload.html'
                }).
                when('/dropboxupload', {
                    templateUrl: '/views/dropboxupload.html'
                }).
                when('/writeblog', {
                    templateUrl: '/views/writeblog.html'
                }).
                when('/writeblog/savedid/:SavedDataId/uploaded/:UploadedUrls', {
                	templateUrl: '/views/writeblog.html'
                }).
                when('/products', {
                    templateUrl: '/views/products.html'
                }).
                when('/addproduct', {
					controller: 'ProductAddController',
                    templateUrl: '/views/editproduct.html'
                }).
                when('/addproduct/savedid/:SavedDataId/uploaded/:UploadedUrls', {
					controller: 'ProductAddController',
                    templateUrl: '/views/editproduct.html'
                }).
                when('/product/:id', {
					controller: 'ProductEditController',
                    templateUrl: '/views/editproduct.html'
                }).
                when('/product/:id/savedid/:SavedDataId/uploaded/:UploadedUrls', {
					controller: 'ProductEditController',
                    templateUrl: '/views/editproduct.html'
                }).
                when('/404', {
                    templateUrl: '/views/404.html'
                }).
                when('/', {
                    templateUrl: '/views/dashboard.html'
                }).
                otherwise({
                    redirectTo: '/'
                });

//                if (window.history && window.history.pushState) {
//                    $locationProvider.html5Mode(true);
//                }
    }
]);


//Setting up route
window.app.config(['$routeProvider',
    function($routeProvider) {
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
                otherwise({
                    redirectTo: '/products'
                });
    }
]);


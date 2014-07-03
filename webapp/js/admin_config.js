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
                    templateUrl: '/views/addproduct.html'
                }).
                when('/addproduct/savedid/:SavedDataId/uploaded/:UploadedUrls', {
                    templateUrl: '/views/addproduct.html'
                }).
                when('/product/:Name', {
                    templateUrl: '/views/editproduct.html'
                }).
                when('/product/:Name/savedid/:SavedDataId/uploaded/:UploadedUrls', {
                    templateUrl: '/views/editproduct.html'
                }).
                when('/addcustomer', {
                    templateUrl: '/views/addcustomer.html'
                }).
                when('/addcustomer/savedid/:SavedDataId/uploaded/:UploadedUrls', {
                    templateUrl: '/views/addcustomer.html'
                }).
                when('/customers', {
                    templateUrl: '/views/customers.html'
                }).
                when('/customer/:Id', {
                    templateUrl: '/views/editcustomer.html'
                }).
                when('/customer/:Id/savedid/:SavedDataId/uploaded/:UploadedUrlsId', {
                    templateUrl: '/views/editcustomer.html'
                }).
                when('/', {
                    templateUrl: '/views/customers.html'
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

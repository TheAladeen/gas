angular.module('featen.agent').controller('AgentsController', ['$scope', '$routeParams', '$location', 'Global', 'Agents', function($scope, $routeParams, $location, Global, Agents) {
        $scope.global = Global;


        $scope.render = function(c){
            return marked(c);
        };
        $scope.getall = function() {
            Agents.getall(function(ps) {
                $scope.agents = ps;
            });

        };


    }]);

angular.module('featen.agent').controller('AgentViewController', ['$scope', '$routeParams', '$location', 'Global', 'Agents', function($scope, $routeParams, $location, Global, Agents) {
        $scope.global = Global;

        $scope.get = function() {
            var id = $routeParams.Id;

            Agents.get(id, function(a) {
                $scope.info = JSON.parse(a.Info);
                var htmlcontent = marked($scope.info.Content);
                $('#htmlcontentdiv').html(htmlcontent);
            });
        };
    }]);

angular.module('featen.agent').controller('AgentEditController', ['$scope', '$window', '$document', '$routeParams', '$location', 'Global',  'Agents', function($scope, $window, $document, $routeParams, $location, Global, Agents) {

        //$scope.global = Global;
        $scope.data = {};

        $scope.create = function() {
            Agents.create({"Info":JSON.stringify({"Title": $scope.data.title, "NavName": $scope.data.navname.replace(/ /g,"_"), "CoverPhoto": $scope.data.coverphoto, "Intro": $scope.data.intro, "Content": $scope.data.content})}, function(l) {
                $location.path('/');
            });
        };
        }]);

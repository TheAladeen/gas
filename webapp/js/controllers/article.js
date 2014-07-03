angular.module('featen.article').controller('ArticleEditController', ['$scope', '$window', '$document', '$routeParams', '$location', 'Global', 'StageData', 'Articles', function($scope, $window, $document, $routeParams, $location, Global, StageData, Articles) {
        $scope.data = {};
        $scope.editstate = true;
        
        $scope.preview = function(){
        	if ($scope.data.content !== undefined && $scope.data.content !== "") {
        		var htmlcontent = marked($scope.data.content);
                $('#htmlcontentdiv').html(htmlcontent);	
        	} else {
        		$('#htmlcontentdiv').html("");
        	}
            
            $scope.editstate = false;
        };

        $scope.edit = function() {
            $scope.editstate = true;
        };

        var savedDataId = $routeParams.SavedDataId;
        var uploadedUrlsId = $routeParams.UploadedUrls;
        $scope.getEditingArticle = function() {
            if (savedDataId !== undefined && savedDataId !== '') {
                var stagedata = StageData.get(savedDataId);
                if (stagedata !== undefined) {
                    $scope.data = stagedata;
                    StageData.del(savedDataId);
                } else {
                    var r = $location.path().split("/")[1];
                    $location.path("/" + r);
                }
                if (uploadedUrlsId !== undefined && uploadedUrlsId !== '') {
                    $scope.data.blogphotos = StageData.get(uploadedUrlsId).split(";");
                    $scope.data.coverphoto = $scope.data.blogphotos[0];
                    StageData.del(uploadedUrlsId);
                }
            } 
        };
        
        
        $scope.jumptoupload = function() {
            var stageDataId = StageData.add($scope.data);
            var r = $location.path().split("/")[1];
            var redirecturl = "/" + r + "/savedid/" + stageDataId;
            $location.path('/uploadfile/redirect/'+Base64.encode(redirecturl));
        };
      
        $scope.create = function() {
            Articles.create({"Info":JSON.stringify({"Title": $scope.data.title, "Nav": $scope.data.navname.replace(/ /g,"_"), "CoverPhoto": $scope.data.coverphoto, "Intro": $scope.data.intro, "Content": $scope.data.content, "CreateTime": Date.now()})}, function(l) {
                $location.path('/');
            });
        }; 
        }]);
